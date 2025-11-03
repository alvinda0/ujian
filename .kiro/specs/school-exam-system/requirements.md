# Requirements Document

## Introduction

The School Exam System is a comprehensive web-based application that manages academic examinations, student assessments, and educational data for schools. The system supports role-based access control with three primary user types: System Administrators, Teachers, and Students. It manages classes, subjects, exams, grades, and user accounts while providing secure access to educational resources based on user roles and class assignments.

## Glossary

- **School_Exam_System**: The complete web application for managing school examinations and academic data
- **System_Administrator**: A user with full system access who manages all data and configurations
- **Teacher**: An educator user who manages specific classes and subjects assigned to them
- **Student**: A learner user who accesses their own academic information and takes exams
- **Class**: An academic group (e.g., X IPA 1, XI IPS 2) containing students and assigned subjects
- **Subject**: An academic course or mata pelajaran (e.g., Mathematics, Indonesian Language)
- **Exam**: An assessment or ujian created by teachers for specific subjects and classes
- **Grade**: A numerical or letter score assigned to student exam performance
- **User_Account**: Authentication credentials and profile information for system access
- **Class_Assignment**: The relationship between teachers, subjects, and classes they are responsible for
- **Exam_Schedule**: Time-based scheduling for when exams are available to students
- **Academic_Data**: All educational information including grades, exam results, and student records

## Requirements

### Requirement 1

**User Story:** As a System Administrator, I want to manage all user accounts and academic data, so that I can maintain complete control over the school's examination system.

#### Acceptance Criteria

1. THE School_Exam_System SHALL provide System_Administrator access to view all Teacher accounts, Student accounts, Subject data, Exam data, Grade data, and Class data
2. WHEN System_Administrator creates new accounts, THE School_Exam_System SHALL store Teacher accounts and Student accounts with complete profile information
3. THE School_Exam_System SHALL allow System_Administrator to modify existing Teacher accounts, Student accounts, Subject data, and Class data
4. THE School_Exam_System SHALL enable System_Administrator to remove Teacher accounts, Student accounts, Subject data, and Class data from the system
5. THE School_Exam_System SHALL provide System_Administrator with export functionality for Grade data and Exam results

### Requirement 2

**User Story:** As a System Administrator, I want to configure class structures and teacher assignments, so that I can organize the academic structure effectively.

#### Acceptance Criteria

1. THE School_Exam_System SHALL allow System_Administrator to create Class entities with names like "X IPA 1" and "XI IPS 2"
2. THE School_Exam_System SHALL enable System_Administrator to assign Teacher accounts to specific Subject and Class combinations
3. WHEN System_Administrator creates Class_Assignment, THE School_Exam_System SHALL establish relationships between Teacher, Subject, and Class entities
4. THE School_Exam_System SHALL provide System_Administrator with comprehensive Grade reports organized by Class and Subject
5. THE School_Exam_System SHALL allow System_Administrator to configure global Exam_Schedule settings

### Requirement 3

**User Story:** As a Teacher, I want to manage my assigned classes and subjects, so that I can create and evaluate exams for my students.

#### Acceptance Criteria

1. THE School_Exam_System SHALL display only Class and Subject data that match Teacher's Class_Assignment
2. WHEN Teacher accesses student data, THE School_Exam_System SHALL show only Student accounts enrolled in Teacher's assigned Class entities
3. THE School_Exam_System SHALL allow Teacher to create, modify, and delete Exam entities for their assigned Subject and Class combinations
4. THE School_Exam_System SHALL enable Teacher to create and edit exam questions within their Exam entities
5. THE School_Exam_System SHALL provide Teacher with Grade assignment capabilities for Student exam submissions in their classes

### Requirement 4

**User Story:** As a Teacher, I want to be restricted to my assigned subjects and classes, so that data security and academic boundaries are maintained.

#### Acceptance Criteria

1. THE School_Exam_System SHALL prevent Teacher access to Class entities not included in their Class_Assignment
2. THE School_Exam_System SHALL prevent Teacher access to Subject data not included in their Class_Assignment
3. WHEN Teacher attempts to access unauthorized data, THE School_Exam_System SHALL deny access and display appropriate error messages
4. THE School_Exam_System SHALL filter all Teacher interface elements to show only authorized Class and Subject information
5. THE School_Exam_System SHALL restrict Teacher Grade viewing to only Student accounts in their assigned Class entities

### Requirement 5

**User Story:** As a Student, I want to access my academic information and take exams, so that I can participate in the school's assessment system.

#### Acceptance Criteria

1. THE School_Exam_System SHALL display Student's personal profile information including name, Class assignment, and student identification number
2. THE School_Exam_System SHALL show Student only Subject data associated with their assigned Class
3. WHEN Student accesses Grade information, THE School_Exam_System SHALL display only their personal Grade data and Exam results
4. THE School_Exam_System SHALL provide Student with access to scheduled Exam entities for their Class and Subject combinations
5. THE School_Exam_System SHALL enable Student to submit Exam responses during scheduled Exam periods

### Requirement 6

**User Story:** As a Student, I want to be restricted to my own academic data, so that privacy and data security are maintained.

#### Acceptance Criteria

1. THE School_Exam_System SHALL prevent Student access to other Student accounts' Grade data and personal information
2. THE School_Exam_System SHALL prevent Student modification of Class data, Subject data, and their own profile information
3. WHEN Student attempts to access unauthorized data, THE School_Exam_System SHALL deny access and redirect to authorized content
4. THE School_Exam_System SHALL filter all Student interface elements to show only their personal Academic_Data
5. THE School_Exam_System SHALL restrict Student Exam access to only those scheduled for their specific Class assignment

### Requirement 7

**User Story:** As a system user, I want secure authentication and role-based access, so that the system maintains data integrity and appropriate access controls.

#### Acceptance Criteria

1. THE School_Exam_System SHALL authenticate all User_Account credentials before granting system access
2. WHEN users log in, THE School_Exam_System SHALL determine user role and apply appropriate access permissions
3. THE School_Exam_System SHALL maintain secure session management for all authenticated users
4. THE School_Exam_System SHALL log all user actions for audit and security purposes
5. THE School_Exam_System SHALL automatically terminate inactive user sessions after a specified timeout period