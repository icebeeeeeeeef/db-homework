package gui;

import javax.swing.*;
import java.awt.*;
import java.util.List;

public class CodeStatsChartPanel extends JPanel {
    public enum ChartMode { BAR, PIE }

    public static class Item {
        public final String label;
        public final int totalLines;
        public final Color color;
        public Item(String label, int totalLines, Color color) {
            this.label = label;
            this.totalLines = totalLines;
            this.color = color;
        }
    }

    private final List<Item> items;
    private final int padding = 40;
    private final int axisPadding = 50;
    private final ChartMode mode;

    public CodeStatsChartPanel(List<Item> items) {
        this(items, ChartMode.BAR);
    }

    public CodeStatsChartPanel(List<Item> items, ChartMode mode) {
        this.items = items;
        this.mode = mode == null ? ChartMode.BAR : mode;
        setBackground(Color.WHITE);
        setPreferredSize(new Dimension(800, 400));
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);
        Graphics2D g2 = (Graphics2D) g;
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);

        if (mode == ChartMode.PIE) {
            paintPieChart(g2);
        } else {
            paintBarChart(g2);
        }
    }

    private void paintBarChart(Graphics2D g2) {
        int width = getWidth();
        int height = getHeight();

        int chartX = axisPadding;
        int chartY = padding;
        int chartW = width - axisPadding - padding;
        int chartH = height - padding * 2;

        int maxVal = 0;
        for (Item it : items) maxVal = Math.max(maxVal, it.totalLines);
        if (maxVal == 0) maxVal = 1;

        g2.setColor(Color.DARK_GRAY);
        g2.drawLine(chartX, chartY, chartX, chartY + chartH);
        g2.drawLine(chartX, chartY + chartH, chartX + chartW, chartY + chartH);

        g2.setFont(new Font("SansSerif", Font.PLAIN, 12));
        int ticks = 5;
        for (int i = 0; i <= ticks; i++) {
            int val = maxVal * i / ticks;
            int y = chartY + chartH - (int) (chartH * (val / (double) maxVal));
            g2.setColor(new Color(230, 230, 230));
            g2.drawLine(chartX, y, chartX + chartW, y);
            g2.setColor(Color.DARK_GRAY);
            String label = String.valueOf(val);
            int sw = g2.getFontMetrics().stringWidth(label);
            g2.drawString(label, chartX - sw - 8, y + 4);
        }

        int n = Math.max(1, items.size());
        int gap = 18;
        int barW = Math.max(20, (chartW - gap * (n + 1)) / n);
        int x = chartX + gap;

        for (Item it : items) {
            int barH = (int) Math.round(chartH * (it.totalLines / (double) maxVal));
            int y = chartY + chartH - barH;
            g2.setColor(it.color);
            g2.fillRoundRect(x, y, barW, barH, 8, 8);
            g2.setColor(it.color.darker());
            g2.drawRoundRect(x, y, barW, barH, 8, 8);

            g2.setColor(Color.DARK_GRAY);
            String lbl = it.label;
            int sw = g2.getFontMetrics().stringWidth(lbl);
            int lx = x + (barW - sw) / 2;
            int ly = chartY + chartH + 18;
            g2.drawString(lbl, lx, ly);

            String valStr = String.valueOf(it.totalLines);
            int vsw = g2.getFontMetrics().stringWidth(valStr);
            g2.drawString(valStr, x + (barW - vsw) / 2, y - 6);

            x += barW + gap;
        }
    }

    private void paintPieChart(Graphics2D g2) {
        int width = getWidth();
        int height = getHeight();
        int diameter = Math.min(width, height) - padding * 2;
        int x = (width - diameter) / 2;
        int y = padding;
        int total = items.stream().mapToInt(it -> Math.max(0, it.totalLines)).sum();
        if (total == 0) total = 1;

        double startAngle = 0;
        for (Item item : items) {
            double sweep = 360.0 * item.totalLines / total;
            g2.setColor(item.color);
            g2.fillArc(x, y, diameter, diameter, (int) Math.round(startAngle), (int) Math.round(sweep));
            startAngle += sweep;
        }

        int legendX = padding;
        int legendY = y + diameter + 20;
        g2.setFont(new Font("SansSerif", Font.PLAIN, 13));
        for (Item item : items) {
            g2.setColor(item.color);
            g2.fillRect(legendX, legendY, 20, 12);
            g2.setColor(Color.DARK_GRAY);
            g2.drawRect(legendX, legendY, 20, 12);
            g2.drawString(item.label + " (" + item.totalLines + ")", legendX + 28, legendY + 11);
            legendY += 18;
        }
    }
}
