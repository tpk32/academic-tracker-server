package router
import(
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tpk32/academic-tracker-server/controllers"
)

func Routes() http.Handler{
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
        AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders: []string{"Link"},
        AllowCredentials: true,
        MaxAge: 300,
	}))

	router.Post("/api/v1/student/login", controllers.StoreStudentInformation)
	router.Post("/api/v1/student", controllers.CreateStudent)
	router.Delete("/api/v1/student/{student_id}", controllers.DeleteStudent)
	router.Post("/api/v1/student/subject", controllers.CreateSubject)
	router.Delete("/api/v1/student/subject/{student_id}/{subject_name}", controllers.DeleteSubject)
	router.Delete("/api/v1/student/subject/{student_id}", controllers.DeleteAllSubjectsByStudentId)
	router.Get("/api/v1/student/subject/{student_id}", controllers.GetAllSubjectsByStudentId)
	router.Post("/api/v1/student/subject/test/{student_id}/{subject_name}", controllers.CreateTest)
	router.Delete("/api/v1/student/subject/test/{student_id}/{subject_name}/{test_name}/{test_date}", controllers.DeleteTest)
	router.Get("/api/v1/student/test/{student_id}", controllers.GetAllTestsByStudentId)
	router.Get("/api/v1/student/test/{student_id}/{subject_name}", controllers.GetAllTestsBySubjectName)
	return router
}
