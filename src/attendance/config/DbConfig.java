package attendance.config;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

/**
 * Centralized JDBC configuration. Update host/port/database/user/password to switch DB in one place.
 */
public class DbConfig {
    private static final String HOST = System.getProperty("db.host",
            System.getenv().getOrDefault("DB_HOST", "localhost"));
    private static final String PORT = System.getProperty("db.port",
            System.getenv().getOrDefault("DB_PORT", "3366"));
    private static final String DATABASE = System.getProperty("db.name",
            System.getenv().getOrDefault("DB_NAME", "attendance_db"));
    private static final String USER = System.getProperty("db.user",
            System.getenv().getOrDefault("DB_USER", "attendance_user"));
    private static final String PASSWORD = System.getProperty("db.password",
            System.getenv().getOrDefault("DB_PASSWORD", "attendance_pass"));

    static {
        try {
            Class.forName("com.mysql.cj.jdbc.Driver");
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("MySQL driver not found", e);
        }
    }

    private DbConfig() {
    }

    public static Connection getConnection() throws SQLException {
        String url = String.format(
                "jdbc:mysql://%s:%s/%s?useSSL=false&allowPublicKeyRetrieval=true&serverTimezone=UTC",
                HOST, PORT, DATABASE);
        return DriverManager.getConnection(url, USER, PASSWORD);
    }
}
