package model.duck.behavior;

public class ChirpSound implements SoundBehavior {
    @Override
    public String getName() {
        return "Gentle Chirp";
    }

    @Override
    public String makeSound() {
        return "\"Chirp chirp!\"";
    }

    @Override
    public String toString() {
        return getName();
    }
}
