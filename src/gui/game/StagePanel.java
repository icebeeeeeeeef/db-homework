package gui;

import model.duck.DuckCharacter;
import model.duck.customization.DuckOutfit;
import model.duck.customization.DuckAppearance;
import model.duck.behavior.DuckBehaviorProfile;

import javax.swing.*;
import java.awt.*;
import java.awt.event.MouseAdapter;
import java.awt.event.MouseEvent;
import java.util.ArrayList;
import java.util.EnumMap;
import java.util.List;
import java.util.Map;
import java.util.Random;

/**
 * Playground panel rendering Donald Duck and the three ducklings.
 */
public class StagePanel extends JPanel {
    public interface StageListener {
        void onDonaldClicked();
    }

    private final EnumMap<DuckCharacter, DuckOutfit> outfits = new EnumMap<>(DuckCharacter.class);
    private final EnumMap<DuckCharacter, DuckBehaviorProfile> behaviors = new EnumMap<>(DuckCharacter.class);
    private final List<RainDrop> rainDrops = new ArrayList<>();
    private StageListener listener;
    private String speech = "Tap Donald to issue a command";
    private Rectangle donaldBounds = new Rectangle();
    private javax.swing.Timer rainTimer;

    public StagePanel() {
        setOpaque(false);
        setPreferredSize(new Dimension(360, 320));
        setMinimumSize(new Dimension(320, 280));
        setBackground(new Color(245, 250, 255));

        addMouseListener(new MouseAdapter() {
            @Override
            public void mouseClicked(MouseEvent e) {
                if (donaldBounds.contains(e.getPoint()) && listener != null) {
                    listener.onDonaldClicked();
                }
            }
        });
    }

    public void setStageListener(StageListener listener) {
        this.listener = listener;
    }

    public void setSpeech(String speech) {
        this.speech = speech == null ? "" : speech;
        repaint();
    }

    public void setOutfit(DuckCharacter character, DuckOutfit outfit) {
        DuckOutfit copy = outfit == null ? new DuckOutfit() : outfit.clone();
        outfits.put(character, copy);
        repaint();
    }

    public DuckOutfit getOutfit(DuckCharacter character) {
        DuckOutfit outfit = outfits.get(character);
        if (outfit == null) {
            outfit = new DuckOutfit();
            outfits.put(character, outfit);
        }
        return outfit;
    }

    public void setBehavior(DuckCharacter character, DuckBehaviorProfile profile) {
        if (profile == null) {
            behaviors.remove(character);
        } else {
            behaviors.put(character, profile.clone());
        }
        repaint();
    }

    public DuckBehaviorProfile getBehavior(DuckCharacter character) {
        DuckBehaviorProfile profile = behaviors.get(character);
        return profile == null ? null : profile.clone();
    }

