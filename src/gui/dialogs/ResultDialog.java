package gui;

import model.redpacket.RedPacketStatistics;
import javax.swing.*;
import java.awt.*;

public class ResultDialog extends JDialog {
    private boolean retry = false;

    public ResultDialog(Frame owner, int count, double amount, RedPacketStatistics statistics) {
        super(owner, "Game Result", true);
        setLayout(new BorderLayout(16, 16));

        Font font = pickUiFont(24f);

        JLabel title = new JLabel("Game Over", SwingConstants.CENTER);
        title.setFont(font.deriveFont(Font.BOLD, 34f));
        title.setBorder(BorderFactory.createEmptyBorder(20, 20, 0, 20));
        add(title, BorderLayout.NORTH);

        JPanel center = new JPanel();
        center.setLayout(new BorderLayout());
        center.setBorder(BorderFactory.createEmptyBorder(10, 40, 10, 40));
        
        // 基本统计
        JPanel basicStats = new JPanel();
        basicStats.setLayout(new GridLayout(0, 1, 8, 8));
        JLabel l1 = new JLabel(String.format("Collected: %d", count), SwingConstants.CENTER);
        JLabel l2 = new JLabel(String.format("Total Amount: %.2f", amount), SwingConstants.CENTER);
        l1.setFont(font);
        l2.setFont(font);
        basicStats.add(l1);
        basicStats.add(l2);
        
        // 详细统计
        JPanel detailedStats = createDetailedStatsPanel(statistics, font);
        
        center.add(basicStats, BorderLayout.NORTH);
        center.add(detailedStats, BorderLayout.CENTER);
        add(center, BorderLayout.CENTER);

        JPanel btns = new JPanel(new FlowLayout(FlowLayout.CENTER, 24, 16));
        JButton again = new JButton("Play Again");
        JButton close = new JButton("Close");
        again.setFont(font.deriveFont(Font.BOLD, 26f));
        close.setFont(font);
        btns.add(again);
        btns.add(close);
        add(btns, BorderLayout.SOUTH);

        again.addActionListener(e -> { retry = true; setVisible(false); });
        close.addActionListener(e -> { retry = false; setVisible(false); });

        // Size: aim roughly 2x of a typical option pane
        setResizable(true);
        pack();
        Dimension d = getSize();
        int w = Math.max(800, d.width * 2);
        int h = Math.max(600, d.height * 2);
        setSize(w, h);
        setLocationRelativeTo(owner);
    }

    public boolean isRetry() {
        return retry;
    }

    public static boolean show(Frame owner, int count, double amount, RedPacketStatistics statistics) {
        ResultDialog rd = new ResultDialog(owner, count, amount, statistics);
        rd.setVisible(true);
        return rd.isRetry();
    }
    
    public static boolean show(Frame owner, int count, double amount) {
        return show(owner, count, amount, null);
    }

    private JPanel createDetailedStatsPanel(RedPacketStatistics statistics, Font font) {
        JPanel panel = new JPanel();
        panel.setLayout(new BorderLayout());
        
        if (statistics == null || statistics.getTotalCount() == 0) {
            JLabel noStats = new JLabel("No detailed statistics available", SwingConstants.CENTER);
            noStats.setFont(font.deriveFont(16f));
            panel.add(noStats, BorderLayout.CENTER);
            return panel;
        }
        
        JTabbedPane tabbedPane = new JTabbedPane();
        tabbedPane.setFont(font.deriveFont(14f));
        
        // 按形状统计
        JPanel shapePanel = new JPanel();
        shapePanel.setLayout(new GridLayout(0, 1, 4, 4));
        shapePanel.setBorder(BorderFactory.createEmptyBorder(10, 20, 10, 20));
        
        JLabel shapeTitle = new JLabel("By Shape:", SwingConstants.LEFT);
        shapeTitle.setFont(font.deriveFont(Font.BOLD, 16f));
        shapePanel.add(shapeTitle);
        
        statistics.getShapeAmounts().forEach((shape, amount) -> {
            int count = statistics.getShapeCounts().get(shape);
            JLabel label = new JLabel(String.format("%s: %d pcs, Amount: %.2f", shape, count, amount));
            label.setFont(font.deriveFont(14f));
            shapePanel.add(label);
        });
        
        // 按大小统计
        JPanel sizePanel = new JPanel();
        sizePanel.setLayout(new GridLayout(0, 1, 4, 4));
        sizePanel.setBorder(BorderFactory.createEmptyBorder(10, 20, 10, 20));
        
        JLabel sizeTitle = new JLabel("By Size:", SwingConstants.LEFT);
        sizeTitle.setFont(font.deriveFont(Font.BOLD, 16f));
        sizePanel.add(sizeTitle);
        
        statistics.getSizeAmounts().forEach((size, amount) -> {
            int count = statistics.getSizeCounts().get(size);
            JLabel label = new JLabel(String.format("%s: %d pcs, Amount: %.2f", size, count, amount));
            label.setFont(font.deriveFont(14f));
            sizePanel.add(label);
        });
        
        // 按形状+大小组合统计
        JPanel comboPanel = new JPanel();
        comboPanel.setLayout(new GridLayout(0, 1, 4, 4));
        comboPanel.setBorder(BorderFactory.createEmptyBorder(10, 20, 10, 20));
        
        JLabel comboTitle = new JLabel("Shape + Size Combination:", SwingConstants.LEFT);
        comboTitle.setFont(font.deriveFont(Font.BOLD, 16f));
        comboPanel.add(comboTitle);
        
        statistics.getShapeSizeAmounts().forEach((combo, amount) -> {
            int count = statistics.getShapeSizeCounts().get(combo);
            JLabel label = new JLabel(String.format("%s: %d pcs, Amount: %.2f", combo, count, amount));
            label.setFont(font.deriveFont(14f));
            comboPanel.add(label);
        });
        
        tabbedPane.addTab("Shapes", shapePanel);
        tabbedPane.addTab("Sizes", sizePanel);
        tabbedPane.addTab("Combinations", comboPanel);
        
        panel.add(tabbedPane, BorderLayout.CENTER);
        return panel;
    }

    private static Font pickUiFont(float size) {
        String[] candidates = new String[]{
                "DejaVu Sans", "Arial", "Liberation Sans", "Helvetica", Font.SANS_SERIF
        };
        for (String name : candidates) {
            try {
                return new Font(name, Font.PLAIN, Math.round(size));
            } catch (Exception ignored) {}
        }
        return new Font(Font.SANS_SERIF, Font.PLAIN, Math.round(size));
    }
}
