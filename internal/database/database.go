package database

import (
	"fmt"
	"school-exam-system/internal/config"
	"school-exam-system/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initialize(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Seed initial data
	if err := SeedDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	// Migrate tables in correct order to avoid foreign key issues
	// First migrate base tables without foreign keys
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.Class{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.Subject{})
	if err != nil {
		return err
	}
	
	// Then migrate tables with foreign keys
	err = db.AutoMigrate(&models.Teacher{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.ClassSubject{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.Exam{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.Question{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.ExamSubmission{})
	if err != nil {
		return err
	}
	
	err = db.AutoMigrate(&models.SubmissionAnswer{})
	if err != nil {
		return err
	}

	// Add custom indexes and constraints
	return addIndexesAndConstraints(db)
}

func addIndexesAndConstraints(db *gorm.DB) error {
	// Add unique constraint for class_subjects (one teacher per subject per class)
	if err := db.Exec(`
		ALTER TABLE class_subjects 
		ADD CONSTRAINT unique_class_subject_teacher 
		UNIQUE (class_id, subject_id)
	`).Error; err != nil {
		// Ignore error if constraint already exists
	}

	// Add unique constraint for exam_submissions (prevent duplicate submissions)
	if err := db.Exec(`
		ALTER TABLE exam_submissions 
		ADD CONSTRAINT unique_student_exam_attempt 
		UNIQUE (exam_id, student_id, attempt_num)
	`).Error; err != nil {
		// Ignore error if constraint already exists
	}

	// Add unique constraint for submission_answers
	if err := db.Exec(`
		ALTER TABLE submission_answers 
		ADD CONSTRAINT unique_submission_question 
		UNIQUE (submission_id, question_id)
	`).Error; err != nil {
		// Ignore error if constraint already exists
	}

	// Add performance indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
		"CREATE INDEX IF NOT EXISTS idx_teachers_employee_id ON teachers(employee_id)",
		"CREATE INDEX IF NOT EXISTS idx_students_student_id ON students(student_id)",
		"CREATE INDEX IF NOT EXISTS idx_students_class_id ON students(class_id)",
		"CREATE INDEX IF NOT EXISTS idx_exams_class_subject ON exams(class_id, subject_id)",
		"CREATE INDEX IF NOT EXISTS idx_exams_teacher_id ON exams(teacher_id)",
		"CREATE INDEX IF NOT EXISTS idx_exams_status ON exams(status)",
		"CREATE INDEX IF NOT EXISTS idx_questions_exam_id ON questions(exam_id)",
		"CREATE INDEX IF NOT EXISTS idx_submissions_student_exam ON exam_submissions(student_id, exam_id)",
		"CREATE INDEX IF NOT EXISTS idx_submissions_status ON exam_submissions(status)",
		"CREATE INDEX IF NOT EXISTS idx_answers_submission_id ON submission_answers(submission_id)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			// Log error but continue with other indexes
			continue
		}
	}

	return nil
}