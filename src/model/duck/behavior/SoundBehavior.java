package model.duck.behavior;

/**
 * Defines a duck sound (quack, chirp, whistle...).
 */
public interface SoundBehavior {
    String getName();

    /**
     * Returns the vocalization line (for chat/stage bubbles).
     */
    String makeSound();

    default String toDisplayString() {
        return getName();
    }
}
