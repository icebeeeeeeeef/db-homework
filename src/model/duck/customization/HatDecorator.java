package model.duck.customization;

import java.awt.*;

public class HatDecorator extends AccessoryDecorator {
    private final Color hatColor;
    private final Color bandColor;

    public HatDecorator(DuckAppearance inner, Color hatColor, Color bandColor) {
        super(inner);
        this.hatColor = hatColor == null ? new Color(30, 144, 255) : hatColor;
        this.bandColor = bandColor == null ? Color.BLACK : bandColor;
    }

    @Override
    protected void paintAccessory(Graphics2D g2, Rectangle bounds) {
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
        int w = bounds.width;
        int h = bounds.height;
        int hatW = (int) (w * 0.35);
        int hatH = (int) (h * 0.25);
        int brimH = Math.max(6, hatH / 4);
        int x = bounds.x + (w - hatW) / 2;
        int y = bounds.y + h / 3 - hatH;

        g2.setColor(hatColor);
        g2.fillRoundRect(x, y, hatW, hatH, 18, 18);
        g2.fillRoundRect(x - hatW / 5, y + hatH - brimH, hatW + hatW / 3, brimH, brimH, brimH);

        g2.setColor(bandColor);
        g2.fillRoundRect(x + hatW / 6, y + hatH / 2, hatW - hatW / 3, brimH / 2, brimH / 2, brimH / 2);
    }
}
