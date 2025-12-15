package gui;

import model.redpacket.Player;
import model.redpacket.RedPacket;
import model.redpacket.RedPacketStatistics;
import geom.Vector2;

import javax.swing.*;
import java.awt.*;
import java.awt.event.KeyAdapter;
import java.awt.event.KeyEvent;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ThreadLocalRandom;
import java.awt.geom.Point2D;

public class GuiGame extends JPanel {
    private final int width;
    private final int height;
    private final int redPacketCount;
    private final double playerRadius;
    private final long durationMillis;
    private final long startAt;

    private final List<RedPacket> redPackets = new ArrayList<>();
    private final Player player;
    private int collectedCount = 0;
    private double collectedAmount = 0.0;
    private final RedPacketStatistics statistics = new RedPacketStatistics();

    private long endAt;
    private final Timer timer;
    private int dx = 0;
    private int dy = 0;
    private GameOverListener gameOverListener;
    private final GameScene scene;
    private final List<Point2D.Double> decorations = new ArrayList<>();

    public GuiGame(int width, int height, int redPacketCount, double playerRadius, long durationMillis, int fps) {
        this.width = width;
        this.height = height;
        this.redPacketCount = redPacketCount;
        this.playerRadius = playerRadius;
        this.durationMillis = durationMillis;
        // 放大画布尺寸
        setPreferredSize(new Dimension(Math.max(1000, width * 12), Math.max(700, height * 18)));
        setBackground(Color.BLACK);
        setFocusable(true);
        requestFocusInWindow();

        this.player = new Player(width / 2.0, height - 2.0, playerRadius);
        initRedPackets();
        this.scene = GameScene.randomScene();
        this.startAt = System.currentTimeMillis();
        initDecorations();
        this.endAt = System.currentTimeMillis() + durationMillis;

        addKeyListener(new KeyAdapter() {
            @Override
            public void keyPressed(KeyEvent e) {
                switch (e.getKeyCode()) {
                    case KeyEvent.VK_LEFT: case KeyEvent.VK_A: dx = -1; break;
                    case KeyEvent.VK_RIGHT: case KeyEvent.VK_D: dx = 1; break;
                    case KeyEvent.VK_UP: case KeyEvent.VK_W: dy = -1; break;
                    case KeyEvent.VK_DOWN: case KeyEvent.VK_S: dy = 1; break;
                }
            }

            @Override
            public void keyReleased(KeyEvent e) {
                switch (e.getKeyCode()) {
                    case KeyEvent.VK_LEFT: case KeyEvent.VK_A: if (dx < 0) dx = 0; break;
                    case KeyEvent.VK_RIGHT: case KeyEvent.VK_D: if (dx > 0) dx = 0; break;
                    case KeyEvent.VK_UP: case KeyEvent.VK_W: if (dy < 0) dy = 0; break;
                    case KeyEvent.VK_DOWN: case KeyEvent.VK_S: if (dy > 0) dy = 0; break;
                }
            }
        });

        int delay = Math.max(5, 1000 / Math.max(1, fps));
        this.timer = new Timer(delay, e -> {
            step();
            repaint();
            if (System.currentTimeMillis() >= endAt) {
                ((Timer) e.getSource()).stop();
                if (gameOverListener != null) {
                    SwingUtilities.invokeLater(() -> gameOverListener.onGameOver(collectedCount, collectedAmount));
                }
            }
        });
    }

    public void start() {
        timer.start();
    }

    public void setGameOverListener(GameOverListener listener) {
        this.gameOverListener = listener;
    }

    public RedPacketStatistics getStatistics() {
        return statistics;
    }

    private void initRedPackets() {
        for (int i = 0; i < redPacketCount; i++) {
            RedPacket rp = new RedPacket(i, 0.1, 10.0, width, height);
            // 下落初始速度
            rp.vel = new Vector2(0, ThreadLocalRandom.current().nextDouble(0.1, 0.6));
            // 出生在顶部
            rp.pos = new Vector2(ThreadLocalRandom.current().nextDouble(0, width), 0);
            redPackets.add(rp);
        }
    }

