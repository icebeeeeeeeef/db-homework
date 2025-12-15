package gui;

import javax.swing.*;
import java.awt.*;

public class DuckAvatarPanel extends JPanel {
    private boolean showHat = true;
    private boolean showScarf = false;
    private boolean showGlasses = false;
    private boolean showBowtie = true;

    private Color hatColor = new Color(30, 144, 255);     // 道具蓝色
    private Color scarfColor = new Color(255, 99, 71);    // 番茄红
    private Color bowtieColor = new Color(220, 20, 60);   // 唐老鸭经典红色

    public DuckAvatarPanel() {
        setPreferredSize(new Dimension(360, 260));
        setOpaque(false);
    }

    public void setShowHat(boolean show) { this.showHat = show; repaint(); }
    public void setShowScarf(boolean show) { this.showScarf = show; repaint(); }
    public void setShowGlasses(boolean show) { this.showGlasses = show; repaint(); }
    public void setShowBowtie(boolean show) { this.showBowtie = show; repaint(); }

    public boolean isShowHat() { return showHat; }
    public boolean isShowScarf() { return showScarf; }
    public boolean isShowGlasses() { return showGlasses; }
    public boolean isShowBowtie() { return showBowtie; }

    public void setHatColor(Color c) { if (c != null) { this.hatColor = c; repaint(); } }
    public void setScarfColor(Color c) { if (c != null) { this.scarfColor = c; repaint(); } }
    public void setBowtieColor(Color c) { if (c != null) { this.bowtieColor = c; repaint(); } }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);
        Graphics2D g2 = (Graphics2D) g;
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);

        int w = getWidth();
        int h = getHeight();
        int cx = w / 2;
        int baseY = h - 20;

        // 身体（白色羽毛）
        int bodyW = Math.min(220, (int)(w * 0.75));
        int bodyH = (int)(bodyW * 0.55);
        int bodyX = cx - bodyW / 2;
        int bodyY = baseY - bodyH;
        g2.setColor(Color.WHITE);
        g2.fillOval(bodyX, bodyY, bodyW, bodyH);

        // 头部
        int headW = (int)(bodyW * 0.52);
        int headH = headW;
        int headX = cx - headW / 3;
        int headY = bodyY - headH / 2;
        g2.fillOval(headX, headY, headW, headH);

        // 眼睛
        int eyeR = Math.max(4, headW / 12);
        g2.setColor(Color.BLACK);
        g2.fillOval(headX + headW/3 - eyeR, headY + headH/3, eyeR, eyeR);
        g2.fillOval(headX + headW*2/3 - eyeR, headY + headH/3, eyeR, eyeR);

        // 嘴（喙）
        int beakW = headW / 2;
        int beakH = Math.max(10, headH / 5);
        int beakX = headX + headW/2 - beakW/2;
        int beakY = headY + headH/2;
        g2.setColor(new Color(255, 140, 0));
        g2.fillRoundRect(beakX, beakY, beakW, beakH, beakH, beakH);

        // 翅膀
        g2.setColor(Color.WHITE);
        g2.fillOval(bodyX - bodyW/6, bodyY + bodyH/3, bodyW/3, bodyH/2);
        g2.fillOval(bodyX + bodyW - bodyW/6, bodyY + bodyH/3, bodyW/3, bodyH/2);

        // 水手服领口（默认蓝色）
        int collarH = bodyH / 3;
        int collarY = bodyY + bodyH/3;
        g2.setColor(new Color(25, 116, 210));
        g2.fillRoundRect(bodyX + bodyW/6, collarY, bodyW*2/3, collarH, 20, 20);
        g2.setColor(new Color(255, 215, 0));
        g2.fillRoundRect(bodyX + bodyW/6, collarY + collarH - 12, bodyW*2/3, 12, 12, 12);

        // 围巾
        if (showScarf) {
            g2.setColor(scarfColor);
            int scarfY = headY + headH - headH/6;
            g2.fillRoundRect(headX + headW/6, scarfY, headW*2/3, headH/6, 10, 10);
            g2.fillRoundRect(headX + headW/2, scarfY + headH/6 - 4, headW/3, headH/5, 10, 10);
        }

        // 眼镜
        if (showGlasses) {
            g2.setStroke(new BasicStroke(3f));
            g2.setColor(Color.DARK_GRAY);
            int r = headW/6;
            int gx1 = headX + headW/3 - r;
            int gy = headY + headH/3 - r/2;
            g2.drawOval(gx1, gy, r*2, r*2);
            int gx2 = headX + headW*2/3 - r;
            g2.drawOval(gx2, gy, r*2, r*2);
            g2.drawLine(gx1 + r*2, gy + r, gx2, gy + r);
        }

        // 领结
        if (showBowtie) {
            g2.setColor(bowtieColor);
            int by = collarY + collarH/2;
            int bw = headW/5;
            Polygon left = new Polygon(new int[]{cx - bw - 12, cx - 12, cx - bw - 12}, new int[]{by, by - bw/2, by - bw}, 3);
            Polygon right = new Polygon(new int[]{cx + bw + 12, cx + 12, cx + bw + 12}, new int[]{by, by - bw/2, by - bw}, 3);
            g2.fillPolygon(left);
            g2.fillPolygon(right);
            g2.fillOval(cx - 12, by - bw/2 - 6, 24, 24);
        }

        // 帽子
        if (showHat) {
            g2.setColor(hatColor);
            int brimW = headW;
            int brimH = headH/9;
            int brimX = headX + (headW - brimW)/2;
            int brimY = headY - brimH/2 + 6;
            g2.fillRoundRect(brimX, brimY, brimW, brimH, brimH, brimH);
            g2.setColor(Color.BLACK);
            g2.fillRoundRect(brimX, brimY + brimH - 4, brimW, 4, 4, 4);
            int hatW = headW/2;
            int hatH = headH/1;
            int hatX = headX + headW/2 - hatW/2;
            int hatY = brimY - hatH + 6;
            g2.setColor(hatColor);
            g2.fillRoundRect(hatX, hatY, hatW, hatH, 18, 18);
        }
    }
}
