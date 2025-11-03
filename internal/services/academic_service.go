package services

import (
	"errors"
	"school-exam-system/internal/models"

	"gorm.io/gorm"
)

type AcademicService struct {
	db *gorm.DB
}

type CreateClassRequest struct {
	Name        string `json:"name" binding:"required"`
	Grade       int    `json:"grade" binding:"required"`
	Stream      string `json:"stream" binding:"required"`
	Section     string `json:"section" binding:"required"`
	Description string `json:"description"`
}

type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

type AssignTeacherRequest struct {
	TeacherID uint `json:"teacher_id" binding:"required"`
	ClassID   uint `json:"class_id" binding:"required"`
	SubjectID uint `json:"subject_id" binding:"required"`
}

func NewAcademicService(db *gorm.DB) *AcademicService {
	return &AcademicService{db: db}
}

// Class Management
func (s *AcademicService) CreateClass(req CreateClassRequest) (*models.Class, error) {
	// Check if class name already exists
	var existingClass models.Class
	if err := s.db.Where("name = ?", req.Name).First(&existingClass).Error; err == nil {
		return nil, errors.New("class name already exists")
	}

	class := models.Class{
		Name:        req.Name,
		Grade:       req.Grade,
		Stream:      req.Stream,
		Section:     req.Section,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.db.Create(&class).Error; err != nil {
		return nil, err
	}

	return &class, nil
}

func (s *AcademicService) GetClassByID(id uint) (*models.Class, error) {
	var class models.Class
	if err := s.db.Preload("Students").First(&class, id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (s *AcademicService) ListClasses() ([]models.Class, error) {
	var classes []models.Class
	if err := s.db.Where("is_active = ?", true).Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (s *AcademicService) UpdateClass(id uint, req CreateClassRequest) (*models.Class, error) {
	var class models.Class
	if err := s.db.First(&class, id).Error; err != nil {
		return nil, err
	}

	class.Name = req.Name
	class.Grade = req.Grade
	class.Stream = req.Stream
	class.Section = req.Section
	class.Description = req.Description

	if err := s.db.Save(&class).Error; err != nil {
		return nil, err
	}

	return &class, nil
}

func (s *AcademicService) DeleteClass(id uint) error {
	return s.db.Model(&models.Class{}).Where("id = ?", id).Update("is_active", false).Error
}

// Subject Management
func (s *AcademicService) CreateSubject(req CreateSubjectRequest) (*models.Subject, error) {
	// Check if subject code already exists
	var existingSubject models.Subject
	if err := s.db.Where("code = ?", req.Code).First(&existingSubject).Error; err == nil {
		return nil, errors.New("subject code already exists")
	}

	subject := models.Subject{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.db.Create(&subject).Error; err != nil {
		return nil, err
	}

	return &subject, nil
}

func (s *AcademicService) ListSubjects() ([]models.Subject, error) {
	var subjects []models.Subject
	if err := s.db.Where("is_active = ?", true).Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (s *AcademicService) GetSubjectByID(id uint) (*models.Subject, error) {
	var subject models.Subject
	if err := s.db.First(&subject, id).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

// Teacher Assignment
func (s *AcademicService) AssignTeacherToClass(req AssignTeacherRequest) error {
	// Check if assignment already exists
	var existingAssignment models.ClassSubject
	if err := s.db.Where("class_id = ? AND subject_id = ?", req.ClassID, req.SubjectID).First(&existingAssignment).Error; err == nil {
		return errors.New("subject already assigned to this class")
	}

	assignment := models.ClassSubject{
		ClassID:   req.ClassID,
		SubjectID: req.SubjectID,
		TeacherID: req.TeacherID,
		IsActive:  true,
	}

	return s.db.Create(&assignment).Error
}

func (s *AcademicService) GetTeacherAssignments(teacherID uint) ([]models.ClassSubject, error) {
	var assignments []models.ClassSubject
	if err := s.db.Preload("Class").Preload("Subject").Where("teacher_id = ? AND is_active = ?", teacherID, true).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (s *AcademicService) GetStudentsByClass(classID uint) ([]models.Student, error) {
	var students []models.Student
	if err := s.db.Preload("User").Where("class_id = ?", classID).Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (s *AcademicService) GetStudentsByTeacher(teacherID uint) ([]models.Student, error) {
	var students []models.Student
	
	// Get classes taught by teacher
	var assignments []models.ClassSubject
	if err := s.db.Where("teacher_id = ? AND is_active = ?", teacherID, true).Find(&assignments).Error; err != nil {
		return nil, err
	}

	if len(assignments) == 0 {
		return students, nil
	}

	// Get unique class IDs
	classIDs := make([]uint, 0)
	classMap := make(map[uint]bool)
	for _, assignment := range assignments {
		if !classMap[assignment.ClassID] {
			classIDs = append(classIDs, assignment.ClassID)
			classMap[assignment.ClassID] = true
		}
	}

	// Get students from those classes
	if err := s.db.Preload("User").Preload("Class").Where("class_id IN ?", classIDs).Find(&students).Error; err != nil {
		return nil, err
	}

	return students, nil
}

// Student Enrollment System
func (s *AcademicService) EnrollStudent(studentID uint, classID uint) error {
	// Check if student exists
	var student models.Student
	if err := s.db.First(&student, studentID).Error; err != nil {
		return errors.New("student not found")
	}

	// Check if class exists
	var class models.Class
	if err := s.db.First(&class, classID).Error; err != nil {
		return errors.New("class not found")
	}

	// Check if student is already enrolled in another class
	if student.ClassID != 0 && student.ClassID != classID {
		return errors.New("student is already enrolled in another class")
	}

	// Update student's class
	student.ClassID = classID
	return s.db.Save(&student).Error
}

func (s *AcademicService) UnenrollStudent(studentID uint) error {
	var student models.Student
	if err := s.db.First(&student, studentID).Error; err != nil {
		return errors.New("student not found")
	}

	// Remove from class
	student.ClassID = 0
	return s.db.Save(&student).Error
}

func (s *AcademicService) TransferStudent(studentID uint, newClassID uint) error {
	// Check if student exists
	var student models.Student
	if err := s.db.First(&student, studentID).Error; err != nil {
		return errors.New("student not found")
	}

	// Check if new class exists
	var newClass models.Class
	if err := s.db.First(&newClass, newClassID).Error; err != nil {
		return errors.New("new class not found")
	}

	// Update student's class
	student.ClassID = newClassID
	return s.db.Save(&student).Error
}

func (s *AcademicService) GetClassRoster(classID uint) ([]models.Student, error) {
	var students []models.Student
	if err := s.db.Preload("User").Where("class_id = ?", classID).
		Order("full_name ASC").Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (s *AcademicService) GetUnassignedStudents() ([]models.Student, error) {
	var students []models.Student
	if err := s.db.Preload("User").Where("class_id = 0 OR class_id IS NULL").
		Order("full_name ASC").Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (s *AcademicService) BulkEnrollStudents(studentIDs []uint, classID uint) error {
	// Check if class exists
	var class models.Class
	if err := s.db.First(&class, classID).Error; err != nil {
		return errors.New("class not found")
	}

	// Update all students
	return s.db.Model(&models.Student{}).
		Where("id IN ?", studentIDs).
		Update("class_id", classID).Error
}

func (s *AcademicService) GetClassStatistics(classID uint) (map[string]interface{}, error) {
	var class models.Class
	if err := s.db.First(&class, classID).Error; err != nil {
		return nil, errors.New("class not found")
	}

	// Count students
	var studentCount int64
	s.db.Model(&models.Student{}).Where("class_id = ?", classID).Count(&studentCount)

	// Count subjects
	var subjectCount int64
	s.db.Model(&models.ClassSubject{}).Where("class_id = ? AND is_active = ?", classID, true).Count(&subjectCount)

	// Count exams
	var examCount int64
	s.db.Model(&models.Exam{}).Where("class_id = ?", classID).Count(&examCount)

	return map[string]interface{}{
		"class_name":     class.Name,
		"student_count":  studentCount,
		"subject_count":  subjectCount,
		"exam_count":     examCount,
		"grade":          class.Grade,
		"stream":         class.Stream,
	}, nil
}

// Teacher Assignment System
func (s *AcademicService) AssignTeacherToSubject(teacherID uint, classID uint, subjectID uint) error {
	// Check if teacher exists
	var teacher models.Teacher
	if err := s.db.First(&teacher, teacherID).Error; err != nil {
		return errors.New("teacher not found")
	}

	// Check if class exists
	var class models.Class
	if err := s.db.First(&class, classID).Error; err != nil {
		return errors.New("class not found")
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.First(&subject, subjectID).Error; err != nil {
		return errors.New("subject not found")
	}

	// Check if assignment already exists
	var existingAssignment models.ClassSubject
	if err := s.db.Where("class_id = ? AND subject_id = ?", classID, subjectID).First(&existingAssignment).Error; err == nil {
		// Update existing assignment
		existingAssignment.TeacherID = teacherID
		existingAssignment.IsActive = true
		return s.db.Save(&existingAssignment).Error
	}

	// Create new assignment
	assignment := models.ClassSubject{
		ClassID:   classID,
		SubjectID: subjectID,
		TeacherID: teacherID,
		IsActive:  true,
	}

	return s.db.Create(&assignment).Error
}

func (s *AcademicService) RemoveTeacherAssignment(assignmentID uint) error {
	var assignment models.ClassSubject
	if err := s.db.First(&assignment, assignmentID).Error; err != nil {
		return errors.New("assignment not found")
	}

	// Check if there are active exams for this assignment
	var examCount int64
	s.db.Model(&models.Exam{}).
		Where("teacher_id = ? AND class_id = ? AND subject_id = ? AND status != ?", 
			assignment.TeacherID, assignment.ClassID, assignment.SubjectID, models.ExamStatusClosed).
		Count(&examCount)

	if examCount > 0 {
		return errors.New("cannot remove assignment with active exams")
	}

	// Soft delete by setting inactive
	assignment.IsActive = false
	return s.db.Save(&assignment).Error
}

func (s *AcademicService) GetTeacherWorkload(teacherID uint) (map[string]interface{}, error) {
	var teacher models.Teacher
	if err := s.db.Preload("User").First(&teacher, teacherID).Error; err != nil {
		return nil, errors.New("teacher not found")
	}

	// Get assignments
	assignments, err := s.GetTeacherAssignments(teacherID)
	if err != nil {
		return nil, err
	}

	// Count classes and subjects
	classMap := make(map[uint]bool)
	subjectMap := make(map[uint]bool)
	
	for _, assignment := range assignments {
		classMap[assignment.ClassID] = true
		subjectMap[assignment.SubjectID] = true
	}

	// Count students taught
	students, err := s.GetStudentsByTeacher(teacherID)
	if err != nil {
		return nil, err
	}

	// Count exams created
	var examCount int64
	s.db.Model(&models.Exam{}).Where("teacher_id = ?", teacherID).Count(&examCount)

	return map[string]interface{}{
		"teacher_name":      teacher.FullName,
		"employee_id":       teacher.EmployeeID,
		"total_assignments": len(assignments),
		"classes_taught":    len(classMap),
		"subjects_taught":   len(subjectMap),
		"students_taught":   len(students),
		"exams_created":     examCount,
		"assignments":       assignments,
	}, nil
}

func (s *AcademicService) GetSubjectTeachers(subjectID uint) ([]models.Teacher, error) {
	var teachers []models.Teacher
	if err := s.db.Joins("JOIN class_subjects ON teachers.id = class_subjects.teacher_id").
		Where("class_subjects.subject_id = ? AND class_subjects.is_active = ?", subjectID, true).
		Preload("User").
		Distinct().Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

func (s *AcademicService) GetClassTeachers(classID uint) ([]models.Teacher, error) {
	var teachers []models.Teacher
	if err := s.db.Joins("JOIN class_subjects ON teachers.id = class_subjects.teacher_id").
		Where("class_subjects.class_id = ? AND class_subjects.is_active = ?", classID, true).
		Preload("User").
		Distinct().Find(&teachers).Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

func (s *AcademicService) ValidateTeacherAccess(teacherID uint, classID uint, subjectID uint) error {
	var assignment models.ClassSubject
	if err := s.db.Where("teacher_id = ? AND class_id = ? AND subject_id = ? AND is_active = ?", 
		teacherID, classID, subjectID, true).First(&assignment).Error; err != nil {
		return errors.New("teacher not assigned to this class and subject")
	}
	return nil
}

func (s *AcademicService) GetUnassignedSubjects(classID uint) ([]models.Subject, error) {
	var subjects []models.Subject
	
	// Get subjects that are not assigned to the class
	if err := s.db.Where("id NOT IN (SELECT subject_id FROM class_subjects WHERE class_id = ? AND is_active = ?)", 
		classID, true).Find(&subjects).Error; err != nil {
		return nil, err
	}
	
	return subjects, nil
}

func (s *AcademicService) GetTeacherSchedule(teacherID uint) ([]map[string]interface{}, error) {
	assignments, err := s.GetTeacherAssignments(teacherID)
	if err != nil {
		return nil, err
	}

	var schedule []map[string]interface{}
	
	for _, assignment := range assignments {
		// Get active exams for this assignment
		var exams []models.Exam
		s.db.Where("teacher_id = ? AND class_id = ? AND subject_id = ? AND status IN ?", 
			teacherID, assignment.ClassID, assignment.SubjectID, 
			[]models.ExamStatus{models.ExamStatusPublished, models.ExamStatusActive}).
			Find(&exams)

		scheduleItem := map[string]interface{}{
			"class_name":    assignment.Class.Name,
			"subject_name":  assignment.Subject.Name,
			"subject_code":  assignment.Subject.Code,
			"upcoming_exams": len(exams),
			"exams":         exams,
		}
		
		schedule = append(schedule, scheduleItem)
	}

	return schedule, nil
}

func (s *AcademicService) BulkAssignTeacher(teacherID uint, assignments []AssignTeacherRequest) error {
	// Check if teacher exists
	var teacher models.Teacher
	if err := s.db.First(&teacher, teacherID).Error; err != nil {
		return errors.New("teacher not found")
	}

	tx := s.db.Begin()
	
	for _, assignment := range assignments {
		// Validate class and subject exist
		var class models.Class
		if err := tx.First(&class, assignment.ClassID).Error; err != nil {
			tx.Rollback()
			return errors.New("class not found: " + string(rune(assignment.ClassID)))
		}

		var subject models.Subject
		if err := tx.First(&subject, assignment.SubjectID).Error; err != nil {
			tx.Rollback()
			return errors.New("subject not found: " + string(rune(assignment.SubjectID)))
		}

		// Check if assignment already exists
		var existingAssignment models.ClassSubject
		if err := tx.Where("class_id = ? AND subject_id = ?", assignment.ClassID, assignment.SubjectID).
			First(&existingAssignment).Error; err == nil {
			// Update existing
			existingAssignment.TeacherID = teacherID
			existingAssignment.IsActive = true
			if err := tx.Save(&existingAssignment).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Create new
			newAssignment := models.ClassSubject{
				ClassID:   assignment.ClassID,
				SubjectID: assignment.SubjectID,
				TeacherID: teacherID,
				IsActive:  true,
			}
			if err := tx.Create(&newAssignment).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}