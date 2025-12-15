package attendance.dao;

import attendance.config.DbConfig;
import attendance.model.Student;

import java.sql.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

public class StudentDao {
    static {
        SchemaInitializer.ensureSchema();
    }

    public Student insert(Student student) throws SQLException {
        String sql = "INSERT INTO student (student_no, name, photo_path) VALUES (?, ?, ?)";
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql, Statement.RETURN_GENERATED_KEYS)) {
            ps.setString(1, student.getStudentNo());
            ps.setString(2, student.getName());
            ps.setString(3, student.getPhotoPath());
            ps.executeUpdate();
            try (ResultSet rs = ps.getGeneratedKeys()) {
                if (rs.next()) {
                    student.setId(rs.getInt(1));
                }
            }
            return student;
        }
    }

    public List<Student> findAll() throws SQLException {
        String sql = "SELECT id, student_no, name, photo_path FROM student ORDER BY id";
        List<Student> result = new ArrayList<>();
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql);
             ResultSet rs = ps.executeQuery()) {
            while (rs.next()) {
                result.add(map(rs));
            }
        }
        return result;
    }

    public Optional<Student> findById(int id) throws SQLException {
        String sql = "SELECT id, student_no, name, photo_path FROM student WHERE id = ?";
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql)) {
            ps.setInt(1, id);
            try (ResultSet rs = ps.executeQuery()) {
                if (rs.next()) {
                    return Optional.of(map(rs));
                }
            }
        }
        return Optional.empty();
    }

    private Student map(ResultSet rs) throws SQLException {
        return new Student(
                rs.getInt("id"),
                rs.getString("student_no"),
                rs.getString("name"),
                rs.getString("photo_path")
        );
    }
}
