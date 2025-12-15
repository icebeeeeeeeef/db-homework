package model.duck.behavior;

public class RunAction implements ActionBehavior {
    @Override
    public String getName() {
        return "Run";
    }

    @Override
    public String perform() {
        return "dashes across the stage with quick steps.";
    }

    @Override
    public String toString() {
        return getName();
    }
}