    public void triggerRain(int durationMillis, int density) {
        rainDrops.clear();
        Random random = new Random();
        int dropCount = Math.max(20, density * 30);
        for (int i = 0; i < dropCount; i++) {
            rainDrops.add(new RainDrop(random.nextInt(Math.max(1, getWidth())), random.nextInt(Math.max(1, getHeight() / 2))));
        }
        if (rainTimer != null) {
            rainTimer.stop();
        }
        rainTimer = new javax.swing.Timer(30, e -> {
            for (RainDrop drop : rainDrops) {
                drop.y += drop.speed;
                if (drop.y > getHeight()) {
                    drop.y = -10;
                    drop.x = random.nextInt(Math.max(1, getWidth()));
                }
            }
            repaint();
        });
        rainTimer.start();
        if (durationMillis > 0) {
            new javax.swing.Timer(durationMillis, e -> {
                rainDrops.clear();
                if (rainTimer != null) {
                    rainTimer.stop();
                }
                repaint();
            }) {{
                setRepeats(false);
            }}.start();
        }
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);
        Graphics2D g2 = (Graphics2D) g;
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);

        paintBackground(g2);

        Map<DuckCharacter, Rectangle> poses = computeCharacterBounds();
        for (Map.Entry<DuckCharacter, Rectangle> entry : poses.entrySet()) {
            DuckCharacter character = entry.getKey();
            Rectangle rect = entry.getValue();
            DuckOutfit outfit = getOutfit(character);
            DuckAppearance appearance = outfit.buildAppearance();
            appearance.paint(g2, rect);

            if (character == DuckCharacter.DONALD) {
                donaldBounds = new Rectangle(rect);
                paintSpeechBubble(g2, rect);
            }

            g2.setColor(new Color(0, 0, 0, 120));
            g2.setFont(g2.getFont().deriveFont(Font.BOLD, 14f));
            String label = character.getDisplayName();
            int sw = g2.getFontMetrics().stringWidth(label);
            int baseY = rect.y + rect.height + 16;
            g2.drawString(label, rect.x + (rect.width - sw) / 2, baseY);

            DuckBehaviorProfile profile = behaviors.get(character);
            if (profile != null) {
                g2.setFont(g2.getFont().deriveFont(Font.PLAIN, 12f));
                String actionName = profile.getActionBehavior() != null ? profile.getActionBehavior().getName() : "No action";
                String soundName = profile.getSoundBehavior() != null ? profile.getSoundBehavior().getName() : "No sound";
                String behaviorLine = actionName + " | " + soundName;
                int bw = g2.getFontMetrics().stringWidth(behaviorLine);
                g2.drawString(behaviorLine, rect.x + (rect.width - bw) / 2, baseY + 14);
            }
        }

        paintRain(g2);
    }

    private void paintBackground(Graphics2D g2) {
        GradientPaint paint = new GradientPaint(0, 0, new Color(173, 216, 255), 0, getHeight(), new Color(255, 255, 255));
        g2.setPaint(paint);
        g2.fillRect(0, 0, getWidth(), getHeight());

        g2.setColor(new Color(180, 220, 240));
        g2.fillRoundRect(20, getHeight() - 80, getWidth() - 40, 60, 40, 40);
    }

    private Map<DuckCharacter, Rectangle> computeCharacterBounds() {
        EnumMap<DuckCharacter, Rectangle> positions = new EnumMap<>(DuckCharacter.class);
        int width = getWidth();
        int height = getHeight();
        int baseW = Math.min(140, Math.max(100, width / 3));
        int baseH = Math.min(160, Math.max(120, height / 2));

        Rectangle donaldRect = new Rectangle(width / 2 - baseW / 2, height / 2 - baseH / 2, baseW, baseH);
        positions.put(DuckCharacter.DONALD, donaldRect);

        int ducklingW = (int) (baseW * 0.6);
        int ducklingH = (int) (baseH * 0.7);

        positions.put(DuckCharacter.DUCKLING_ONE,
                new Rectangle(width / 6 - ducklingW / 2, height - ducklingH - 30, ducklingW, ducklingH));
        positions.put(DuckCharacter.DUCKLING_TWO,
                new Rectangle(width / 2 - ducklingW / 2, height - ducklingH - 10, ducklingW, ducklingH));
        positions.put(DuckCharacter.DUCKLING_THREE,
                new Rectangle(width * 5 / 6 - ducklingW / 2, height - ducklingH - 30, ducklingW, ducklingH));
        return positions;
    }

    private void paintSpeechBubble(Graphics2D g2, Rectangle anchor) {
        if (speech == null || speech.isBlank()) return;
        int padding = 12;
        g2.setFont(g2.getFont().deriveFont(Font.BOLD, 14f));
        int sw = g2.getFontMetrics().stringWidth(speech);
        int bubbleW = sw + padding * 2;
        int bubbleH = 40;
        int x = Math.max(20, anchor.x + anchor.width / 2 - bubbleW / 2);
        int y = Math.max(10, anchor.y - bubbleH - 10);
        g2.setColor(new Color(255, 255, 255, 220));
        g2.fillRoundRect(x, y, bubbleW, bubbleH, 20, 20);
        g2.setColor(new Color(0, 0, 0, 120));
        g2.drawRoundRect(x, y, bubbleW, bubbleH, 20, 20);
        g2.drawString(speech, x + padding, y + bubbleH / 2 + 5);
    }

    private void paintRain(Graphics2D g2) {
        if (rainDrops.isEmpty()) return;
        Color packetColor = new Color(255, 0, 0, 180);
        for (RainDrop drop : rainDrops) {
            g2.setColor(packetColor);
            g2.fillRoundRect((int) drop.x, (int) drop.y, drop.size, drop.size * 2, 6, 6);
            g2.setColor(Color.YELLOW);
            g2.drawLine((int) drop.x + drop.size / 2, (int) drop.y, (int) drop.x + drop.size / 2, (int) drop.y + drop.size * 2);
        }
    }

    private static class RainDrop {
        double x;
        double y;
        final double speed;
        final int size;

        RainDrop(int startX, int startY) {
            this.x = startX;
            this.y = startY;
            this.speed = 4 + Math.random() * 4;
            this.size = 8 + (int) (Math.random() * 6);
        }
    }
}
