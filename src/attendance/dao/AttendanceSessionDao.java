package attendance.dao;

import attendance.config.DbConfig;
import attendance.model.AttendanceSession;

import java.sql.*;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Optional;

public class AttendanceSessionDao {
    public AttendanceSession insert(AttendanceSession session) throws SQLException {
        String sql = "INSERT INTO attendance_session (session_time, mode, total_count) VALUES (?, ?, ?)";
        try (Connection conn = DbConfig.getConnection();
             PreparedStatement ps = conn.prepareStatement(sql, Statement.RETURN_GENERATED_KEYS)) {
            ps.setTimestamp(1, Timestamp.valueOf(session.getSessionTime()));
            ps.setString(2, session.getMode());
            ps.setInt(3, session.getTotalCount());
            ps.executeUpdate();
            try (ResultSet rs = ps.getGeneratedKeys()) {
                if (rs.next()) {
                    session.setId(rs.getInt(1));
                }
            }
            return session;
        }
    }

    public Optional<AttendanceSession> findById(int id) throws SQLException {
        String sql = "SELECT id, session_time, mode, total_count FROM attendance_session WHERE id = ?";
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

    private AttendanceSession map(ResultSet rs) throws SQLException {
        Timestamp ts = rs.getTimestamp("session_time");
        LocalDateTime time = ts.toInstant().atZone(ZoneId.systemDefault()).toLocalDateTime();
        return new AttendanceSession(
                rs.getInt("id"),
                time,
                rs.getString("mode"),
                rs.getInt("total_count")
        );
    }
}
