package attendance.strategy;

import attendance.model.Student;

import java.util.List;

public interface SelectionStrategy {
    List<Student> select(List<Student> allStudents, int count);
}
