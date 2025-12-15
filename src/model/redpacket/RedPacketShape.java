package model.redpacket;

/**
 * 红包形状枚举
 */
public enum RedPacketShape {
    CIRCLE("圆形", "●"),
    SQUARE("方形", "■"),
    TRIANGLE("三角形", "▲"),
    DIAMOND("菱形", "♦"),
    HEART("心形", "♥"),
    STAR("星形", "★");

    private final String name;
    private final String symbol;

    RedPacketShape(String name, String symbol) {
        this.name = name;
        this.symbol = symbol;
    }

    public String getName() {
        return name;
    }

    public String getSymbol() {
        return symbol;
    }
}
