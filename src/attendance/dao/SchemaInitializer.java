package attendance.dao;

import attendance.config.DbConfig;

import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

/**
 * Ensures required schema and seed data exist even if Docker init scripts did not run.
 */
public class SchemaInitializer {
    private static boolean initialized = false;

    private SchemaInitializer() {
    }

    public static synchronized void ensureSchema() {
        if (initialized) {
            return;
        }
        try (Connection conn = DbConfig.getConnection();
             Statement st = conn.createStatement()) {

            st.executeUpdate("CREATE TABLE IF NOT EXISTS student (" +
                    "id INT PRIMARY KEY AUTO_INCREMENT," +
                    "student_no VARCHAR(20) UNIQUE NOT NULL," +
                    "name VARCHAR(50) NOT NULL," +
                    "photo_path VARCHAR(255)" +
                    ")");

            st.executeUpdate("CREATE TABLE IF NOT EXISTS attendance_session (" +
                    "id INT PRIMARY KEY AUTO_INCREMENT," +
                    "session_time DATETIME NOT NULL," +
                    "mode VARCHAR(20) NOT NULL," +
                    "total_count INT NOT NULL" +
                    ")");

            st.executeUpdate("CREATE TABLE IF NOT EXISTS attendance_record (" +
                    "id INT PRIMARY KEY AUTO_INCREMENT," +
                    "session_id INT NOT NULL," +
                    "student_id INT NOT NULL," +
                    "status VARCHAR(20) NOT NULL," +
                    "FOREIGN KEY (session_id) REFERENCES attendance_session(id)," +
                    "FOREIGN KEY (student_id) REFERENCES student(id)" +
                    ")");

            seedStudentsIfEmpty(st);
            initialized = true;
        } catch (SQLException e) {
            throw new RuntimeException("Failed to initialize schema", e);
        }
    }

    private static void seedStudentsIfEmpty(Statement st) throws SQLException {
        try (ResultSet rs = st.executeQuery("SELECT COUNT(*) FROM student")) {
            if (rs.next() && rs.getInt(1) > 0) {
                return;
            }
        } catch (SQLException ignored) {
            // if table was just created, continue to insert seed data
        }
        st.executeUpdate("INSERT INTO student (student_no, name, photo_path) VALUES " +
                "('2023001', 'Alice', NULL)," +
                "('2023002', 'Bob', NULL)," +
                "('2023003', 'Charlie', NULL)," +
                "('2023004', 'Diana', NULL)," +
                "('2023005', 'Ethan', NULL)," +
                "('2023006', 'Fiona', NULL)," +
                "('2023007', 'George', NULL)," +
                "('2023008', 'Hannah', NULL)," +
                "('2023009', 'Ivan', NULL)," +
                "('2023010', 'Julia', NULL)");
    }
}
