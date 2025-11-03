package services

import (
	"school-exam-system/internal/models"

	"gorm.io/gorm"
)

type StudentService struct {
	db *gorm.DB
}

func NewStudentService(db *gorm.DB) *StudentService {
	return &StudentService{db: db}
}

func (s *StudentService) GetStudentByUserID(userID uint) (*models.Student, error) {
	var student models.Student
	if err := s.db.Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (s *StudentService) GetStudentProfile(userID uint) (*models.Student, error) {
	var student models.Student
	if err := s.db.Preload("User").Preload("Class").Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (s *StudentService) GetStudentSubjects(userID uint) ([]models.Subject, error) {
	student, err := s.GetStudentByUserID(userID)
	if err != nil {
		return nil, err
	}

	var subjects []models.Subject
	if err := s.db.Table("subjects").
		Joins("JOIN class_subjects ON subjects.id = class_subjects.subject_id").
		Where("class_subjects.class_id = ? AND class_subjects.is_active = ? AND subjects.is_active = ?", 
			student.ClassID, true, true).
		Find(&subjects).Error; err != nil {
		return nil, err
	}

	return subjects, nil
}

func (s *StudentService) GetClassmates(userID uint) ([]models.Student, error) {
	student, err := s.GetStudentByUserID(userID)
	if err != nil {
		return nil, err
	}

	var classmates []models.Student
	if err := s.db.Preload("User").Where("class_id = ? AND user_id != ?", student.ClassID, userID).Find(&classmates).Error; err != nil {
		return nil, err
	}

	return classmates, nil
}