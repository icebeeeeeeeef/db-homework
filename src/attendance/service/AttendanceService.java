package attendance.service;

import attendance.dao.AttendanceRecordDao;
import attendance.dao.AttendanceSessionDao;
import attendance.model.AttendanceRecord;
import attendance.model.AttendanceSession;
import attendance.model.SessionSummary;
import attendance.model.Student;
import attendance.model.StudentStatistics;

import java.sql.SQLException;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class AttendanceService {
    private final AttendanceSessionDao sessionDao;
    private final AttendanceRecordDao recordDao;
    private AttendanceSession currentSession;
    private List<Student> selectedStudents = new ArrayList<>();
    private final Map<Integer, String> statusMap = new HashMap<>();
    private Integer lastSessionId;

    public AttendanceService() {
        this.sessionDao = new AttendanceSessionDao();
        this.recordDao = new AttendanceRecordDao();
    }

    public AttendanceSession startSession(String mode, List<Student> students) throws SQLException {
        if (students == null || students.isEmpty()) {
            throw new IllegalArgumentException("没有可用的学生进行点名");
        }
        if (currentSession != null) {
            throw new IllegalStateException("上一场点名尚未结束");
        }
        this.selectedStudents = new ArrayList<>(students);
        this.statusMap.clear();
        for (Student s : students) {
            // 默认未到，教师确认后再置为到
            statusMap.put(s.getId(), "absent");
        }
        AttendanceSession session = new AttendanceSession(LocalDateTime.now(), mode, students.size());
        this.currentSession = sessionDao.insert(session);
        return currentSession;
    }

    public void markAttendance(int studentId, String status) {
        if (currentSession == null) {
            throw new IllegalStateException("点名尚未开始");
        }
        if (!"present".equalsIgnoreCase(status) && !"absent".equalsIgnoreCase(status)) {
            throw new IllegalArgumentException("状态必须是 present 或 absent");
        }
        statusMap.put(studentId, status.toLowerCase());
    }

    public List<AttendanceRecord> finishSession() throws SQLException {
        if (currentSession == null) {
            throw new IllegalStateException("点名尚未开始");
        }
        List<AttendanceRecord> records = new ArrayList<>();
        for (Student student : selectedStudents) {
            String status = statusMap.getOrDefault(student.getId(), "absent");
            records.add(new AttendanceRecord(currentSession.getId(), student.getId(), status));
        }
        recordDao.insertRecords(records);
        lastSessionId = currentSession.getId();
        currentSession = null;
        selectedStudents = new ArrayList<>();
        statusMap.clear();
        return records;
    }

    public SessionSummary getSessionSummary() throws SQLException {
        if (lastSessionId == null) {
            return null;
        }
        return getSessionSummary(lastSessionId);
    }

    public SessionSummary getSessionSummary(int sessionId) throws SQLException {
        AttendanceSession session = sessionDao.findById(sessionId)
                .orElseThrow(() -> new IllegalArgumentException("无效的 sessionId"));
        List<AttendanceRecord> records = recordDao.findBySessionId(sessionId);
        int present = 0;
        int absent = 0;
        for (AttendanceRecord record : records) {
            if ("present".equalsIgnoreCase(record.getStatus())) {
                present++;
            } else {
                absent++;
            }
        }
        return new SessionSummary(session, present, absent);
    }

    public List<StudentStatistics> getStudentStatistics() throws SQLException {
        return recordDao.fetchStudentStatistics();
    }

    public AttendanceSession getCurrentSession() {
        return currentSession;
    }

    public List<Student> getSelectedStudents() {
        return new ArrayList<>(selectedStudents);
    }
}
