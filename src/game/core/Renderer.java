package game.core;

import model.redpacket.Player;
import model.redpacket.RedPacket;
import model.redpacket.RedPacketShape;

import java.util.List;

public class Renderer {
    private final GameConfig cfg;

    public Renderer(GameConfig cfg) {
        this.cfg = cfg;
    }

    public void clear() {
        // ANSI 清屏并移动到左上角
        System.out.print("\u001b[2J\u001b[H");
    }

    public void drawHUD(int collected, double amount) {
        System.out.printf("红包数: %d  金额: %.2f  (WSAD/方向键移动, Q 退出)\n", collected, amount);
    }

    public void drawWorld(Player player, List<RedPacket> rps) {
        char[][] buf = new char[cfg.height][cfg.width];
        for (int y = 0; y < cfg.height; y++) {
            for (int x = 0; x < cfg.width; x++) {
                buf[y][x] = ' ';
            }
        }
        // 绘制红包
        for (RedPacket rp : rps) {
            if (rp.collected) continue;
            drawRedPacket(buf, rp);
        }
        // 绘制玩家
        int px = (int)Math.round(player.pos.x);
        int py = (int)Math.round(player.pos.y);
        if (px >= 0 && px < cfg.width && py >= 0 && py < cfg.height) {
            buf[py][px] = '@';
        }

        // 输出画面
        for (int y = 0; y < cfg.height; y++) {
            System.out.println(new String(buf[y]));
        }
    }

    /**
     * 绘制单个红包，根据形状和大小显示不同的符号
     */
    private void drawRedPacket(char[][] buf, RedPacket rp) {
        int centerX = (int)Math.round(rp.pos.x);
        int centerY = (int)Math.round(rp.pos.y);
        int radius = (int)Math.ceil(rp.radius);
        
        // 根据大小调整绘制范围
        for (int dy = -radius; dy <= radius; dy++) {
            for (int dx = -radius; dx <= radius; dx++) {
                int x = centerX + dx;
                int y = centerY + dy;
                
                if (x >= 0 && x < cfg.width && y >= 0 && y < cfg.height) {
                    // 根据形状决定是否绘制该点
                    if (shouldDrawPoint(rp.shape, dx, dy, radius)) {
                        buf[y][x] = rp.getDisplaySymbol().charAt(0);
                    }
                }
            }
        }
    }

    /**
     * 根据形状判断是否应该绘制该点
     */
    private boolean shouldDrawPoint(RedPacketShape shape, int dx, int dy, int radius) {
        double distance = Math.sqrt(dx * dx + dy * dy);
        
        switch (shape) {
            case CIRCLE:
                return distance <= radius;
            case SQUARE:
                return Math.abs(dx) <= radius && Math.abs(dy) <= radius;
            case TRIANGLE:
                return dy >= -radius && dy <= 0 && Math.abs(dx) <= (radius + dy);
            case DIAMOND:
                return Math.abs(dx) + Math.abs(dy) <= radius;
            case HEART:
                return (dx * dx + (dy + radius/2) * (dy + radius/2)) <= radius * radius ||
                       (dx * dx + (dy - radius/2) * (dy - radius/2)) <= radius * radius;
            case STAR:
                return distance <= radius && (dx == 0 || dy == 0 || Math.abs(dx) == Math.abs(dy));
            default:
                return distance <= radius;
        }
    }
}
