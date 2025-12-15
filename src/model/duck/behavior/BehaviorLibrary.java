package model.duck.behavior;

import java.util.List;

/**
 * Provides predefined action and sound behaviors.
 */
public final class BehaviorLibrary {
    private BehaviorLibrary() {}

    public static List<ActionBehavior> availableActions() {
        return List.of(
                new FlyAction(),
                new RunAction(),
                new SwimAction()
        );
    }

    public static List<SoundBehavior> availableSounds() {
        return List.of(
                new QuackSound(),
                new ChirpSound(),
                new WhistleSound()
        );
    }
}
