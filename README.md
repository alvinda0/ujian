# School Exam System

A comprehensive web-based application for managing school examinations and academic assessments built with Go.

## Features

- **Role-based Access Control**: System Admin, Teacher, and Student roles
- **Class Management**: Organize students into classes (e.g., X IPA 1, XI IPS 2)
- **Subject Management**: Manage academic subjects and teacher assignments
- **Exam System**: Create, manage, and take online exams
- **Grading System**: Automatic and manual grading capabilities
- **Reporting**: Generate grade reports and export data

## Technology Stack

- **Backend**: Go with Gin framework
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT tokens
- **Frontend**: HTML templates with vanilla JavaScript

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create MySQL database:
   ```sql
   CREATE DATABASE ujian CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

4. Configure database connection in `configs/config.yaml`

5. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

## Configuration

Edit `configs/config.yaml` to configure:
- Server port and mode
- Database connection
- JWT settings

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout

### System Admin
- `GET /api/admin/users` - List all users
- `POST /api/admin/users` - Create new user
- `PUT /api/admin/users/:id` - Update user
- `DELETE /api/admin/users/:id` - Delete user

### Teacher
- `GET /api/teacher/classes` - Get assigned classes
- `GET /api/teacher/exams` - List teacher's exams
- `POST /api/teacher/exams` - Create new exam

### Student
- `GET /api/student/profile` - Get student profile
- `GET /api/student/exams` - Get available exams
- `GET /api/student/grades` - Get personal grades

## Project Structure

```
school-exam-system/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── models/          # Database models
│   ├── services/        # Business logic
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   └── utils/           # Utility functions
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS, JS, images
├── configs/             # Configuration files
└── README.md
```

## License

This project is licensed under the MIT License.