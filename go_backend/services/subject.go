package services

import (
	"context"	
	"time"
)

type Subject struct {
	SubjectID   int64  `json:"subject_id"`
	StudentID   string `json:"student_id"`
	SubjectName string `json:"subject_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SyncStatus   string `json:"sync_status"`
}

func (s* Subject) CreateSubject(subject Subject) (*Subject, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	subjectQuery := `
		INSERT INTO Subject (student_id, subject_name, created_at, updated_at, sync_status)
		VALUES ($1, $2, $3, $4, $5) RETURNING *
	`

	_, err := db.ExecContext(
		ctx,
		subjectQuery,
		subject.StudentID,
		subject.SubjectName,
		time.Now(),
		time.Now(),
		"synced",
	)
	if err != nil{
		return nil, err
	}

	return &subject, nil
}

func (s* Subject) GetAllSubjectsByStudentId(student_id string) ([]*Subject, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT subject_id, student_id, subject_name, created_at, updated_at, sync_status
		FROM Subject
		WHERE student_id = $1
	`
	rows, err := db.QueryContext(ctx, query, student_id)
	if err != nil{
		return nil, err
	}
	
	var subjects []*Subject
	for rows.Next(){
		var subject Subject
		err := rows.Scan(
			&subject.SubjectID,
			&subject.StudentID,
			&subject.SubjectName,
			&subject.CreatedAt,
			&subject.UpdatedAt,
			&subject.SyncStatus,
		)
		if err!=nil {
			return nil ,err
		}
		subjects = append(subjects, &subject)
	}
	
	return subjects, nil
} 

func (s *Subject) DeleteSubject(student_id string, subject_name string) error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM Subject WHERE student_id=$1 AND subject_name=$2`
	_, err := db.ExecContext(ctx, query, student_id, subject_name)
	if err!=nil {
		return err
	}
	return nil
}

func(s *Subject) DeleteAllSubjectsByStudentId(student_id string) error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM Subject WHERE student_id=$1`
	_, err := db.ExecContext(ctx, query, student_id)
	if err!=nil {
		return err
	}
	return nil
}