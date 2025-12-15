package model.redpacket;

import geom.Vector2;
import java.util.concurrent.ThreadLocalRandom;

public class RedPacket {
    public final int id;
    public final double amount;
    public final RedPacketShape shape;
    public final RedPacketSize size;
    public final double radius;
    public Vector2 pos;
    public Vector2 vel;
    public boolean collected = false;

    public RedPacket(int id, double minAmount, double maxAmount, int width, int height) {
        this.id = id;
        this.amount = ThreadLocalRandom.current().nextDouble(minAmount, maxAmount);
        
        // 随机选择形状和大小
        RedPacketShape[] shapes = RedPacketShape.values();
        RedPacketSize[] sizes = RedPacketSize.values();
        this.shape = shapes[ThreadLocalRandom.current().nextInt(shapes.length)];
        this.size = sizes[ThreadLocalRandom.current().nextInt(sizes.length)];
        this.radius = this.size.getRandomRadius();
        
        this.pos = new Vector2(
                ThreadLocalRandom.current().nextDouble(0, width),
                ThreadLocalRandom.current().nextDouble(0, height)
        );
        this.vel = new Vector2(
                ThreadLocalRandom.current().nextDouble(-0.8, 0.8),
                ThreadLocalRandom.current().nextDouble(-0.8, 0.8)
        );
    }

    /**
     * 获取红包的显示符号
     */
    public String getDisplaySymbol() {
        return shape.getSymbol();
    }

    /**
     * 获取红包的完整描述
     */
    public String getDescription() {
        return shape.getName() + size.getName() + "红包";
    }
}

