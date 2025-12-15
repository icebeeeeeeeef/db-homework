package attendance.ui;

import attendance.controller.AttendanceController;
import attendance.model.SessionSummary;
import attendance.model.Student;
import attendance.model.StudentStatistics;

import javax.swing.*;
import javax.swing.table.DefaultTableModel;
import java.awt.*;
import java.sql.SQLException;
import java.time.format.DateTimeFormatter;
import java.util.List;

public class AttendanceApp extends JFrame {
    private final AttendanceController controller = new AttendanceController();
    private DefaultTableModel studentTableModel;
    private DefaultTableModel summaryTableModel;
    private DefaultTableModel statsTableModel;
    private JLabel summaryLabel;
    private JButton startButton;
    private JButton presentButton;
    private JButton absentButton;
    private JRadioButton fullRadio;
    private JRadioButton randomRadio;
    private JSpinner randomCountSpinner;
    private JLabel currentStudentLabel;
    private JLabel currentIndexLabel;
    private List<Student> activeStudents = List.of();
    private int currentIndex = -1;
    private boolean sessionActive = false;
    private final java.util.Map<Integer, String> localStatus = new java.util.HashMap<>();

    public AttendanceApp() {
        setTitle("Classroom Attendance");
        setSize(1400, 900);
        setLocationRelativeTo(null);
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        initUI();
        refreshStudentTable();
        refreshStatistics();
    }

    private void initUI() {
        JTabbedPane tabs = new JTabbedPane();
        tabs.addTab("Students", buildStudentPanel());
        tabs.addTab("Attendance", buildAttendancePanel());
        tabs.addTab("Statistics", buildStatisticsPanel());
        setContentPane(tabs);
    }

    private JPanel buildStudentPanel() {
        JPanel panel = new JPanel(new BorderLayout());
        studentTableModel = new DefaultTableModel(new Object[]{"ID", "Student No", "Name", "Photo Path"}, 0) {
            @Override
            public boolean isCellEditable(int row, int column) {
                return false;
            }
        };
        JTable table = new JTable(studentTableModel);
        panel.add(new JScrollPane(table), BorderLayout.CENTER);

        JPanel form = new JPanel(new GridLayout(2, 4, 8, 8));
        JTextField studentNoField = new JTextField();
        JTextField nameField = new JTextField();
        JTextField photoField = new JTextField();
        JButton addButton = new JButton("Add Student");

        form.add(new JLabel("Student No"));
        form.add(studentNoField);
        form.add(new JLabel("Name"));
        form.add(nameField);
        form.add(new JLabel("Photo Path (optional)"));
        form.add(photoField);
        form.add(new JLabel());
        form.add(addButton);

        addButton.addActionListener(e -> {
            String studentNo = studentNoField.getText().trim();
            String name = nameField.getText().trim();
            String photo = photoField.getText().trim();
            if (studentNo.isEmpty() || name.isEmpty()) {
                showError("Student No and name are required");
                return;
            }
            try {
                controller.addStudent(studentNo, name, photo.isEmpty() ? null : photo);
                refreshStudentTable();
                studentNoField.setText("");
                nameField.setText("");
                photoField.setText("");
                showInfo("Student added");
            } catch (SQLException ex) {
                showError("Failed to add student: " + ex.getMessage());
            }
        });

        panel.add(form, BorderLayout.SOUTH);
        return panel;
    }

