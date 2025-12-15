package attendance.controller;

import attendance.model.AttendanceRecord;
import attendance.model.AttendanceSession;
import attendance.model.SessionSummary;
import attendance.model.Student;
import attendance.model.StudentStatistics;
import attendance.service.AttendanceService;
import attendance.service.StudentService;
import attendance.strategy.FullSelectionStrategy;
import attendance.strategy.RandomSelectionStrategy;
import attendance.strategy.SelectionStrategy;

import java.sql.SQLException;
import java.util.List;

public class AttendanceController {
    private final StudentService studentService;
    private final AttendanceService attendanceService;

    public AttendanceController() {
        this.studentService = new StudentService();
        this.attendanceService = new AttendanceService();
    }

    public Student addStudent(String studentNo, String name, String photoPath) throws SQLException {
        return studentService.addStudent(studentNo, name, photoPath);
    }

    public List<Student> getAllStudents() throws SQLException {
        return studentService.getAllStudents();
    }

    public List<Student> startSession(String mode, int randomCount) throws SQLException {
        List<Student> all = studentService.getAllStudents();
        SelectionStrategy strategy = resolveStrategy(mode);
        List<Student> selected = strategy.select(all, randomCount);
        attendanceService.startSession(mode, selected);
        return selected;
    }

    public void markAttendance(int studentId, boolean present) {
        attendanceService.markAttendance(studentId, present ? "present" : "absent");
    }

    public SessionSummary finishSession() throws SQLException {
        attendanceService.finishSession();
        return attendanceService.getSessionSummary();
    }

    public List<StudentStatistics> getStudentStatistics() throws SQLException {
        return attendanceService.getStudentStatistics();
    }

    public AttendanceSession getCurrentSession() {
        return attendanceService.getCurrentSession();
    }

    private SelectionStrategy resolveStrategy(String mode) {
        if ("random".equalsIgnoreCase(mode)) {
            return new RandomSelectionStrategy();
        }
        return new FullSelectionStrategy();
    }
}
