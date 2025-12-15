package model.duck.behavior;

public class DuckBehaviorProfile implements Cloneable {
    private ActionBehavior actionBehavior;
    private SoundBehavior soundBehavior;

    public DuckBehaviorProfile(ActionBehavior actionBehavior, SoundBehavior soundBehavior) {
        this.actionBehavior = actionBehavior;
        this.soundBehavior = soundBehavior;
    }

    public DuckBehaviorProfile() {
        this(null, null);
    }

    public ActionBehavior getActionBehavior() {
        return actionBehavior;
    }

    public void setActionBehavior(ActionBehavior actionBehavior) {
        this.actionBehavior = actionBehavior;
    }

    public SoundBehavior getSoundBehavior() {
        return soundBehavior;
    }

    public void setSoundBehavior(SoundBehavior soundBehavior) {
        this.soundBehavior = soundBehavior;
    }

    @Override
    public DuckBehaviorProfile clone() {
        try {
            return (DuckBehaviorProfile) super.clone();
        } catch (CloneNotSupportedException e) {
            return new DuckBehaviorProfile(actionBehavior, soundBehavior);
        }
    }
}
