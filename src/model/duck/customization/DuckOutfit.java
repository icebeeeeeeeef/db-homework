package model.duck.customization;

import java.awt.Color;

/**
 * Simple DTO describing selected accessories for a duck.
 */
public class DuckOutfit implements Cloneable {
    private boolean hat;
    private boolean scarf;
    private boolean eyes;
    private boolean tie;
    private boolean cane;
    private Color hatColor = new Color(30, 144, 255);
    private Color scarfColor = new Color(255, 99, 71);
    private Color eyeFrameColor = Color.DARK_GRAY;
    private Color tieColor = new Color(220, 20, 60);
    private Color caneColor = new Color(139, 69, 19);

    public DuckAppearance buildAppearance() {
        DuckAppearance appearance = new BaseDuck();
        if (scarf) appearance = new ScarfDecorator(appearance, scarfColor);
        if (tie) appearance = new TieDecorator(appearance, tieColor);
        if (eyes) appearance = new EyeDecorator(appearance, eyeFrameColor, new Color(255, 255, 255, 80));
        if (hat) appearance = new HatDecorator(appearance, hatColor, Color.BLACK);
        if (cane) appearance = new CaneDecorator(appearance, caneColor);
        return appearance;
    }

    public boolean hasHat() { return hat; }
    public void setHat(boolean hat) { this.hat = hat; }
    public boolean hasScarf() { return scarf; }
    public void setScarf(boolean scarf) { this.scarf = scarf; }
    public boolean hasEyes() { return eyes; }
    public void setEyes(boolean eyes) { this.eyes = eyes; }
    public boolean hasTie() { return tie; }
    public void setTie(boolean tie) { this.tie = tie; }
    public boolean hasCane() { return cane; }
    public void setCane(boolean cane) { this.cane = cane; }
    public Color getHatColor() { return hatColor; }
    public void setHatColor(Color hatColor) { if (hatColor != null) this.hatColor = hatColor; }
    public Color getScarfColor() { return scarfColor; }
    public void setScarfColor(Color scarfColor) { if (scarfColor != null) this.scarfColor = scarfColor; }
    public Color getEyeFrameColor() { return eyeFrameColor; }
    public void setEyeFrameColor(Color eyeFrameColor) { if (eyeFrameColor != null) this.eyeFrameColor = eyeFrameColor; }
    public Color getTieColor() { return tieColor; }
    public void setTieColor(Color tieColor) { if (tieColor != null) this.tieColor = tieColor; }
    public Color getCaneColor() { return caneColor; }
    public void setCaneColor(Color caneColor) { if (caneColor != null) this.caneColor = caneColor; }

    @Override
    public DuckOutfit clone() {
        try {
            return (DuckOutfit) super.clone();
        } catch (CloneNotSupportedException e) {
            DuckOutfit copy = new DuckOutfit();
            copy.hat = hat;
            copy.scarf = scarf;
            copy.eyes = eyes;
            copy.tie = tie;
            copy.cane = cane;
            copy.hatColor = hatColor;
            copy.scarfColor = scarfColor;
            copy.eyeFrameColor = eyeFrameColor;
            copy.tieColor = tieColor;
            copy.caneColor = caneColor;
            return copy;
        }
    }
}
