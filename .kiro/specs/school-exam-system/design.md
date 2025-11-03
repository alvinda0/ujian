# School Exam System Design Document

## Overview

The School Exam System is a web-based application built with Golang that provides role-based access to academic examination and assessment functionality. The system follows a clean architecture pattern with clear separation of concerns, implementing RESTful APIs for frontend communication and using PostgreSQL for data persistence. The application supports three distinct user roles (System Administrator, Teacher, Student) with granular permission controls and secure authentication mechanisms.

## Architecture

### High-Level Architecture

The system follows a layered architecture pattern:

```
┌─────────────────────────────────────────┐
│              Frontend Layer             │
│         (Web UI - HTML/CSS/JS)          │
└─────────────────────────────────────────┘
                     │ HTTP/REST
┌─────────────────────────────────────────┐
│             API Gateway Layer           │
│        (Gin HTTP Router + Auth)         │
└─────────────────────────────────────────┘
                     │
┌─────────────────────────────────────────┐
│            Business Logic Layer         │
│           (Services + Use Cases)        │
└─────────────────────────────────────────┘
                     │
┌─────────────────────────────────────────┐
│            Data Access Layer            │
│        (Repositories + GORM ORM)        │
└─────────────────────────────────────────┘
                     │
┌─────────────────────────────────────────┐
│             Database Layer              │
│              (PostgreSQL)               │
└─────────────────────────────────────────┘
```

### Technology Stack

