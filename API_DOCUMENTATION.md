# School Exam System API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Semua endpoint yang dilindungi memerlukan header Authorization:
```
Authorization: Bearer <jwt_token>
```

## Endpoints yang Sudah Tersedia

### 1. Health Check
**GET** `/health`

Response:
```json
{
  "status": "ok",
  "message": "School Exam System is running"
}
```

### 2. Authentication

#### Login
**POST** `/api/auth/login`

Request Body:
```json
{
  "email": "admin@school.com",
  "password": "admin123"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-02T10:00:00Z",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@school.com",
      "role": "system_admin"
    }
  }
}
```

#### Logout
**POST** `/api/auth/logout`

Headers: `Authorization: Bearer <token>`

Response:
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### Get Profile
**GET** `/api/auth/profile`

Headers: `Authorization: Bearer <token>`

Response:
```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "username": "admin",
    "role": "system_admin"
  }
}
```

### 3. Admin Endpoints (Hanya untuk System Admin)

#### Create User
**POST** `/api/admin/users`

Headers: `Authorization: Bearer <admin_token>`

Request Body untuk Teacher:
```json
{
  "username": "teacher1",
  "email": "teacher1@school.com",
  "password": "password123",
  "role": "teacher",
  "full_name": "Budi Santoso",
  "phone": "081234567890",
  "address": "Jl. Pendidikan No. 1",
  "employee_id": "T001"
}
```

Request Body untuk Student:
```json
{
  "username": "student1",
  "email": "student1@school.com",
  "password": "password123",
  "role": "student",
  "full_name": "Siti Nurhaliza",
  "phone": "081234567891",
  "address": "Jl. Siswa No. 1",
  "student_id": "S001",
  "class_id": 1
}
```

#### Get All Users
**GET** `/api/admin/users`

Headers: `Authorization: Bearer <admin_token>`

Query Parameters:
- `role` (optional): Filter by role (system_admin, teacher, student)

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@school.com",
      "role": "system_admin",
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

#### Get User by ID
**GET** `/api/admin/users/:id`

Headers: `Authorization: Bearer <admin_token>`

#### Update User
**PUT** `/api/admin/users/:id`

Headers: `Authorization: Bearer <admin_token>`

Request Body:
```json
{
  "username": "new_username",
  "email": "new_email@school.com",
  "full_name": "New Full Name",
  "phone": "081234567890",
  "address": "New Address",
  "is_active": true
}
```

#### Delete User
**DELETE** `/api/admin/users/:id`

Headers: `Authorization: Bearer <admin_token>`

### 4. Teacher Endpoints (Hanya untuk Teacher)

#### Get Classes (Coming Soon)
**GET** `/api/teacher/classes`

Headers: `Authorization: Bearer <teacher_token>`

### 5. Student Endpoints (Hanya untuk Student)

#### Get Profile (Coming Soon)
**GET** `/api/student/profile`

Headers: `Authorization: Bearer <student_token>`

## Default Users

Setelah menjalankan aplikasi, user default berikut akan tersedia:

1. **Admin**
   - Email: `admin@school.com`
   - Password: `admin123`
   - Role: `system_admin`

2. **Teacher** (Sample)
   - Email: `teacher1@school.com`
   - Password: `teacher123`
   - Role: `teacher`

3. **Student** (Sample)
   - Email: `student1@school.com`
   - Password: `student123`
   - Role: `student`

## Error Responses

### 400 Bad Request
```json
{
  "error": "Validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 403 Forbidden
```json
{
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "User not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Cara Menjalankan

1. Pastikan MySQL sudah berjalan
2. Buat database `ujian`
3. Jalankan aplikasi:
   ```bash
   go run cmd/server/main.go
   ```
4. Aplikasi akan berjalan di `http://localhost:8080`

## Testing dengan cURL

### Login sebagai Admin
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@school.com","password":"admin123"}'
```

### Get All Users (dengan token admin)
```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Create Teacher
```bash
curl -X POST http://localhost:8080/api/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "username": "teacher2",
    "email": "teacher2@school.com",
    "password": "password123",
    "role": "teacher",
    "full_name": "Ahmad Wijaya",
    "phone": "081234567892",
    "address": "Jl. Guru No. 2",
    "employee_id": "T002"
  }'
```
## 
Updated Endpoints (Teacher & Student)

### 👨‍🏫 Teacher Endpoints (Butuh token teacher)

#### Get Teacher Profile
**GET** `/api/teacher/profile`

Headers: `Authorization: Bearer <teacher_token>`

Response:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 2,
    "employee_id": "T001",
    "full_name": "Budi Santoso",
    "phone": "081234567890",
    "address": "Jl. Pendidikan No. 1",
    "user": {
      "id": 2,
      "username": "teacher1",
      "email": "teacher1@school.com",
      "role": "teacher"
    }
  }
}
```

#### Get My Classes
**GET** `/api/teacher/classes`

Headers: `Authorization: Bearer <teacher_token>`

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "class_id": 1,
      "subject_id": 1,
      "teacher_id": 1,
      "class": {
        "id": 1,
        "name": "X IPA 1",
        "grade": 10,
        "stream": "IPA",
        "section": "1"
      },
      "subject": {
        "id": 1,
        "name": "Matematika",
        "code": "MTK"
      }
    }
  ]
}
```

