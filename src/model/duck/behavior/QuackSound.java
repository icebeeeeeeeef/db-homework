package model.duck.behavior;

public class QuackSound implements SoundBehavior {
    @Override
    public String getName() {
        return "Classic Quack";
    }

    @Override
    public String makeSound() {
        return "\"Quaaack quack!\"";
    }

    @Override
    public String toString() {
        return getName();
    }
}
