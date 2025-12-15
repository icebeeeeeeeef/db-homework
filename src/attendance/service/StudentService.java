package attendance.service;

import attendance.dao.StudentDao;
import attendance.model.Student;

import java.sql.SQLException;
import java.util.List;

public class StudentService {
    private final StudentDao studentDao;

    public StudentService() {
        this.studentDao = new StudentDao();
    }

    public Student addStudent(String studentNo, String name, String photoPath) throws SQLException {
        Student student = new Student(studentNo, name, photoPath);
        return studentDao.insert(student);
    }

    public List<Student> getAllStudents() throws SQLException {
        return studentDao.findAll();
    }
}