    private JPanel buildAttendancePanel() {
        JPanel panel = new JPanel(new BorderLayout());

        JPanel topConfig = new JPanel(new FlowLayout(FlowLayout.LEFT));
        fullRadio = new JRadioButton("Full", true);
        randomRadio = new JRadioButton("Random");
        ButtonGroup group = new ButtonGroup();
        group.add(fullRadio);
        group.add(randomRadio);
        randomCountSpinner = new JSpinner(new SpinnerNumberModel(1, 1, 200, 1));
        randomCountSpinner.setEnabled(false);
        startButton = new JButton("Start");
        presentButton = new JButton("Present");
        absentButton = new JButton("Absent");
        presentButton.setEnabled(false);
        absentButton.setEnabled(false);

        topConfig.add(new JLabel("Mode:"));
        topConfig.add(fullRadio);
        topConfig.add(randomRadio);
        topConfig.add(new JLabel("Random count:"));
        topConfig.add(randomCountSpinner);
        topConfig.add(startButton);
        topConfig.add(new JLabel(" | Mark:"));
        topConfig.add(presentButton);
        topConfig.add(absentButton);

        JPanel currentPanel = new JPanel(new BorderLayout());
        currentPanel.setBorder(BorderFactory.createTitledBorder("Current Student"));
        currentStudentLabel = new JLabel("No active session", SwingConstants.CENTER);
        currentStudentLabel.setFont(currentStudentLabel.getFont().deriveFont(Font.BOLD, 28f));
        currentIndexLabel = new JLabel("", SwingConstants.CENTER);
        currentIndexLabel.setFont(currentIndexLabel.getFont().deriveFont(Font.PLAIN, 18f));
        currentPanel.add(currentStudentLabel, BorderLayout.CENTER);
        currentPanel.add(currentIndexLabel, BorderLayout.SOUTH);

        summaryTableModel = new DefaultTableModel(new Object[]{"Student No", "Name", "Status"}, 0) {
            @Override
            public boolean isCellEditable(int row, int column) {
                return false;
            }
        };
        JTable summaryTable = new JTable(summaryTableModel);
        JScrollPane summaryScroll = new JScrollPane(summaryTable);
        summaryScroll.setBorder(BorderFactory.createTitledBorder("Session Summary"));
        summaryScroll.setPreferredSize(new Dimension(200, 260));

        startButton.addActionListener(e -> startAttendance());
        presentButton.addActionListener(e -> markAndNext(true));
        absentButton.addActionListener(e -> markAndNext(false));
        fullRadio.addActionListener(e -> randomCountSpinner.setEnabled(false));
        randomRadio.addActionListener(e -> randomCountSpinner.setEnabled(true));

        panel.add(topConfig, BorderLayout.NORTH);
        panel.add(currentPanel, BorderLayout.CENTER);
        panel.add(summaryScroll, BorderLayout.SOUTH);
        return panel;
    }

    private JPanel buildStatisticsPanel() {
        JPanel panel = new JPanel(new BorderLayout());
        statsTableModel = new DefaultTableModel(new Object[]{"Student No", "Name", "Total", "Absent", "Absence Rate (%)"}, 0) {
            @Override
            public boolean isCellEditable(int row, int column) {
                return false;
            }
        };
        JTable statsTable = new JTable(statsTableModel);
        panel.add(new JScrollPane(statsTable), BorderLayout.CENTER);

        JPanel footer = new JPanel(new BorderLayout());
        summaryLabel = new JLabel("No attendance yet");
        JButton refreshButton = new JButton("Refresh");
        refreshButton.addActionListener(e -> refreshStatistics());
        footer.add(summaryLabel, BorderLayout.CENTER);
        footer.add(refreshButton, BorderLayout.EAST);
        panel.add(footer, BorderLayout.SOUTH);
        return panel;
    }

    private void refreshStudentTable() {
        try {
            List<Student> students = controller.getAllStudents();
            studentTableModel.setRowCount(0);
            for (Student s : students) {
                studentTableModel.addRow(new Object[]{s.getId(), s.getStudentNo(), s.getName(), s.getPhotoPath()});
            }
        } catch (SQLException e) {
            showError("Failed to load students: " + e.getMessage());
        }
    }

    private void startAttendance() {
        String mode = fullRadio.isSelected() ? "full" : "random";
        int count = (Integer) randomCountSpinner.getValue();
        try {
            List<Student> selected = controller.startSession(mode, count);
            if (selected.isEmpty()) {
                showError("No students available. Please add students first.");
                return;
            }
            activeStudents = selected;
            currentIndex = 0;
            localStatus.clear();
            summaryTableModel.setRowCount(0);
            sessionActive = true;
            toggleControls(false);
            updateCurrentStudent();
            summaryLabel.setText("Attendance in progress...");
        } catch (Exception e) {
            showError("Cannot start attendance: " + e.getMessage());
        }
    }

