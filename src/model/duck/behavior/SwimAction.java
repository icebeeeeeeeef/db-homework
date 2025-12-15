package model.duck.behavior;

public class SwimAction implements ActionBehavior {
    @Override
    public String getName() {
        return "Swim";
    }

    @Override
    public String perform() {
        return "paddles gracefully as if the stage were a lake.";
    }

    @Override
    public String toString() {
        return getName();
    }
}
