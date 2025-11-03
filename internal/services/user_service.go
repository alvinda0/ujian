package services

import (
	"errors"
	"school-exam-system/internal/models"
	"school-exam-system/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

type CreateUserRequest struct {
	Username string          `json:"username" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Role     models.UserRole `json:"role" binding:"required"`
	FullName string          `json:"full_name" binding:"required"`
	Phone    string          `json:"phone"`
	Address  string          `json:"address"`
	// For teacher
	EmployeeID string `json:"employee_id"`
	// For student
	StudentID string `json:"student_id"`
	ClassID   *uint  `json:"class_id"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	IsActive *bool  `json:"is_active"`
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(req CreateUserRequest) (*models.User, error) {
	// Check if username or email already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("username or email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	tx := s.db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create profile based on role
	switch req.Role {
	case models.RoleTeacher:
		if req.EmployeeID == "" {
			tx.Rollback()
			return nil, errors.New("employee_id is required for teacher")
		}
		teacher := models.Teacher{
			UserID:     user.ID,
			EmployeeID: req.EmployeeID,
			FullName:   req.FullName,
			Phone:      req.Phone,
			Address:    req.Address,
		}
		if err := tx.Create(&teacher).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	case models.RoleStudent:
		if req.StudentID == "" || req.ClassID == nil {
			tx.Rollback()
			return nil, errors.New("student_id and class_id are required for student")
		}
		student := models.Student{
			UserID:    user.ID,
			StudentID: req.StudentID,
			FullName:  req.FullName,
			ClassID:   *req.ClassID,
			Phone:     req.Phone,
			Address:   req.Address,
		}
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return &user, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) ListUsers(role string) ([]models.User, error) {
	var users []models.User
	query := s.db.Model(&models.User{})
	
	if role != "" {
		query = query.Where("role = ?", role)
	}
	
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(id uint, req UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Update user fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// Update profile based on role
	switch user.Role {
	case models.RoleTeacher:
		var teacher models.Teacher
		if err := s.db.Where("user_id = ?", user.ID).First(&teacher).Error; err == nil {
			if req.FullName != "" {
				teacher.FullName = req.FullName
			}
			if req.Phone != "" {
				teacher.Phone = req.Phone
			}
			if req.Address != "" {
				teacher.Address = req.Address
			}
			s.db.Save(&teacher)
		}

	case models.RoleStudent:
		var student models.Student
		if err := s.db.Where("user_id = ?", user.ID).First(&student).Error; err == nil {
			if req.FullName != "" {
				student.FullName = req.FullName
			}
			if req.Phone != "" {
				student.Phone = req.Phone
			}
			if req.Address != "" {
				student.Address = req.Address
			}
			s.db.Save(&student)
		}
	}

	return &user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}