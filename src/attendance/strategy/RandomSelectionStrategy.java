package attendance.strategy;

import attendance.model.Student;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Random;

public class RandomSelectionStrategy implements SelectionStrategy {
    private final Random random = new Random();

    @Override
    public List<Student> select(List<Student> allStudents, int count) {
        List<Student> copy = new ArrayList<>(allStudents);
        Collections.shuffle(copy, random);
        if (count <= 0 || count >= copy.size()) {
            return copy;
        }
        return new ArrayList<>(copy.subList(0, count));
    }
}
