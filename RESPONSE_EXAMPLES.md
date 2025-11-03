# API Response Examples

## New Standardized Response Format

All API responses now follow a consistent format:

### Success Response Format
```json
{
  "status": 200,
  "message": "Operation successful",
  "data": {
    // Response data here
  }
}
```

### Error Response Format
```json
{
  "status": 400,
  "message": "Bad Request",
  "error": "Detailed error message"
}
```

## Authentication Responses

### Login Success
**POST** `/api/auth/login`
```json
{
  "status": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-02T10:00:00Z"
  }
}
```

### Login Error
```json
{
  "status": 401,
  "message": "Unauthorized",
  "error": "Invalid email or password"
}
```

### Profile Success
**GET** `/api/auth/profile`
```json
{
  "status": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "user_id": 1,
    "username": "admin",
    "role": "system_admin"
  }
}
```

### Logout Success
**POST** `/api/auth/logout`
```json
{
  "status": 200,
  "message": "Logged out successfully"
}
```

## Admin Responses

### Create User Success
**POST** `/api/admin/users`
```json
{
  "status": 201,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "username": "teacher1",
    "email": "teacher1@school.com",
    "role": "teacher",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### Get Users Success
**GET** `/api/admin/users`
```json
{
  "status": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@school.com",
      "role": "system_admin",
      "is_active": true
    },
    {
      "id": 2,
      "username": "teacher1",
      "email": "teacher1@school.com",
      "role": "teacher",
      "is_active": true
    }
  ]
}
```

### Create Class Success
**POST** `/api/admin/classes`
```json
{
  "status": 201,
  "message": "Class created successfully",
  "data": {
    "id": 1,
    "name": "X IPA 1",
    "grade": 10,
    "stream": "IPA",
    "section": "1",
    "description": "Kelas 10 IPA 1",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### Get Classes Success
**GET** `/api/admin/classes`
```json
{
  "status": 200,
  "message": "Classes retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "X IPA 1",
      "grade": 10,
      "stream": "IPA",
      "section": "1",
      "is_active": true
    }
  ]
}
```

### Create Subject Success
**POST** `/api/admin/subjects`
```json
{
  "status": 201,
  "message": "Subject created successfully",
  "data": {
    "id": 1,
    "name": "Matematika",
    "code": "MTK",
    "description": "Mata pelajaran Matematika",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

### Assign Teacher Success
**POST** `/api/admin/assign-teacher`
```json
{
  "status": 200,
  "message": "Teacher assigned successfully"
}
```

### Enroll Student Success
**POST** `/api/admin/students/1/enroll`
```json
{
  "status": 200,
  "message": "Student enrolled successfully"
}
```

## Error Responses

### Bad Request (400)
```json
{
  "status": 400,
  "message": "Bad Request",
  "error": "Invalid input data"
}
```

### Unauthorized (401)
```json
{
  "status": 401,
  "message": "Unauthorized",
  "error": "Authorization header required"
}
```

### Forbidden (403)
```json
{
  "status": 403,
  "message": "Forbidden",
  "error": "Insufficient permissions"
}
```

### Not Found (404)
```json
{
  "status": 404,
  "message": "Not Found",
  "error": "User not found"
}
```

### Internal Server Error (500)
```json
{
  "status": 500,
  "message": "Internal Server Error",
  "error": "Database connection failed"
}
```

## Teacher Responses

### Get Teacher Profile Success
**GET** `/api/teacher/profile`
```json
{
  "status": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "employee_id": "T001",
    "full_name": "Budi Santoso",
    "phone": "081234567890",
    "user": {
      "username": "teacher1",
      "email": "teacher1@school.com",
      "role": "teacher"
    }
  }
}
```

### Get My Classes Success
**GET** `/api/teacher/classes`
```json
{
  "status": 200,
  "message": "Classes retrieved successfully",
  "data": [
    {
      "id": 1,
      "class_id": 1,
      "subject_id": 1,
      "class": {
        "name": "X IPA 1",
        "grade": 10,
        "stream": "IPA"
      },
      "subject": {
        "name": "Matematika",
        "code": "MTK"
      }
    }
  ]
}
```

## Student Responses

### Get Student Profile Success
**GET** `/api/student/profile`
```json
{
  "status": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "student_id": "S001",
    "full_name": "Siti Nurhaliza",
    "phone": "081234567891",
    "user": {
      "username": "student1",
      "email": "student1@school.com",
      "role": "student"
    },
    "class": {
      "name": "X IPA 1",
      "grade": 10,
      "stream": "IPA"
    }
  }
}
```

### Get My Subjects Success
**GET** `/api/student/subjects`
```json
{
  "status": 200,
  "message": "Subjects retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Matematika",
      "code": "MTK",
      "description": "Mata pelajaran Matematika"
    }
  ]
}
```

## Key Improvements

1. **Consistent Structure**: All responses follow the same format
2. **Clear Status Codes**: HTTP status codes in response body
3. **Descriptive Messages**: Human-readable success/error messages
4. **Proper Error Handling**: Detailed error information
5. **Clean Data**: Well-structured response data
6. **No Extra Fields**: Removed unnecessary "success" boolean field