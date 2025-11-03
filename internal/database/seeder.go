package database

import (
	"school-exam-system/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	// Check if admin user already exists
	var count int64
	db.Model(&models.User{}).Where("role = ?", models.RoleSystemAdmin).Count(&count)
	
	if count > 0 {
		return nil // Admin already exists
	}

	// Create default admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Username: "admin",
		Email:    "admin@school.com",
		Password: string(hashedPassword),
		Role:     models.RoleSystemAdmin,
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	// Seed sample subjects
	subjects := []models.Subject{
		{Name: "Matematika", Code: "MTK", Description: "Mata pelajaran Matematika"},
		{Name: "Bahasa Indonesia", Code: "BIND", Description: "Mata pelajaran Bahasa Indonesia"},
		{Name: "Bahasa Inggris", Code: "BING", Description: "Mata pelajaran Bahasa Inggris"},
		{Name: "Fisika", Code: "FIS", Description: "Mata pelajaran Fisika"},
		{Name: "Kimia", Code: "KIM", Description: "Mata pelajaran Kimia"},
		{Name: "Biologi", Code: "BIO", Description: "Mata pelajaran Biologi"},
		{Name: "Sejarah", Code: "SEJ", Description: "Mata pelajaran Sejarah"},
		{Name: "Geografi", Code: "GEO", Description: "Mata pelajaran Geografi"},
		{Name: "Ekonomi", Code: "EKO", Description: "Mata pelajaran Ekonomi"},
		{Name: "Sosiologi", Code: "SOS", Description: "Mata pelajaran Sosiologi"},
	}

	for _, subject := range subjects {
		var existingSubject models.Subject
		if err := db.Where("code = ?", subject.Code).First(&existingSubject).Error; err != nil {
			if err := db.Create(&subject).Error; err != nil {
				return err
			}
		}
	}

	// Seed sample classes
	classes := []models.Class{
		{Name: "X IPA 1", Grade: 10, Stream: "IPA", Section: "1"},
		{Name: "X IPA 2", Grade: 10, Stream: "IPA", Section: "2"},
		{Name: "X IPS 1", Grade: 10, Stream: "IPS", Section: "1"},
		{Name: "X IPS 2", Grade: 10, Stream: "IPS", Section: "2"},
		{Name: "XI IPA 1", Grade: 11, Stream: "IPA", Section: "1"},
		{Name: "XI IPA 2", Grade: 11, Stream: "IPA", Section: "2"},
		{Name: "XI IPS 1", Grade: 11, Stream: "IPS", Section: "1"},
		{Name: "XI IPS 2", Grade: 11, Stream: "IPS", Section: "2"},
		{Name: "XII IPA 1", Grade: 12, Stream: "IPA", Section: "1"},
		{Name: "XII IPA 2", Grade: 12, Stream: "IPA", Section: "2"},
		{Name: "XII IPS 1", Grade: 12, Stream: "IPS", Section: "1"},
		{Name: "XII IPS 2", Grade: 12, Stream: "IPS", Section: "2"},
	}

	for _, class := range classes {
		var existingClass models.Class
		if err := db.Where("name = ?", class.Name).First(&existingClass).Error; err != nil {
			if err := db.Create(&class).Error; err != nil {
				return err
			}
		}
	}

	return nil
}