#### Get My Students
**GET** `/api/teacher/students`

Headers: `Authorization: Bearer <teacher_token>`

#### Get Students by Class
**GET** `/api/teacher/classes/:classId/students`

Headers: `Authorization: Bearer <teacher_token>`

### 👨‍🎓 Student Endpoints (Butuh token student)

#### Get Student Profile
**GET** `/api/student/profile`

Headers: `Authorization: Bearer <student_token>`

Response:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 3,
    "student_id": "S001",
    "full_name": "Siti Nurhaliza",
    "class_id": 1,
    "phone": "081234567891",
    "address": "Jl. Siswa No. 1",
    "user": {
      "id": 3,
      "username": "student1",
      "email": "student1@school.com",
      "role": "student"
    },
    "class": {
      "id": 1,
      "name": "X IPA 1",
      "grade": 10,
      "stream": "IPA",
      "section": "1"
    }
  }
}
```

#### Get My Subjects
**GET** `/api/student/subjects`

Headers: `Authorization: Bearer <student_token>`

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Matematika",
      "code": "MTK",
      "description": "Mata pelajaran Matematika"
    },
    {
      "id": 2,
      "name": "Bahasa Indonesia",
      "code": "BIND",
      "description": "Mata pelajaran Bahasa Indonesia"
    }
  ]
}
```

#### Get My Class
**GET** `/api/student/class`

Headers: `Authorization: Bearer <student_token>`

#### Get Classmates
**GET** `/api/student/classmates`

Headers: `Authorization: Bearer <student_token>`

## Testing New Endpoints

### Login sebagai Teacher
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher1@school.com","password":"teacher123"}'
```

### Get Teacher Profile
```bash
curl -X GET http://localhost:8080/api/teacher/profile \
  -H "Authorization: Bearer YOUR_TEACHER_TOKEN"
```

### Login sebagai Student
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"student1@school.com","password":"student123"}'
```

### Get Student Profile
```bash
curl -X GET http://localhost:8080/api/student/profile \
  -H "Authorization: Bearer YOUR_STUDENT_TOKEN"
```
#
# 📝 Exam System Endpoints (NEW!)

### 👨‍🏫 Teacher Exam Management

#### Create Exam
**POST** `/api/teacher/exams`

Headers: `Authorization: Bearer <teacher_token>`

Request Body:
```json
{
  "title": "Ujian Matematika Semester 1",
  "description": "Ujian tengah semester matematika",
  "subject_id": 1,
  "class_id": 1,
  "start_time": "2024-01-15T08:00:00Z",
  "end_time": "2024-01-15T10:00:00Z",
  "duration": 90,
  "max_attempts": 1
}
```

#### Get My Exams
**GET** `/api/teacher/exams`

#### Get Exam Details
**GET** `/api/teacher/exams/:id`

#### Update Exam
**PUT** `/api/teacher/exams/:id`

#### Delete Exam
**DELETE** `/api/teacher/exams/:id`

#### Publish Exam
**POST** `/api/teacher/exams/:id/publish`

#### Add Question to Exam
**POST** `/api/teacher/exams/:examId/questions`

Request Body:
```json
{
  "text": "Berapa hasil dari 2 + 2?",
  "type": "multiple_choice",
  "options": ["3", "4", "5", "6"],
  "answer": "4",
  "points": 10,
  "order_index": 1
}
```

#### Get Exam Questions
**GET** `/api/teacher/exams/:examId/questions`

#### Get Exam Submissions
**GET** `/api/teacher/exams/:examId/submissions`

### 👨‍🎓 Student Exam Taking

#### Get Available Exams
**GET** `/api/student/exams`

Headers: `Authorization: Bearer <student_token>`

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Ujian Matematika Semester 1",
      "description": "Ujian tengah semester matematika",
      "start_time": "2024-01-15T08:00:00Z",
      "end_time": "2024-01-15T10:00:00Z",
      "duration": 90,
      "status": "published",
      "subject": {
        "name": "Matematika",
        "code": "MTK"
      }
    }
  ]
}
```

#### Start Exam
**POST** `/api/student/exams/:id/start`

Response:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "exam_id": 1,
    "student_id": 1,
    "started_at": "2024-01-15T08:00:00Z",
    "status": "in_progress",
    "attempt_num": 1
  }
}
```

