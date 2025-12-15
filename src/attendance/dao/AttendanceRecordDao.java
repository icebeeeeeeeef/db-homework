package attendance.dao;

import attendance.config.DbConfig;
import attendance.model.AttendanceRecord;
import attendance.model.Student;
import attendance.model.StudentStatistics;

import java.sql.*;
import java.util.ArrayList;
import java.util.List;

public class AttendanceRecordDao {
    public void insertRecords(List<AttendanceRecord> records) throws SQLException {
        if (records.isEmpty()) {
            return;
        }
        String sql = "INSERT INTO attendance_record (session_id, student_id, status) VALUES (?, ?, ?)";
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql)) {
            for (AttendanceRecord record : records) {
                ps.setInt(1, record.getSessionId());
                ps.setInt(2, record.getStudentId());
                ps.setString(3, record.getStatus());
                ps.addBatch();
            }
            ps.executeBatch();
        }
    }

    public List<AttendanceRecord> findBySessionId(int sessionId) throws SQLException {
        String sql = "SELECT id, session_id, student_id, status FROM attendance_record WHERE session_id = ?";
        List<AttendanceRecord> list = new ArrayList<>();
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql)) {
            ps.setInt(1, sessionId);
            try (ResultSet rs = ps.executeQuery()) {
                while (rs.next()) {
                    list.add(map(rs));
                }
            }
        }
        return list;
    }

    public List<StudentStatistics> fetchStudentStatistics() throws SQLException {
        String sql = "SELECT s.id, s.student_no, s.name, s.photo_path, " +
                "COUNT(r.id) AS total_count, " +
                "SUM(CASE WHEN r.status = 'absent' THEN 1 ELSE 0 END) AS absent_count " +
                "FROM student s " +
                "LEFT JOIN attendance_record r ON s.id = r.student_id " +
                "GROUP BY s.id, s.student_no, s.name, s.photo_path " +
                "ORDER BY s.id";
        List<StudentStatistics> stats = new ArrayList<>();
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql);
             ResultSet rs = ps.executeQuery()) {
            while (rs.next()) {
                Student student = new Student(
                        rs.getInt("id"),
                        rs.getString("student_no"),
                        rs.getString("name"),
                        rs.getString("photo_path")
                );
                int total = rs.getInt("total_count");
                int absent = rs.getInt("absent_count");
                stats.add(new StudentStatistics(student, total, absent));
            }
        }
        return stats;
    }

    private AttendanceRecord map(ResultSet rs) throws SQLException {
        return new AttendanceRecord(
                rs.getInt("id"),
                rs.getInt("session_id"),
                rs.getInt("student_id"),
                rs.getString("status")
        );
    }
}
