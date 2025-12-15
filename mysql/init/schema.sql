-- Database and tables for attendance system
CREATE DATABASE IF NOT EXISTS attendance_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE attendance_db;

CREATE TABLE IF NOT EXISTS student (
  id INT PRIMARY KEY AUTO_INCREMENT,
  student_no VARCHAR(20) UNIQUE NOT NULL,
  name VARCHAR(50) NOT NULL,
  photo_path VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS attendance_session (
  id INT PRIMARY KEY AUTO_INCREMENT,
  session_time DATETIME NOT NULL,
  mode VARCHAR(20) NOT NULL,
  total_count INT NOT NULL
);

CREATE TABLE IF NOT EXISTS attendance_record (
  id INT PRIMARY KEY AUTO_INCREMENT,
  session_id INT NOT NULL,
  student_id INT NOT NULL,
  status VARCHAR(20) NOT NULL,
  FOREIGN KEY (session_id) REFERENCES attendance_session(id),
  FOREIGN KEY (student_id) REFERENCES student(id)
);

-- ensure app user exists and has privileges (works even if entrypoint already created user)
CREATE USER IF NOT EXISTS 'attendance_user'@'%' IDENTIFIED BY 'attendance_pass';
GRANT ALL PRIVILEGES ON attendance_db.* TO 'attendance_user'@'%';
FLUSH PRIVILEGES;