#### Submit Answer
**POST** `/api/student/submissions/:submissionId/answer`

Request Body:
```json
{
  "question_id": 1,
  "answer": "4"
}
```

#### Finish Exam
**POST** `/api/student/submissions/:submissionId/finish`

#### Get My Grades
**GET** `/api/student/grades`

## 📊 Question Types Supported

### 1. Multiple Choice
```json
{
  "type": "multiple_choice",
  "options": ["Option A", "Option B", "Option C", "Option D"],
  "answer": "Option B"
}
```

### 2. True/False
```json
{
  "type": "true_false",
  "answer": "true"
}
```

### 3. Short Answer
```json
{
  "type": "short_answer",
  "answer": "Expected answer"
}
```

### 4. Essay
```json
{
  "type": "essay",
  "answer": "Sample answer for reference"
}
```

## 🎯 Exam Workflow

### For Teachers:
1. **Create Exam** → `POST /api/teacher/exams`
2. **Add Questions** → `POST /api/teacher/exams/:examId/questions`
3. **Publish Exam** → `POST /api/teacher/exams/:id/publish`
4. **Monitor Submissions** → `GET /api/teacher/exams/:examId/submissions`
5. **Grade Essays** (if any) → Manual grading interface

### For Students:
1. **View Available Exams** → `GET /api/student/exams`
2. **Start Exam** → `POST /api/student/exams/:id/start`
3. **Submit Answers** → `POST /api/student/submissions/:submissionId/answer`
4. **Finish Exam** → `POST /api/student/submissions/:submissionId/finish`
5. **View Grades** → `GET /api/student/grades`

## 🔄 Exam Status Flow

1. **draft** → Exam created, can add/edit questions
2. **published** → Exam ready, students can see it
3. **active** → Exam is currently running
4. **closed** → Exam finished, no more submissions

## ⏰ Auto-Grading Features

- **Multiple Choice**: Automatically graded
- **True/False**: Automatically graded  
- **Short Answer**: Simple text matching
- **Essay**: Requires manual grading by teacher

## 🧪 Testing Exam System

### Create and Take an Exam:

1. **Login as Teacher:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher1@school.com","password":"teacher123"}'
```

2. **Create Exam:**
```bash
curl -X POST http://localhost:8080/api/teacher/exams \
  -H "Authorization: Bearer TEACHER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Exam",
    "subject_id": 1,
    "class_id": 1,
    "start_time": "2024-01-15T08:00:00Z",
    "end_time": "2024-01-15T10:00:00Z",
    "duration": 60
  }'
```

3. **Add Question:**
```bash
curl -X POST http://localhost:8080/api/teacher/exams/1/questions \
  -H "Authorization: Bearer TEACHER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "What is 2+2?",
    "type": "multiple_choice",
    "options": ["3", "4", "5"],
    "answer": "4",
    "points": 10
  }'
```

4. **Publish Exam:**
```bash
curl -X POST http://localhost:8080/api/teacher/exams/1/publish \
  -H "Authorization: Bearer TEACHER_TOKEN"
```

5. **Login as Student and Take Exam:**
```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"student1@school.com","password":"student123"}'

# Start exam
curl -X POST http://localhost:8080/api/student/exams/1/start \
  -H "Authorization: Bearer STUDENT_TOKEN"

# Submit answer
curl -X POST http://localhost:8080/api/student/submissions/1/answer \
  -H "Authorization: Bearer STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"question_id": 1, "answer": "4"}'

# Finish exam
curl -X POST http://localhost:8080/api/student/submissions/1/finish \
  -H "Authorization: Bearer STUDENT_TOKEN"
```## 🏫 A
cademic Management System (NEW!)

### 👨‍💼 Admin Academic Management

#### Create Class
**POST** `/api/admin/classes`

Headers: `Authorization: Bearer <admin_token>`

Request Body:
```json
{
  "name": "X IPA 3",
  "grade": 10,
  "stream": "IPA",
  "section": "3",
  "description": "Kelas 10 IPA 3"
}
```

#### Get All Classes
**GET** `/api/admin/classes`

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "X IPA 1",
      "grade": 10,
      "stream": "IPA",
      "section": "1",
      "description": "Kelas 10 IPA 1",
      "is_active": true
    }
  ]
}
```

#### Create Subject
**POST** `/api/admin/subjects`

Request Body:
```json
{
  "name": "Fisika",
  "code": "FIS",
  "description": "Mata pelajaran Fisika"
}
```

#### Get All Subjects
**GET** `/api/admin/subjects`

#### Assign Teacher to Class & Subject
**POST** `/api/admin/assign-teacher`

Request Body:
```json
{
  "teacher_id": 1,
  "class_id": 1,
  "subject_id": 1
}
```

