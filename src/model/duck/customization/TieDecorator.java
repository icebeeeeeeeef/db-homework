package model.duck.customization;

import java.awt.*;

public class TieDecorator extends AccessoryDecorator {
    private final Color color;

    public TieDecorator(DuckAppearance inner, Color color) {
        super(inner);
        this.color = color == null ? new Color(255, 0, 60) : color;
    }

    @Override
    protected void paintAccessory(Graphics2D g2, Rectangle bounds) {
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
        int w = bounds.width;
        int h = bounds.height;
        int tieWidth = (int) (w * 0.08);
        int tieHeight = (int) (h * 0.25);
        int x = bounds.x + w / 2 - tieWidth / 2;
        int y = bounds.y + h / 2;

        g2.setColor(color);
        Polygon knot = new Polygon(
                new int[]{x - tieWidth, x + tieWidth, x + tieWidth * 2, x},
                new int[]{y, y - tieWidth, y, y + tieWidth}, 4);
        g2.fillPolygon(knot);

        Polygon body = new Polygon(
                new int[]{x - tieWidth / 2, x + tieWidth / 2, x + tieWidth, x - tieWidth},
                new int[]{y + tieWidth, y + tieWidth, y + tieHeight, y + tieHeight},
                4);
        g2.fillPolygon(body);
    }
}
