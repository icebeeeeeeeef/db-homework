package attendance.model;

public class AttendanceRecord {
    private Integer id;
    private int sessionId;
    private int studentId;
    private String status; // present / absent

    public AttendanceRecord(Integer id, int sessionId, int studentId, String status) {
        this.id = id;
        this.sessionId = sessionId;
        this.studentId = studentId;
        this.status = status;
    }

    public AttendanceRecord(int sessionId, int studentId, String status) {
        this(null, sessionId, studentId, status);
    }

    public Integer getId() {
        return id;
    }

    public int getSessionId() {
        return sessionId;
    }

    public int getStudentId() {
        return studentId;
    }

    public String getStatus() {
        return status;
    }
}
