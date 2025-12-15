package attendance.strategy;

import attendance.model.Student;

import java.util.ArrayList;
import java.util.List;

public class FullSelectionStrategy implements SelectionStrategy {
    @Override
    public List<Student> select(List<Student> allStudents, int count) {
        return new ArrayList<>(allStudents);
    }
}
