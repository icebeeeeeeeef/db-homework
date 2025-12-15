package gui;

import javax.swing.*;
import java.awt.*;

public class StartDialog extends JDialog {
    private final JTextField tfWidth = new JTextField("80");
    private final JTextField tfHeight = new JTextField("45");
    private final JTextField tfCount = new JTextField("120");
    private final JTextField tfDuration = new JTextField("20000");
    private final JTextField tfFps = new JTextField("60");
    private final JTextField tfRadius = new JTextField("1.6");
    private boolean confirmed = false;

    public StartDialog(Frame owner) {
        super(owner, "Game Configuration", true);
        setLayout(new BorderLayout(16, 16));

        // Use a widely available Latin UI font to avoid missing glyphs/boxes
        Font font = pickUiFont(24f);

        // 设置输入框列宽
        tfWidth.setColumns(14);
        tfHeight.setColumns(14);
        tfCount.setColumns(14);
        tfDuration.setColumns(14);
        tfFps.setColumns(14);
        tfRadius.setColumns(14);

        JPanel head = new JPanel(new BorderLayout());
        JLabel title = new JLabel("Configure game parameters, then click Start", SwingConstants.LEFT);
        title.setFont(font.deriveFont(Font.BOLD, 28f));
        title.setBorder(BorderFactory.createEmptyBorder(16, 24, 0, 24));
        head.add(title, BorderLayout.CENTER);
        add(head, BorderLayout.NORTH);

        JPanel form = new JPanel(new GridLayout(0, 2, 12, 12));
        form.setBorder(BorderFactory.createEmptyBorder(16, 24, 8, 24));
        JLabel l1 = new JLabel("Width (cells):"); l1.setFont(font); form.add(l1); tfWidth.setFont(font); form.add(tfWidth);
        JLabel l2 = new JLabel("Height (cells):"); l2.setFont(font); form.add(l2); tfHeight.setFont(font); form.add(tfHeight);
        JLabel l3 = new JLabel("Red packets count:"); l3.setFont(font); form.add(l3); tfCount.setFont(font); form.add(tfCount);
        JLabel l4 = new JLabel("Duration (ms):"); l4.setFont(font); form.add(l4); tfDuration.setFont(font); form.add(tfDuration);
        JLabel l5 = new JLabel("FPS:"); l5.setFont(font); form.add(l5); tfFps.setFont(font); form.add(tfFps);
        JLabel l6 = new JLabel("Player radius (cells):"); l6.setFont(font); form.add(l6); tfRadius.setFont(font); form.add(tfRadius);
        add(form, BorderLayout.CENTER);

        JPanel btns = new JPanel(new FlowLayout(FlowLayout.RIGHT, 16, 16));
        JButton ok = new JButton("Start");
        JButton cancel = new JButton("Cancel");
        ok.setToolTipText("Start game");
        cancel.setToolTipText("Close dialog");
        ok.setMnemonic('K');
        cancel.setMnemonic('Q');
        ok.setFont(font.deriveFont(Font.BOLD, 32f));
        cancel.setFont(font);
        btns.add(ok); btns.add(cancel);
        add(btns, BorderLayout.SOUTH);

        ok.addActionListener(e -> { confirmed = true; setVisible(false); });
        cancel.addActionListener(e -> { confirmed = false; setVisible(false); });

        // 将“开始游戏”设为默认按钮，回车可触发
        getRootPane().setDefaultButton(ok);

        // Force a large dialog size (~80% of screen). Use setSize to avoid pack shrinking.
        Dimension screen = Toolkit.getDefaultToolkit().getScreenSize();
        int w = (int) (screen.width * 0.8);
        int h = (int) (screen.height * 0.8);
        setPreferredSize(new Dimension(w, h));
        setMinimumSize(new Dimension(Math.min(w, 1200), Math.min(h, 800)));
        setResizable(true);
        pack();
        setSize(w, h);
        setLocationRelativeTo(owner);
    }

    private static Font pickUiFont(float size) {
        String[] candidates = new String[]{
                "DejaVu Sans", "Arial", "Liberation Sans", "Helvetica", Font.SANS_SERIF
        };
        GraphicsEnvironment ge = GraphicsEnvironment.getLocalGraphicsEnvironment();
        for (String name : candidates) {
            try {
                Font f = new Font(name, Font.PLAIN, Math.round(size));
                return f;
            } catch (Exception ignored) { }
        }
        return new Font(Font.SANS_SERIF, Font.PLAIN, Math.round(size));
    }

    public boolean isConfirmed() { return confirmed; }

    public int getWorldWidth() { return Integer.parseInt(tfWidth.getText().trim()); }
    public int getWorldHeight() { return Integer.parseInt(tfHeight.getText().trim()); }
    public int getCount() { return Integer.parseInt(tfCount.getText().trim()); }
    public long getDuration() { return Long.parseLong(tfDuration.getText().trim()); }
    public int getFps() { return Integer.parseInt(tfFps.getText().trim()); }
    public double getRadius() { return Double.parseDouble(tfRadius.getText().trim()); }
}

