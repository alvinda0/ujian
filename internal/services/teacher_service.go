package services

import (
	"school-exam-system/internal/models"

	"gorm.io/gorm"
)

type TeacherService struct {
	db *gorm.DB
}

func NewTeacherService(db *gorm.DB) *TeacherService {
	return &TeacherService{db: db}
}

func (s *TeacherService) GetTeacherByUserID(userID uint) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := s.db.Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (s *TeacherService) GetTeacherProfile(userID uint) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := s.db.Preload("User").Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}