package router

import (
	"school-exam-system/internal/config"
	"school-exam-system/internal/handlers"
	"school-exam-system/internal/middleware"
	"school-exam-system/internal/models"
	"school-exam-system/internal/services"
	"school-exam-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Initialize services
	jwtService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHour)
	authService := services.NewAuthService(db, jwtService)
	userService := services.NewUserService(db)
	academicService := services.NewAcademicService(db)
	teacherService := services.NewTeacherService(db)
	studentService := services.NewStudentService(db)
	examService := services.NewExamService(db)
	questionService := services.NewQuestionService(db)
	submissionService := services.NewSubmissionService(db)
	gradingService := services.NewGradingService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(userService, academicService)
	teacherHandler := handlers.NewTeacherHandler(academicService, teacherService)
	studentHandler := handlers.NewStudentHandler(studentService, academicService)
	examHandler := handlers.NewExamHandler(examService, questionService, submissionService, gradingService, teacherService, studentService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "School Exam System is running",
		})
	})

	// Auth routes without /api prefix for compatibility
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(authService), authHandler.Profile)
	}

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		authAPI := api.Group("/auth")
		{
			authAPI.POST("/login", authHandler.Login)
			authAPI.POST("/logout", authHandler.Logout)
			authAPI.GET("/profile", middleware.AuthMiddleware(authService), authHandler.Profile)
		}

		// Admin routes (protected)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.RequireRole(models.RoleSystemAdmin))
		{
			// User management
			admin.POST("/users", adminHandler.CreateUser)
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)
			
			// Class management
			admin.POST("/classes", adminHandler.CreateClass)
			admin.GET("/classes", adminHandler.GetClasses)
			admin.GET("/classes/:classId/roster", adminHandler.GetClassRoster)
			
			// Subject management
			admin.POST("/subjects", adminHandler.CreateSubject)
			admin.GET("/subjects", adminHandler.GetSubjects)
			
			// Teacher assignment and workload
			admin.POST("/assign-teacher", adminHandler.AssignTeacher)
			admin.GET("/teachers/:teacherId/workload", adminHandler.GetTeacherWorkload)
			
			// Student enrollment
			admin.POST("/students/:studentId/enroll", adminHandler.EnrollStudent)
			admin.POST("/students/:studentId/transfer", adminHandler.TransferStudent)
		}

		// Teacher routes (protected)
		teacher := api.Group("/teacher")
		teacher.Use(middleware.AuthMiddleware(authService))
		teacher.Use(middleware.RequireRole(models.RoleTeacher))
		{
			teacher.GET("/profile", teacherHandler.GetProfile)
			teacher.GET("/classes", teacherHandler.GetMyClasses)
			teacher.GET("/students", teacherHandler.GetMyStudents)
			teacher.GET("/classes/:classId/students", teacherHandler.GetStudentsByClass)
			
			// Exam management
			teacher.POST("/exams", examHandler.CreateExam)
			teacher.GET("/exams", examHandler.GetMyExams)
			teacher.GET("/exams/:id", examHandler.GetExam)
			teacher.PUT("/exams/:id", examHandler.UpdateExam)
			teacher.DELETE("/exams/:id", examHandler.DeleteExam)
			teacher.POST("/exams/:id/publish", examHandler.PublishExam)
			
			// Question management
			teacher.POST("/exams/:id/questions", examHandler.AddQuestion)
			teacher.GET("/exams/:id/questions", examHandler.GetQuestions)
			
			// Grading
			teacher.GET("/exams/:id/submissions", examHandler.GetSubmissions)
		}

		// Student routes (protected)
		student := api.Group("/student")
		student.Use(middleware.AuthMiddleware(authService))
		student.Use(middleware.RequireRole(models.RoleStudent))
		{
			student.GET("/profile", studentHandler.GetProfile)
			student.GET("/subjects", studentHandler.GetMySubjects)
			student.GET("/class", studentHandler.GetMyClass)
			student.GET("/classmates", studentHandler.GetClassmates)
			
			// Exam taking
			student.GET("/exams", examHandler.GetAvailableExams)
			student.POST("/exams/:id/start", examHandler.StartExam)
			student.POST("/submissions/:submissionId/answer", examHandler.SubmitAnswer)
			student.POST("/submissions/:submissionId/finish", examHandler.FinishExam)
			student.GET("/grades", examHandler.GetMyGrades)
		}
	}

	return r
}