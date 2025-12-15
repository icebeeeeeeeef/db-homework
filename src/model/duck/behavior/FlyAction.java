package model.duck.behavior;

public class FlyAction implements ActionBehavior {
    @Override
    public String getName() {
        return "Fly";
    }

    @Override
    public String perform() {
        return "takes off and glides through the air.";
    }

    @Override
    public String toString() {
        return getName();
    }
}
