package game.core;

import model.redpacket.Player;
import model.redpacket.RedPacket;
import model.redpacket.RedPacketStatistics;

import java.util.ArrayList;
import java.util.List;

public class GameEngine {
    private final GameConfig config;
    private final List<RedPacket> redPackets = new ArrayList<>();
    private final Player player;
    private int collectedCount = 0;
    private double collectedAmount = 0.0;
    private final RedPacketStatistics statistics = new RedPacketStatistics();
    private final Renderer renderer;
    private final InputController input;

    public GameEngine(GameConfig config) {
        this.config = config;
        this.player = new Player(config.width / 2.0, config.height / 2.0, config.playerRadius);
        initRedPackets();
        this.renderer = new Renderer(config);
        this.input = new InputController();
    }

    private void initRedPackets() {
        for (int i = 0; i < config.redPacketCount; i++) {
            redPackets.add(new RedPacket(i, 0.1, 10.0, config.width, config.height));
        }
    }

    public void start() {
        Thread t = new Thread(input, "input-thread");
        t.setDaemon(true);
        t.start();

        long endAt = System.currentTimeMillis() + config.durationMillis;
        while (System.currentTimeMillis() < endAt && !input.isQuit()) {
            updatePlayerByInput();
            updateRedPackets();
            detectCollisions();
            render();
            sleep(config.tickMillis);
        }
        input.stop();
    }

    private void updatePlayerByInput() {
        int dx = input.getDx();
        int dy = input.getDy();
        double speed = 1.5;
        player.pos.x = clamp(player.pos.x + dx * speed, 0, config.width - 1);
        player.pos.y = clamp(player.pos.y + dy * speed, 0, config.height - 1);
    }

    private void updateRedPackets() {
        for (RedPacket rp : redPackets) {
            if (rp.collected) continue;
            rp.pos.x += rp.vel.x;
            rp.pos.y += rp.vel.y;
            if (rp.pos.x < 0 || rp.pos.x > config.width) rp.vel.x = -rp.vel.x;
            if (rp.pos.y < 0 || rp.pos.y > config.height) rp.vel.y = -rp.vel.y;
            rp.pos.x = clamp(rp.pos.x, 0, config.width);
            rp.pos.y = clamp(rp.pos.y, 0, config.height);
        }
    }

    private void detectCollisions() {
        for (RedPacket rp : redPackets) {
            if (player.collide(rp)) {
                rp.collected = true;
                collectedCount++;
                collectedAmount += rp.amount;
                statistics.recordCollected(rp);
            }
        }
    }

    private void render() {
        renderer.clear();
        renderer.drawHUD(collectedCount, collectedAmount);
        renderer.drawWorld(player, redPackets);
    }

    private double clamp(double v, double min, double max) {
        if (v < min) return min;
        if (v > max) return max;
        return v;
    }

    private void sleep(long ms) {
        try {
            Thread.sleep(ms);
        } catch (InterruptedException ignored) {
        }
    }

    public void printSummary() {
        System.out.printf("时长: %dms, 红包总数: %d, 撞到: %d, 金额: %.2f\n",
                config.durationMillis, config.redPacketCount, collectedCount, collectedAmount);
        
        System.out.println("\n=== 红包统计详情 ===");
        
        // 按形状统计
        System.out.println("\n按形状统计:");
        statistics.getShapeAmounts().forEach((shape, amount) -> {
            int count = statistics.getShapeCounts().get(shape);
            System.out.printf("  %s: %d个, 金额: %.2f\n", shape, count, amount);
        });
        
        // 按大小统计
        System.out.println("\n按大小统计:");
        statistics.getSizeAmounts().forEach((size, amount) -> {
            int count = statistics.getSizeCounts().get(size);
            System.out.printf("  %s: %d个, 金额: %.2f\n", size, count, amount);
        });
        
        // 按形状+大小组合统计
        System.out.println("\n按形状+大小组合统计:");
        statistics.getShapeSizeAmounts().forEach((combo, amount) -> {
            int count = statistics.getShapeSizeCounts().get(combo);
            System.out.printf("  %s: %d个, 金额: %.2f\n", combo, count, amount);
        });
    }
}

