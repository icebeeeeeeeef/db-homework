package model.duck;

public enum DuckCharacter {
    DONALD,
    DUCKLING_ONE,
    DUCKLING_TWO,
    DUCKLING_THREE;

    public String getDisplayName() {
        switch (this) {
            case DONALD: return "Donald";
            case DUCKLING_ONE: return "Duckling A";
            case DUCKLING_TWO: return "Duckling B";
            case DUCKLING_THREE: return "Duckling C";
            default: return name();
        }
    }
}
