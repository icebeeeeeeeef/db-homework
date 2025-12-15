package model.redpacket;

/**
 * 红包大小枚举
 */
public enum RedPacketSize {
    SMALL("小", 0.5, 0.8),
    MEDIUM("中", 0.8, 1.2),
    LARGE("大", 1.2, 2.0),
    HUGE("超大", 2.0, 3.0);

    private final String name;
    private final double minRadius;
    private final double maxRadius;

    RedPacketSize(String name, double minRadius, double maxRadius) {
        this.name = name;
        this.minRadius = minRadius;
        this.maxRadius = maxRadius;
    }

    public String getName() {
        return name;
    }

    public double getMinRadius() {
        return minRadius;
    }

    public double getMaxRadius() {
        return maxRadius;
    }

    /**
     * 获取随机半径
     */
    public double getRandomRadius() {
        return minRadius + Math.random() * (maxRadius - minRadius);
    }
}
