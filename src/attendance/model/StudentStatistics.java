package attendance.model;

public class StudentStatistics {
    private final Student student;
    private final int totalCount;
    private final int absentCount;

    public StudentStatistics(Student student, int totalCount, int absentCount) {
        this.student = student;
        this.totalCount = totalCount;
        this.absentCount = absentCount;
    }

    public Student getStudent() {
        return student;
    }

    public int getTotalCount() {
        return totalCount;
    }

    public int getAbsentCount() {
        return absentCount;
    }

    public double getAbsenceRate() {
        if (totalCount == 0) {
            return 0.0;
        }
        return (double) absentCount / totalCount;
    }
}
