package model.duck.customization;

import java.awt.*;

/**
 * Basic Donald-style duck without accessories.
 */
public class BaseDuck implements DuckAppearance {
    private final Color bodyColor;
    private final Color beakColor;
    private final Color eyeColor;

    public BaseDuck() {
        this(new Color(245, 248, 255), new Color(255, 166, 0), Color.BLACK);
    }

    public BaseDuck(Color bodyColor, Color beakColor, Color eyeColor) {
        this.bodyColor = bodyColor;
        this.beakColor = beakColor;
        this.eyeColor = eyeColor;
    }

    @Override
    public void paint(Graphics2D g2, Rectangle bounds) {
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);

        int w = bounds.width;
        int h = bounds.height;
        int cx = bounds.x + w / 2;
        int bodyW = (int) (w * 0.7);
        int bodyH = (int) (h * 0.55);
        int bodyX = cx - bodyW / 2;
        int bodyY = bounds.y + h - bodyH;

        g2.setColor(bodyColor);
        g2.fillOval(bodyX, bodyY, bodyW, bodyH);

        int headW = (int) (bodyW * 0.45);
        int headX = cx - headW / 2;
        int headY = bodyY - headW / 2;
        g2.fillOval(headX, headY, headW, headW);

        int eyeR = Math.max(4, headW / 10);
        g2.setColor(eyeColor);
        g2.fillOval(headX + headW / 3 - eyeR / 2, headY + headW / 3, eyeR, eyeR);
        g2.fillOval(headX + headW * 2 / 3 - eyeR / 2, headY + headW / 3, eyeR, eyeR);

        int beakW = headW / 2;
        int beakH = Math.max(10, headW / 4);
        int beakX = headX + headW / 2 - beakW / 2;
        int beakY = headY + headW / 2;
        g2.setColor(beakColor);
        g2.fillRoundRect(beakX, beakY, beakW, beakH, beakH, beakH);

        // Fins
        g2.setColor(bodyColor);
        g2.fillOval(bodyX - bodyW / 6, bodyY + bodyH / 4, bodyW / 3, bodyH / 2);
        g2.fillOval(bodyX + bodyW - bodyW / 6, bodyY + bodyH / 4, bodyW / 3, bodyH / 2);
    }
}
