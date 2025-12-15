package attendance.model;

public class SessionSummary {
    private final AttendanceSession session;
    private final int presentCount;
    private final int absentCount;

    public SessionSummary(AttendanceSession session, int presentCount, int absentCount) {
        this.session = session;
        this.presentCount = presentCount;
        this.absentCount = absentCount;
    }

    public AttendanceSession getSession() {
        return session;
    }

    public int getPresentCount() {
        return presentCount;
    }

    public int getAbsentCount() {
        return absentCount;
    }

    public int getTotal() {
        return presentCount + absentCount;
    }
}
