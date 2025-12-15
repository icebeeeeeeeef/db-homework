package model.duck.customization;

import java.awt.*;

public class EyeDecorator extends AccessoryDecorator {
    private final Color frameColor;
    private final Color lensColor;

    public EyeDecorator(DuckAppearance inner, Color frameColor, Color lensColor) {
        super(inner);
        this.frameColor = frameColor == null ? Color.DARK_GRAY : frameColor;
        this.lensColor = lensColor == null ? new Color(255, 255, 255, 120) : lensColor;
    }

    @Override
    protected void paintAccessory(Graphics2D g2, Rectangle bounds) {
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
        int w = bounds.width;
        int h = bounds.height;

        int radius = Math.max(8, w / 12);
        int cx = bounds.x + w / 2;
        int cy = bounds.y + h / 2 - radius;

        int leftX = cx - radius * 2 - 4;
        int rightX = cx + 4;

        g2.setStroke(new BasicStroke(3f));
        g2.setColor(frameColor);
        g2.drawOval(leftX, cy, radius * 2, radius * 2);
        g2.drawOval(rightX, cy, radius * 2, radius * 2);
        g2.drawLine(leftX + radius * 2, cy + radius, rightX, cy + radius);

        g2.setColor(lensColor);
        g2.fillOval(leftX + 2, cy + 2, radius * 2 - 4, radius * 2 - 4);
        g2.fillOval(rightX + 2, cy + 2, radius * 2 - 4, radius * 2 - 4);
    }
}
