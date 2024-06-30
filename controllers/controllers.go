package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rickb777/date"
	"github.com/tpk32/academic-tracker-server/helpers"
	"github.com/tpk32/academic-tracker-server/services"
)

var student services.Student
var subject services.Subject
var test services.Test

// POST(/api/v1/student/login)
func StoreStudentInformation(w http.ResponseWriter, r *http.Request){
	var request struct{
		IDToken string `json:"idToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err!=nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	client, err := services.InitializeApp()
	if err != nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
	}

	token, err := services.VerifyIDToken(client, request.IDToken)
    if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	studentExists, student, err := services.StoreStudentInformation(token)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := services.JsonResponse{
		Error: false,
        Message: "Student exists",
        Data: student,
    }
    if !studentExists {
        response.Message = "Student created"
    }

    w.Header().Set("Content-Type", "application/json")

	if err := helpers.WriteJSON(w, http.StatusOK, response); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

}

// POST(/api/v1/student)
func CreateStudent(w http.ResponseWriter, r *http.Request){
	var studentData services.Student
	err := json.NewDecoder(r.Body).Decode(&studentData)
	if err != nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	studentCreated, err := student.CreateStudent(studentData)
	if err != nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, studentCreated)
}

// DELETE(/api/v1/student/{student_id})
func DeleteStudent(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	err := student.DeleteStudent(student_id)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"message":"successful student deletion"})
}

// POST(/api/v1/student/subject)
func CreateSubject(w http.ResponseWriter, r *http.Request){
	var subjectData services.Subject
	err := json.NewDecoder(r.Body).Decode(&subjectData)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	subjectCreated, err := subject.CreateSubject(subjectData)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, subjectCreated);
}

// DELETE(/api/v1/student/subject/{student_id}/{subject_name})
func DeleteSubject(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	subject_name := chi.URLParam(r, "subject_name")
	err := subject.DeleteSubject(student_id, subject_name)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"message":"successful subject deletion"})
}

// GET(/api/v1/student/subject/{student_id})
func GetAllSubjectsByStudentId(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	allSubjects, err := subject.GetAllSubjectsByStudentId(student_id)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"subjects": allSubjects})
}

// DELETE(/api/v1/student/subject/{student_id})
func DeleteAllSubjectsByStudentId(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	err := subject.DeleteAllSubjectsByStudentId(student_id)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"message":"successful all subjects deletion"})
}

// POST(/api/v1/student/subject/test/{student_id}/{subject_name})
func CreateTest(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	subject_name := chi.URLParam(r, "subject_name")
	var testData services.Test
	err := json.NewDecoder(r.Body).Decode(&testData)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	testCreated, err := test.CreateTest(student_id, subject_name, testData)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, testCreated);
}

// DELETE(/api/v1/student/subject/test/{student_id}/{subject_name}/{test_name}/{test_date})
func DeleteTest(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	subject_name := chi.URLParam(r, "subject_name")
	test_name := chi.URLParam(r, "test_name")
	test_date_str := chi.URLParam(r, "test_date")
	test_date, err := date.AutoParse(test_date_str)

	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	err = test.DeleteTest(student_id, subject_name, test_name, test_date)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"message":"successful subject deletion"})
}

// GET(/api/v1/student/test/{student_id})
func GetAllTestsByStudentId(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")

	allTests, err := test.GetAllTestsByStudentId(student_id)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"tests": allTests});
}

// GET(/api/v1/student/test/{student_id}/{subject_name})
func GetAllTestsBySubjectName(w http.ResponseWriter, r *http.Request){
	student_id := chi.URLParam(r, "student_id")
	subject_name := chi.URLParam(r, "subject_name")

	allTests, err := test.GetAllTestsBySubjectName(student_id, subject_name)
	if err!=nil{
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"tests": allTests});
}