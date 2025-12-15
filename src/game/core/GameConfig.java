package game.core;

public class GameConfig {
    public final int width;
    public final int height;
    public final int redPacketCount;
    public final long durationMillis;
    public final long tickMillis;
    public final double playerRadius;

    private GameConfig(int width, int height, int redPacketCount, long durationMillis, long tickMillis, double playerRadius) {
        this.width = width;
        this.height = height;
        this.redPacketCount = redPacketCount;
        this.durationMillis = durationMillis;
        this.tickMillis = tickMillis;
        this.playerRadius = playerRadius;
    }

    public static GameConfig fromArgs(String[] args) {
        int width = 100;
        int height = 40;
        int redPacketCount = 30;
        long durationMillis = 5000; // 5s
        long tickMillis = 50;       // 20 FPS
        double playerRadius = 1.5;

        for (String arg : args) {
            if (arg.startsWith("--width=")) {
                width = Integer.parseInt(arg.substring("--width=".length()));
            } else if (arg.startsWith("--height=")) {
                height = Integer.parseInt(arg.substring("--height=".length()));
            } else if (arg.startsWith("--count=")) {
                redPacketCount = Integer.parseInt(arg.substring("--count=".length()));
            } else if (arg.startsWith("--duration=")) {
                durationMillis = Long.parseLong(arg.substring("--duration=".length()));
            } else if (arg.startsWith("--tick=")) {
                tickMillis = Long.parseLong(arg.substring("--tick=".length()));
            } else if (arg.startsWith("--radius=")) {
                playerRadius = Double.parseDouble(arg.substring("--radius=".length()));
            }
        }
        return new GameConfig(width, height, redPacketCount, durationMillis, tickMillis, playerRadius);
    }
}

