package model.duck.customization;

import java.awt.*;

public class CaneDecorator extends AccessoryDecorator {
    private final Color color;

    public CaneDecorator(DuckAppearance inner, Color color) {
        super(inner);
        this.color = color == null ? new Color(139, 69, 19) : color;
    }

    @Override
    protected void paintAccessory(Graphics2D g2, Rectangle bounds) {
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
        int w = bounds.width;
        int h = bounds.height;
        int caneHeight = (int) (h * 0.75);
        int x = bounds.x + (int) (w * 0.8);
        int y = bounds.y + h - caneHeight;

        g2.setColor(color);
        g2.setStroke(new BasicStroke(Math.max(3f, w / 40f)));
        g2.drawLine(x, y, x, y + caneHeight - 10);
        g2.drawArc(x - 10, y - 20, 20, 20, 0, 180);
    }
}
