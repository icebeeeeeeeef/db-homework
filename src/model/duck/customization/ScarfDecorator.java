package model.duck.customization;

import java.awt.*;

public class ScarfDecorator extends AccessoryDecorator {
    private final Color color;

    public ScarfDecorator(DuckAppearance inner, Color color) {
        super(inner);
        this.color = color == null ? new Color(220, 20, 60) : color;
    }

    @Override
    protected void paintAccessory(Graphics2D g2, Rectangle bounds) {
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
        int w = bounds.width;
        int h = bounds.height;
        int scarfW = (int) (w * 0.35);
        int scarfH = Math.max(10, h / 12);
        int x = bounds.x + (w - scarfW) / 2;
        int y = bounds.y + h / 2;

        g2.setColor(color);
        g2.fillRoundRect(x, y, scarfW, scarfH, 10, 10);
        g2.fillRoundRect(x + scarfW / 2, y + scarfH - 4, scarfW / 3, scarfH * 2, 10, 10);
    }
}