- **Backend Framework**: Gin (Golang HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens with bcrypt password hashing
- **Frontend**: Server-side rendered HTML templates with vanilla JavaScript
- **Session Management**: Redis for session storage
- **Configuration**: Viper for environment configuration
- **Logging**: Logrus for structured logging

## Components and Interfaces

### Core Domain Models

#### User Management
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"unique;not null"`
    Email     string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Role      UserRole  `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserRole string
const (
    RoleSystemAdmin UserRole = "system_admin"
    RoleTeacher     UserRole = "teacher"
    RoleStudent     UserRole = "student"
)
```

#### Academic Structure
```go
type Class struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"` // e.g., "X IPA 1"
    Grade       int       `gorm:"not null"` // 10, 11, 12
    Stream      string    `gorm:"not null"` // IPA, IPS, etc.
    Students    []Student `gorm:"foreignKey:ClassID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Subject struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`
    Code        string    `gorm:"unique;not null"`
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ClassSubject struct {
    ID        uint    `gorm:"primaryKey"`
    ClassID   uint    `gorm:"not null"`
    SubjectID uint    `gorm:"not null"`
    TeacherID uint    `gorm:"not null"`
    Class     Class   `gorm:"foreignKey:ClassID"`
    Subject   Subject `gorm:"foreignKey:SubjectID"`
    Teacher   Teacher `gorm:"foreignKey:TeacherID"`
}
```

#### User Profiles
```go
type Teacher struct {
    ID          uint           `gorm:"primaryKey"`
    UserID      uint           `gorm:"unique;not null"`
    EmployeeID  string         `gorm:"unique;not null"`
    FullName    string         `gorm:"not null"`
    User        User           `gorm:"foreignKey:UserID"`
    Assignments []ClassSubject `gorm:"foreignKey:TeacherID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Student struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"unique;not null"`
    StudentID string    `gorm:"unique;not null"` // NIS
    FullName  string    `gorm:"not null"`
    ClassID   uint      `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    Class     Class     `gorm:"foreignKey:ClassID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Examination System
```go
type Exam struct {
    ID          uint           `gorm:"primaryKey"`
    Title       string         `gorm:"not null"`
    Description string
    SubjectID   uint           `gorm:"not null"`
    ClassID     uint           `gorm:"not null"`
    TeacherID   uint           `gorm:"not null"`
    StartTime   time.Time      `gorm:"not null"`
    EndTime     time.Time      `gorm:"not null"`
    Duration    int            `gorm:"not null"` // minutes
    Status      ExamStatus     `gorm:"default:'draft'"`
    Questions   []Question     `gorm:"foreignKey:ExamID"`
    Submissions []ExamSubmission `gorm:"foreignKey:ExamID"`
    Subject     Subject        `gorm:"foreignKey:SubjectID"`
    Class       Class          `gorm:"foreignKey:ClassID"`
    Teacher     Teacher        `gorm:"foreignKey:TeacherID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Question struct {
    ID         uint         `gorm:"primaryKey"`
    ExamID     uint         `gorm:"not null"`
    Text       string       `gorm:"not null"`
    Type       QuestionType `gorm:"not null"`
    Options    string       `gorm:"type:json"` // JSON array for multiple choice
    Answer     string       `gorm:"not null"`
    Points     int          `gorm:"default:1"`
    OrderIndex int          `gorm:"not null"`
    Exam       Exam         `gorm:"foreignKey:ExamID"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type ExamSubmission struct {
    ID          uint                `gorm:"primaryKey"`
    ExamID      uint                `gorm:"not null"`
    StudentID   uint                `gorm:"not null"`
    StartedAt   time.Time           `gorm:"not null"`
    SubmittedAt *time.Time
    Score       *float64
    Status      SubmissionStatus    `gorm:"default:'in_progress'"`
    Answers     []SubmissionAnswer  `gorm:"foreignKey:SubmissionID"`
    Exam        Exam                `gorm:"foreignKey:ExamID"`
    Student     Student             `gorm:"foreignKey:StudentID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Service Layer Interfaces

#### Authentication Service
```go
type AuthService interface {
    Login(username, password string) (*AuthResponse, error)
    Logout(token string) error
    ValidateToken(token string) (*UserClaims, error)
    RefreshToken(token string) (*AuthResponse, error)
}

type AuthResponse struct {
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    User      UserInfo  `json:"user"`
}
```

#### User Management Service
```go
type UserService interface {
    CreateUser(req CreateUserRequest) (*User, error)
    GetUserByID(id uint) (*User, error)
    UpdateUser(id uint, req UpdateUserRequest) (*User, error)
    DeleteUser(id uint) error
    ListUsers(filters UserFilters) ([]User, error)
}
```

#### Academic Management Service
```go
type AcademicService interface {
    // Class Management
    CreateClass(req CreateClassRequest) (*Class, error)
    GetClassByID(id uint) (*Class, error)
    ListClasses() ([]Class, error)
    
    // Subject Management
    CreateSubject(req CreateSubjectRequest) (*Subject, error)
    AssignTeacherToClass(teacherID, classID, subjectID uint) error
    
    // Student Management
    EnrollStudent(studentID, classID uint) error
    GetStudentsByClass(classID uint) ([]Student, error)
}
```

#### Exam Management Service
```go
type ExamService interface {
    CreateExam(req CreateExamRequest) (*Exam, error)
    GetExamByID(id uint) (*Exam, error)
    UpdateExam(id uint, req UpdateExamRequest) (*Exam, error)
    DeleteExam(id uint) error
    ListExamsByTeacher(teacherID uint) ([]Exam, error)
    ListExamsByStudent(studentID uint) ([]Exam, error)
    
    // Question Management
    AddQuestion(examID uint, req AddQuestionRequest) (*Question, error)
    UpdateQuestion(id uint, req UpdateQuestionRequest) (*Question, error)
    DeleteQuestion(id uint) error
    
    // Submission Management
    StartExam(examID, studentID uint) (*ExamSubmission, error)
    SubmitAnswer(submissionID, questionID uint, answer string) error
    FinishExam(submissionID uint) (*ExamSubmission, error)
    GradeSubmission(submissionID uint, score float64) error
}
```

### API Endpoints Structure

#### Authentication Endpoints
- `POST /api/auth/login` - User authentication
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Token refresh

#### System Admin Endpoints
- `GET /api/admin/users` - List all users
- `POST /api/admin/users` - Create new user
- `PUT /api/admin/users/:id` - Update user
- `DELETE /api/admin/users/:id` - Delete user
- `GET /api/admin/classes` - List all classes
- `POST /api/admin/classes` - Create class
- `GET /api/admin/reports/grades` - Export grade reports

#### Teacher Endpoints
- `GET /api/teacher/classes` - Get assigned classes
- `GET /api/teacher/students/:classId` - Get students in class
- `GET /api/teacher/exams` - List teacher's exams
- `POST /api/teacher/exams` - Create new exam
- `PUT /api/teacher/exams/:id` - Update exam
- `POST /api/teacher/exams/:id/questions` - Add question to exam
- `GET /api/teacher/submissions/:examId` - Get exam submissions
- `PUT /api/teacher/submissions/:id/grade` - Grade submission

#### Student Endpoints
- `GET /api/student/profile` - Get student profile
- `GET /api/student/subjects` - Get enrolled subjects
- `GET /api/student/exams` - Get available exams
- `POST /api/student/exams/:id/start` - Start exam
- `POST /api/student/submissions/:id/answer` - Submit answer
- `POST /api/student/submissions/:id/finish` - Finish exam
- `GET /api/student/grades` - Get personal grades

## Data Models

### Database Schema Design

The system uses PostgreSQL with the following key relationships:

1. **Users** have roles and are linked to either Teacher or Student profiles
2. **Classes** contain multiple Students and have multiple Subjects
3. **ClassSubject** junction table links Classes, Subjects, and Teachers
4. **Exams** belong to a specific Subject, Class, and Teacher
5. **Questions** belong to Exams and define the assessment content
6. **ExamSubmissions** track student attempts and scores
7. **SubmissionAnswers** store individual question responses

### Key Constraints and Indexes

```sql
-- Unique constraints
ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);
ALTER TABLE teachers ADD CONSTRAINT unique_employee_id UNIQUE (employee_id);
ALTER TABLE students ADD CONSTRAINT unique_student_id UNIQUE (student_id);

-- Foreign key constraints with cascading
ALTER TABLE class_subjects ADD CONSTRAINT fk_class_subject_teacher 
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE;

-- Indexes for performance
CREATE INDEX idx_exams_class_subject ON exams(class_id, subject_id);
CREATE INDEX idx_submissions_student_exam ON exam_submissions(student_id, exam_id);
CREATE INDEX idx_users_role ON users(role);
```

## Error Handling

### Error Response Structure
```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### Error Categories
- **Authentication Errors**: Invalid credentials, expired tokens
- **Authorization Errors**: Insufficient permissions, role restrictions
- **Validation Errors**: Invalid input data, constraint violations
- **Business Logic Errors**: Exam scheduling conflicts, duplicate enrollments
- **System Errors**: Database connection issues, external service failures

### Error Handling Middleware
```go
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *ValidationError:
                c.JSON(400, ErrorResponse{Code: "VALIDATION_ERROR", Message: e.Message})
            case *AuthorizationError:
                c.JSON(403, ErrorResponse{Code: "FORBIDDEN", Message: e.Message})
            default:
                c.JSON(500, ErrorResponse{Code: "INTERNAL_ERROR", Message: "Internal server error"})
            }
        }
    }
}
```

## Testing Strategy

### Unit Testing
- **Service Layer Testing**: Mock repositories and test business logic
- **Repository Testing**: Test database operations with test database
- **Utility Function Testing**: Test helper functions and validators
- **Coverage Target**: Minimum 80% code coverage for critical paths

### Integration Testing
- **API Endpoint Testing**: Test complete request/response cycles
- **Database Integration**: Test with real PostgreSQL test database
- **Authentication Flow**: Test JWT token generation and validation
- **Role-based Access**: Verify permission enforcement

### End-to-End Testing
- **User Journey Testing**: Complete workflows for each user role
- **Exam Flow Testing**: Full exam creation, taking, and grading process
- **Data Consistency**: Verify data integrity across operations

### Testing Tools
- **Unit Tests**: Go's built-in testing package with testify assertions
- **HTTP Testing**: httptest package for API testing
- **Database Testing**: dockertest for PostgreSQL test containers
- **Mocking**: gomock for interface mocking

### Test Data Management
```go
type TestDataBuilder struct {
    db *gorm.DB
}

func (t *TestDataBuilder) CreateTestUser(role UserRole) *User {
    // Create test user with specified role
}

func (t *TestDataBuilder) CreateTestClass() *Class {
    // Create test class with students
}

func (t *TestDataBuilder) CreateTestExam(teacherID, classID uint) *Exam {
    // Create test exam with questions
}
```

The testing strategy ensures reliability and maintainability while supporting continuous integration and deployment practices.