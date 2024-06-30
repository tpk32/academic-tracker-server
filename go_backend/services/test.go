package services

import (
	"fmt"
	"context"
	"time"

	"github.com/rickb777/date"
)

type Test struct{
	TestId int64 `json:"test_id"`
	SubjectId int64 `json:"subject_id"`
	TestName string `json:"test_name"`
	TestDate date.Date `json:"test_date"`
	MaxMarks int16 `json:"max_marks"`
	ObtainedMarks int16 `json:"obtained_marks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SyncStatus string `json:"sync_status"`
}

func(t* Test) CreateTest(student_id string, subject_name string, test Test) (*Test, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.Begin()
	if err!=nil {
		fmt.Println("error at 1")
		return nil, err

	}

	testQuery := `
		WITH SubjectCTE AS(
			SELECT subject_id
			FROM Subject
			WHERE student_id=$1 AND subject_name=$2
		)
		INSERT INTO Test (subject_id, test_name, test_date, max_marks, obtained_marks, created_at, updated_at, sync_status)
		SELECT subject_id, $3, $4, $5, $6, $7, $8, $9
		FROM SubjectCTE
		RETURNING *
	`

	_, err = tx.ExecContext(
		ctx,
		testQuery,
		student_id,
		subject_name,
		test.TestName,
		test.TestDate.String(),
		test.MaxMarks,
		test.ObtainedMarks,
		time.Now(),
		time.Now(),
		"synced",
	)

	if err != nil{
		tx.Rollback()
		fmt.Println("error at 2")
		return nil, err
	}

	err = tx.Commit()
	if err!=nil{
		fmt.Println("error at 3")
		return nil, err
	}
	return &test, nil
}

func (t* Test) GetAllTestsByStudentId(student_id string) ([]*Test, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT t.test_id, t.subject_id, t.test_name, t.test_date, t.max_marks, t.obtained_marks, t.created_at, t.updated_at, t.sync_status
		FROM Test t
		INNER JOIN Subject s ON t.subject_id = s.subject_id
		INNER JOIN Student stu ON s.student_id = stu.student_id
		WHERE stu.student_id = $1;
	`
	rows, err := db.QueryContext(ctx, query, student_id)
	if err != nil{
		return nil, err
	}
	
	var tests []*Test
	for rows.Next(){
		var test Test
		err := rows.Scan(
			&test.TestId,
			&test.SubjectId,
			&test.TestName,
			&test.TestDate,
			&test.MaxMarks,
			&test.ObtainedMarks,
			&test.CreatedAt,
			&test.UpdatedAt,
			&test.SyncStatus,
		)
		if err!=nil {
			return nil ,err
		}
		tests = append(tests, &test)
	}
	
	return tests, nil
} 

func (t* Test) GetAllTestsBySubjectName(student_id string, subject_name string) ([]*Test, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT t.test_id, t.subject_id, t.test_name,
		t.test_date, t.max_marks, t.obtained_marks,
		t.created_at, t.updated_at, t.sync_status
		FROM Test t
            INNER JOIN Subject s ON t.subject_Id = s.subject_Id
            INNER JOIN Student stu ON stu.student_Id = s.student_Id
            WHERE stu.student_Id = $1 AND s.subject_name = $2
	`
	rows, err := db.QueryContext(ctx, query, student_id, subject_name)
	if err != nil{
		return nil, err
	}
	
	var tests []*Test
	for rows.Next(){
		var test Test
		err := rows.Scan(
			&test.TestId,
			&test.SubjectId,
			&test.TestName,
			&test.TestDate,
			&test.MaxMarks,
			&test.ObtainedMarks,
			&test.CreatedAt,
			&test.UpdatedAt,
			&test.SyncStatus,
		)
		if err!=nil {
			return nil ,err
		}
		tests = append(tests, &test)
	}
	
	return tests, nil
}

func (t *Test) DeleteTest(student_id string, subject_name string, test_name string, test_date date.Date) error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	deleteTestQuery := `
		DELETE FROM Test
		WHERE subject_id = (
			SELECT subject_id
			FROM Subject
			WHERE student_id=$1 AND subject_name=$2
		) AND test_name=$3 AND test_date=$4
	`
	_, err := db.ExecContext(
		ctx, 
		deleteTestQuery, 
		student_id, 
		subject_name, 
		test_name, 
		test_date.String(),
	)
	if err!=nil {
		return err
	}
	return nil
}