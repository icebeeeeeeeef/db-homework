package model.duck.behavior;

/**
 * Defines a duck action such as running or flying.
 */
public interface ActionBehavior {
    String getName();

    /**
     * Short description for chat/stage narration.
     */
    String perform();

    default String toDisplayString() {
        return getName();
    }
}