    private void step() {
        // 玩家移动（限制边界）
        double speed = 0.8;
        player.pos.x = clamp(player.pos.x + dx * speed, 0, width - 1);
        player.pos.y = clamp(player.pos.y + dy * speed, 0, height - 1);

        // 红包下落，越界则重生
        for (RedPacket rp : redPackets) {
            if (rp.collected) continue;
            rp.pos.x += rp.vel.x;
            rp.pos.y += rp.vel.y;
            if (rp.pos.y > height) {
                // 重生
                rp.pos.y = 0;
                rp.pos.x = ThreadLocalRandom.current().nextDouble(0, width);
                rp.collected = false;
            }
        }

        // 碰撞检测
        for (RedPacket rp : redPackets) {
            if (!rp.collected && player.collide(rp)) {
                rp.collected = true;
                collectedCount++;
                collectedAmount += rp.amount;
                statistics.recordCollected(rp);
            }
        }
    }

    private double clamp(double v, double min, double max) {
        if (v < min) return min;
        if (v > max) return max;
        return v;
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);
        Graphics2D g2 = (Graphics2D) g;
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);

        paintSceneBackground(g2);

        int cellW = Math.max(12, getWidth() / width);
        int cellH = Math.max(18, getHeight() / height);

        // 背景网格
        g2.setColor(new Color(30, 30, 30));
        for (int x = 0; x < width; x++) g2.drawLine(x * cellW, 0, x * cellW, getHeight());
        for (int y = 0; y < height; y++) g2.drawLine(0, y * cellH, getWidth(), y * cellH);

        drawDecorations(g2);

        // HUD
        g2.setColor(Color.WHITE);
        double remainingSec = Math.max(0, (endAt - System.currentTimeMillis()) / 1000.0);
        g2.drawString(String.format("Scene: %s | Time Left: %.1fs | Captured: %d | Amount: %.2f",
                scene.displayName, remainingSec, collectedCount, collectedAmount), 10, 20);

        // 红包
        for (RedPacket rp : redPackets) {
            if (rp.collected) continue;
            drawRedPacket(g2, rp, cellW, cellH);
        }

        // 玩家
        int px = (int) Math.round(player.pos.x * cellW);
        int py = (int) Math.round(player.pos.y * cellH);
        g2.setColor(Color.CYAN);
        int pSize = Math.max(20, (int)(Math.min(cellW, cellH) * 1.4)); // 玩家比红包大
        g2.fillRoundRect(px, py, pSize, pSize, pSize/3, pSize/3);
    }

    /**
     * 绘制单个红包，根据形状和大小显示不同的样式
     */
    private void drawRedPacket(Graphics2D g2, RedPacket rp, int cellW, int cellH) {
        int x = (int) Math.round(rp.pos.x * cellW);
        int y = (int) Math.round(rp.pos.y * cellH);
        int baseSize = Math.max(14, Math.min(cellW, cellH));
        int rpSize = (int) (baseSize * rp.radius);
        
        // 根据大小选择颜色
        Color[] sizeColors = {
            new Color(255, 100, 100), // 小 - 浅红
            new Color(255, 50, 50),   // 中 - 中红
            new Color(200, 0, 0),     // 大 - 深红
            new Color(150, 0, 0)      // 超大 - 暗红
        };
        Color[] shapeColors = {
            new Color(255, 215, 0),   // 金色
            new Color(255, 165, 0),   // 橙色
            new Color(255, 20, 147),  // 深粉红
            new Color(138, 43, 226),  // 蓝紫色
            new Color(255, 69, 0),    // 红橙色
            new Color(255, 215, 0)    // 金色
        };
        
        Color bgColor = sizeColors[rp.size.ordinal()];
        Color borderColor = shapeColors[rp.shape.ordinal()];
        
        g2.setColor(bgColor);
        
        // 根据形状绘制不同的图形
        switch (rp.shape) {
            case CIRCLE:
                g2.fillOval(x, y, rpSize, rpSize);
                break;
            case SQUARE:
                g2.fillRect(x, y, rpSize, rpSize);
                break;
            case TRIANGLE:
                int[] xPoints = {x + rpSize/2, x, x + rpSize};
                int[] yPoints = {y, y + rpSize, y + rpSize};
                g2.fillPolygon(xPoints, yPoints, 3);
                break;
            case DIAMOND:
                int[] dxPoints = {x + rpSize/2, x + rpSize, x + rpSize/2, x};
                int[] dyPoints = {y, y + rpSize/2, y + rpSize, y + rpSize/2};
                g2.fillPolygon(dxPoints, dyPoints, 4);
                break;
            case HEART:
                // 简化的心形
                g2.fillOval(x, y + rpSize/4, rpSize/2, rpSize/2);
                g2.fillOval(x + rpSize/2, y + rpSize/4, rpSize/2, rpSize/2);
                g2.fillPolygon(new int[]{x, x + rpSize/2, x + rpSize}, 
                             new int[]{y + rpSize/2, y + rpSize, y + rpSize/2}, 3);
                break;
            case STAR:
                // 简化的星形
                int[] sxPoints = {x + rpSize/2, x + rpSize*3/8, x, x + rpSize/4, x + rpSize/2, 
                                 x + rpSize*3/4, x + rpSize, x + rpSize*5/8, x + rpSize/2};
                int[] syPoints = {y, y + rpSize/3, y + rpSize/2, y + rpSize/3, y + rpSize, 
                                 y + rpSize/3, y + rpSize/2, y + rpSize/3, y};
                g2.fillPolygon(sxPoints, syPoints, 9);
                break;
        }
        
        // 绘制边框
        g2.setColor(borderColor);
        g2.setStroke(new BasicStroke(2));
        switch (rp.shape) {
            case CIRCLE:
                g2.drawOval(x, y, rpSize, rpSize);
                break;
            case SQUARE:
                g2.drawRect(x, y, rpSize, rpSize);
                break;
            case TRIANGLE:
                int[] xPoints = {x + rpSize/2, x, x + rpSize};
                int[] yPoints = {y, y + rpSize, y + rpSize};
                g2.drawPolygon(xPoints, yPoints, 3);
                break;
            case DIAMOND:
                int[] dxPoints = {x + rpSize/2, x + rpSize, x + rpSize/2, x};
                int[] dyPoints = {y, y + rpSize/2, y + rpSize, y + rpSize/2};
                g2.drawPolygon(dxPoints, dyPoints, 4);
                break;
            case HEART:
                g2.drawOval(x, y + rpSize/4, rpSize/2, rpSize/2);
                g2.drawOval(x + rpSize/2, y + rpSize/4, rpSize/2, rpSize/2);
                g2.drawPolygon(new int[]{x, x + rpSize/2, x + rpSize}, 
                             new int[]{y + rpSize/2, y + rpSize, y + rpSize/2}, 3);
                break;
            case STAR:
                int[] sxPoints = {x + rpSize/2, x + rpSize*3/8, x, x + rpSize/4, x + rpSize/2, 
                                 x + rpSize*3/4, x + rpSize, x + rpSize*5/8, x + rpSize/2};
                int[] syPoints = {y, y + rpSize/3, y + rpSize/2, y + rpSize/3, y + rpSize, 
                                 y + rpSize/3, y + rpSize/2, y + rpSize/3, y};
                g2.drawPolygon(sxPoints, syPoints, 9);
                break;
        }
        
        // 绘制金额符号
        g2.setColor(Color.YELLOW);
        g2.setFont(new Font("Arial", Font.BOLD, rpSize/3));
        FontMetrics fm = g2.getFontMetrics();
        String symbol = rp.getDisplaySymbol();
        int textX = x + (rpSize - fm.stringWidth(symbol)) / 2;
        int textY = y + (rpSize + fm.getAscent()) / 2;
        g2.drawString(symbol, textX, textY);
    }

    private void initDecorations() {
        int count = scene.decorationCount;
        for (int i = 0; i < count; i++) {
            decorations.add(new Point2D.Double(Math.random(), Math.random()));
        }
    }

    private void paintSceneBackground(Graphics2D g2) {
        GradientPaint paint = new GradientPaint(0, 0, scene.topColor, 0, getHeight(), scene.bottomColor);
        g2.setPaint(paint);
        g2.fillRect(0, 0, getWidth(), getHeight());
    }

    private void drawDecorations(Graphics2D g2) {
        g2.setColor(scene.decorationColor);
        for (Point2D.Double pt : decorations) {
            int x = (int) (pt.x * getWidth());
            int y = (int) (pt.y * getHeight());
            switch (scene.decorationType) {
                case SNOW:
                    g2.fillOval(x, y, 6, 6);
                    break;
                case STARS:
                    g2.drawLine(x, y, x, y);
                    g2.drawLine(x - 2, y, x + 2, y);
                    g2.drawLine(x, y - 2, x, y + 2);
                    break;
                case BLOSSOM:
                    g2.fillOval(x, y, 8, 8);
                    g2.fillOval(x + 4, y + 4, 8, 8);
                    break;
            }
        }
    }

    private enum DecorationType {
        SNOW, STARS, BLOSSOM
    }

    private enum GameScene {
        SUNSET("Sunset Beach", new Color(255, 153, 102), new Color(255, 94, 98), DecorationType.BLOSSOM, new Color(255, 255, 255, 160), 40),
        SNOW("Snowfield", new Color(120, 180, 255), new Color(200, 220, 255), DecorationType.SNOW, Color.WHITE, 80),
        SPACE("Cosmic Drift", new Color(10, 10, 45), new Color(3, 3, 20), DecorationType.STARS, new Color(255, 255, 200), 70);

        final String displayName;
        final Color topColor;
        final Color bottomColor;
        final DecorationType decorationType;
        final Color decorationColor;
        final int decorationCount;

        GameScene(String displayName, Color topColor, Color bottomColor,
                  DecorationType decorationType, Color decorationColor, int decorationCount) {
            this.displayName = displayName;
            this.topColor = topColor;
            this.bottomColor = bottomColor;
            this.decorationType = decorationType;
            this.decorationColor = decorationColor;
            this.decorationCount = decorationCount;
        }

        static GameScene randomScene() {
            GameScene[] scenes = values();
            return scenes[ThreadLocalRandom.current().nextInt(scenes.length)];
        }
    }

    public static void launchInteractive() {
        SwingUtilities.invokeLater(() -> {
            JFrame frame = new JFrame("Red Packet Game");
            StartDialog dialog = new StartDialog(frame);
            dialog.setVisible(true);
            if (!dialog.isConfirmed()) {
                frame.dispose();
                return;
            }
            frame.setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
            frame.setVisible(true);

            Runnable startNewGame = new Runnable() {
                @Override public void run() {
                    final int worldW, worldH, count, fps; final long duration; final double radius;
                    try {
                        worldW = dialog.getWorldWidth();
                        worldH = dialog.getWorldHeight();
                        count = dialog.getCount();
                        duration = dialog.getDuration();
                        fps = dialog.getFps();
                        radius = dialog.getRadius();
                    } catch (Exception ex) {
                        JOptionPane.showMessageDialog(frame, "Invalid parameters, please check and try again.", "Error", JOptionPane.ERROR_MESSAGE);
                        StartDialog d2 = new StartDialog(frame);
                        d2.setVisible(true);
                        if (!d2.isConfirmed()) { frame.dispose(); return; }
                        SwingUtilities.invokeLater(this);
                        return;
                    }

                    GuiGame panel = new GuiGame(worldW, worldH, count, radius, duration, fps);
                    panel.setGameOverListener((c, amt) -> {
                        boolean again = ResultDialog.show(frame, c, amt, panel.statistics);
                        if (again) {
                            SwingUtilities.invokeLater(this);
                        } else {
                            frame.dispose();
                        }
                    });
                    frame.setContentPane(panel);
                    frame.pack();
                    frame.setLocationRelativeTo(null);
                    frame.revalidate();
                    frame.repaint();
                    panel.start();
                }
            };

            startNewGame.run();
        });
    }
}
