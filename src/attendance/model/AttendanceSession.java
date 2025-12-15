package attendance.model;

import java.time.LocalDateTime;

public class AttendanceSession {
    private Integer id;
    private LocalDateTime sessionTime;
    private String mode;
    private int totalCount;

    public AttendanceSession(Integer id, LocalDateTime sessionTime, String mode, int totalCount) {
        this.id = id;
        this.sessionTime = sessionTime;
        this.mode = mode;
        this.totalCount = totalCount;
    }

    public AttendanceSession(LocalDateTime sessionTime, String mode, int totalCount) {
        this(null, sessionTime, mode, totalCount);
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public LocalDateTime getSessionTime() {
        return sessionTime;
    }

    public String getMode() {
        return mode;
    }

    public int getTotalCount() {
        return totalCount;
    }
}
