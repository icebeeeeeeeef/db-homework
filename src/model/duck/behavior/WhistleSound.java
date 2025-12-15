package model.duck.behavior;

public class WhistleSound implements SoundBehavior {
    @Override
    public String getName() {
        return "Playful Whistle";
    }

    @Override
    public String makeSound() {
        return "\"Fweeeeee!\"";
    }

    @Override
    public String toString() {
        return getName();
    }
}
