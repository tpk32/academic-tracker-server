package services

import (
	"context"
	"database/sql"
	"time"

	"firebase.google.com/go/v4/auth"
)

type Student struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func StoreStudentInformation(token *auth.Token) (bool, *Student, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	studentID := token.UID
	name := token.Claims["name"].(string)
	email := token.Claims["email"].(string)

	var existingStudent Student
	query := `SELECT student_id, name, email FROM Student WHERE student_id=$1`
	err := db.QueryRow(query, studentID).Scan(
		&existingStudent.StudentID,
		&existingStudent.Name,
		&existingStudent.Email,
	)
	if err== sql.ErrNoRows{
		//student does not exist
		newStudent := &Student{
			StudentID: studentID,
			Name: name,
			Email: email,
		}

		insertQuery :=`INSERT INTO Student (student_id, name, email, created_at)
			VALUES ($1, $2, $3, $4) returning *
		`
		_, err := db.ExecContext(
			ctx,
			insertQuery,
			newStudent.StudentID,
			newStudent.Name,
			newStudent.Email,
			time.Now(),
		)
		if err!=nil{
			return false, nil, err
		}

		return false, newStudent, nil

	}else if err!=nil{
		return false, nil, err
	}

	//Student Exists
	return true, &existingStudent, nil

}

func (s *Student) CreateStudent(student Student) (*Student, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	INSERT INTO Student (student_id, name, email, created_at)
	VALUES ($1, $2, $3, $4) returning *
	`
	_, err := db.ExecContext(
		ctx,
		query,
		student.StudentID,
		student.Name,
		student.Email,
		time.Now(),
	)
	if err != nil{
		return nil, err
	}

	return &student, nil
}

func (s *Student) DeleteStudent(student_id string) error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM student WHERE student_id=$1`
	_, err := db.ExecContext(ctx, query, student_id)
	if err!=nil {
		return err
	}
	return nil
}