    private void markAndNext(boolean present) {
        if (!sessionActive || currentIndex < 0 || currentIndex >= activeStudents.size()) {
            return;
        }
        Student current = activeStudents.get(currentIndex);
        controller.markAttendance(current.getId(), present);
        localStatus.put(current.getId(), present ? "Present" : "Absent");
        currentIndex++;
        if (currentIndex >= activeStudents.size()) {
            finishAttendanceInternal();
        } else {
            updateCurrentStudent();
        }
    }

    private void finishAttendanceInternal() {
        if (!sessionActive) return;
        try {
            SessionSummary summary = controller.finishSession();
            sessionActive = false;
            toggleControls(true);
            populateSummaryTable();
            if (summary != null) {
                DateTimeFormatter fmt = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm");
                summaryLabel.setText(String.format("Last session: %s | mode: %s | present: %d absent: %d",
                        summary.getSession().getSessionTime().format(fmt),
                        summary.getSession().getMode(),
                        summary.getPresentCount(),
                        summary.getAbsentCount()));
            }
            refreshStatistics();
            showInfo("Attendance saved");
        } catch (Exception e) {
            showError("Failed to finish: " + e.getMessage());
        } finally {
            currentStudentLabel.setText("Session complete");
            currentIndexLabel.setText("");
        }
    }

    private void updateCurrentStudent() {
        if (currentIndex < 0 || currentIndex >= activeStudents.size()) {
            currentStudentLabel.setText("No active session");
            currentIndexLabel.setText("");
            presentButton.setEnabled(false);
            absentButton.setEnabled(false);
            return;
        }
        Student s = activeStudents.get(currentIndex);
        currentStudentLabel.setText(String.format("%s - %s", s.getStudentNo(), s.getName()));
        currentIndexLabel.setText(String.format("Student %d of %d", currentIndex + 1, activeStudents.size()));
        presentButton.setEnabled(true);
        absentButton.setEnabled(true);
    }

    private void populateSummaryTable() {
        summaryTableModel.setRowCount(0);
        for (Student s : activeStudents) {
            String status = localStatus.getOrDefault(s.getId(), "Absent");
            summaryTableModel.addRow(new Object[]{s.getStudentNo(), s.getName(), status});
        }
    }

    private void toggleControls(boolean enableStart) {
        startButton.setEnabled(enableStart);
        fullRadio.setEnabled(enableStart);
        randomRadio.setEnabled(enableStart);
        randomCountSpinner.setEnabled(enableStart && randomRadio.isSelected());
        presentButton.setEnabled(!enableStart);
        absentButton.setEnabled(!enableStart);
    }

    private void refreshStatistics() {
        try {
            List<StudentStatistics> stats = controller.getStudentStatistics();
            statsTableModel.setRowCount(0);
            for (StudentStatistics stat : stats) {
                statsTableModel.addRow(new Object[]{
                        stat.getStudent().getStudentNo(),
                        stat.getStudent().getName(),
                        stat.getTotalCount(),
                        stat.getAbsentCount(),
                        String.format("%.1f", stat.getAbsenceRate() * 100)
                });
            }
        } catch (SQLException e) {
            showError("Failed to refresh stats: " + e.getMessage());
        }
    }

    private void showError(String msg) {
        JOptionPane.showMessageDialog(this, msg, "Error", JOptionPane.ERROR_MESSAGE);
    }

    private void showInfo(String msg) {
        JOptionPane.showMessageDialog(this, msg, "Info", JOptionPane.INFORMATION_MESSAGE);
    }

    public static void launch() {
        SwingUtilities.invokeLater(() -> new AttendanceApp().setVisible(true));
    }

    public static void main(String[] args) {
        launch();
    }
}