#### Get Class Roster
**GET** `/api/admin/classes/:classId/roster`

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "student_id": "S001",
      "full_name": "Siti Nurhaliza",
      "phone": "081234567891",
      "user": {
        "username": "student1",
        "email": "student1@school.com"
      }
    }
  ]
}
```

#### Get Teacher Workload
**GET** `/api/admin/teachers/:teacherId/workload`

Response:
```json
{
  "success": true,
  "data": {
    "teacher_name": "Budi Santoso",
    "employee_id": "T001",
    "total_assignments": 3,
    "classes_taught": 2,
    "subjects_taught": 2,
    "students_taught": 45,
    "exams_created": 8,
    "assignments": [
      {
        "class": {"name": "X IPA 1"},
        "subject": {"name": "Matematika", "code": "MTK"}
      }
    ]
  }
}
```

#### Enroll Student to Class
**POST** `/api/admin/students/:studentId/enroll`

Request Body:
```json
{
  "class_id": 1
}
```

#### Transfer Student to Another Class
**POST** `/api/admin/students/:studentId/transfer`

Request Body:
```json
{
  "new_class_id": 2
}
```

## 📊 Academic Management Features

### 🎓 Student Enrollment System
- **Enroll students** to classes
- **Transfer students** between classes
- **View class rosters** with student details
- **Bulk enrollment** capabilities
- **Unassigned students** tracking

### 👨‍🏫 Teacher Assignment System
- **Assign teachers** to class-subject combinations
- **Teacher workload** monitoring
- **Access validation** for teacher permissions
- **Bulk assignment** capabilities
- **Assignment history** tracking

### 📈 Academic Analytics
- **Class statistics** (student count, subjects, exams)
- **Teacher workload** analysis
- **Subject coverage** tracking
- **Enrollment reports**

## 🔄 Academic Workflow

### Setting Up Academic Structure:
1. **Create Classes** → `POST /api/admin/classes`
2. **Create Subjects** → `POST /api/admin/subjects`
3. **Assign Teachers** → `POST /api/admin/assign-teacher`
4. **Enroll Students** → `POST /api/admin/students/:studentId/enroll`

### Managing Students:
1. **View Class Roster** → `GET /api/admin/classes/:classId/roster`
2. **Transfer Student** → `POST /api/admin/students/:studentId/transfer`
3. **Monitor Enrollment** → Various analytics endpoints

### Managing Teachers:
1. **View Workload** → `GET /api/admin/teachers/:teacherId/workload`
2. **Assign New Classes** → `POST /api/admin/assign-teacher`
3. **Monitor Performance** → Analytics and reports

## 🧪 Testing Academic Management

### Setup Complete Academic Structure:

1. **Login as Admin:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@school.com","password":"admin123"}'
```

2. **Create Class:**
```bash
curl -X POST http://localhost:8080/api/admin/classes \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "XI IPA 1",
    "grade": 11,
    "stream": "IPA",
    "section": "1",
    "description": "Kelas 11 IPA 1"
  }'
```

3. **Create Subject:**
```bash
curl -X POST http://localhost:8080/api/admin/subjects \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kimia",
    "code": "KIM",
    "description": "Mata pelajaran Kimia"
  }'
```

4. **Assign Teacher:**
```bash
curl -X POST http://localhost:8080/api/admin/assign-teacher \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "teacher_id": 1,
    "class_id": 1,
    "subject_id": 1
  }'
```

5. **Enroll Student:**
```bash
curl -X POST http://localhost:8080/api/admin/students/1/enroll \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"class_id": 1}'
```

6. **Check Class Roster:**
```bash
curl -X GET http://localhost:8080/api/admin/classes/1/roster \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## 📋 Updated Endpoint Summary

### Total Endpoints: **38 endpoints** (+8 new)

#### 👨‍💼 Admin (18 endpoints):
- **User Management**: 5 endpoints
- **Class Management**: 3 endpoints  
- **Subject Management**: 2 endpoints
- **Teacher Assignment**: 2 endpoints
- **Student Enrollment**: 2 endpoints
- **Analytics**: 4 endpoints

#### 👨‍🏫 Teacher (12 endpoints):
- **Profile & Classes**: 4 endpoints
- **Exam Management**: 6 endpoints
- **Question Management**: 2 endpoints

#### 👨‍🎓 Student (9 endpoints):
- **Profile & Academic Info**: 4 endpoints
- **Exam Taking**: 5 endpoints

#### 🔐 Auth & Health (3 endpoints):
- **Authentication**: 3 endpoints
- **Health Check**: 1 endpoint

Sistem sekarang sudah menjadi **platform akademik yang sangat lengkap** dengan manajemen kelas, mata pelajaran, penugasan guru, enrollment siswa, dan sistem ujian online yang terintegrasi!