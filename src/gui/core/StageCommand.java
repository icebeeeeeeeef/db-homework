package gui;

import java.util.List;

public class StageCommand {
    public enum Type {
        RED_PACKET_RAIN,
        CODE_STATS,
        AI_CHAT
    }

    public static class RainOptions {
        public final int durationSeconds;
        public final int density;
        public final boolean showStageRain;
        public final int worldWidth;
        public final int worldHeight;
        public final int redPacketCount;
        public final long gameDurationMillis;
        public final int fps;
        public final double playerRadius;

        public RainOptions(int durationSeconds,
                           int density,
                           boolean showStageRain,
                           int worldWidth,
                           int worldHeight,
                           int redPacketCount,
                           long gameDurationMillis,
                           int fps,
                           double playerRadius) {
            this.durationSeconds = durationSeconds;
            this.density = density;
            this.showStageRain = showStageRain;
            this.worldWidth = worldWidth;
            this.worldHeight = worldHeight;
            this.redPacketCount = redPacketCount;
            this.gameDurationMillis = gameDurationMillis;
            this.fps = fps;
            this.playerRadius = playerRadius;
        }
    }

    public static class CodeStatsOptions {
        public final String directory;
        public final List<String> languages;
        public final boolean includeBlank;
        public final boolean includeComments;
        public final boolean includeFunctionStats;
        public final boolean pieChart;

        public CodeStatsOptions(String directory,
                                List<String> languages,
                                boolean includeBlank,
                                boolean includeComments,
                                boolean includeFunctionStats,
                                boolean pieChart) {
            this.directory = directory;
            this.languages = languages;
            this.includeBlank = includeBlank;
            this.includeComments = includeComments;
            this.includeFunctionStats = includeFunctionStats;
            this.pieChart = pieChart;
        }
    }

    public static class AiOptions {
        public final String prompt;

        public AiOptions(String prompt) {
            this.prompt = prompt;
        }
    }

    public final Type type;
    public final RainOptions rainOptions;
    public final CodeStatsOptions codeStatsOptions;
    public final AiOptions aiOptions;

    private StageCommand(Type type, RainOptions rainOptions, CodeStatsOptions codeStatsOptions, AiOptions aiOptions) {
        this.type = type;
        this.rainOptions = rainOptions;
        this.codeStatsOptions = codeStatsOptions;
        this.aiOptions = aiOptions;
    }

    public static StageCommand rain(RainOptions options) {
        return new StageCommand(Type.RED_PACKET_RAIN, options, null, null);
    }

    public static StageCommand codeStats(CodeStatsOptions options) {
        return new StageCommand(Type.CODE_STATS, null, options, null);
    }

    public static StageCommand ai(AiOptions options) {
        return new StageCommand(Type.AI_CHAT, null, null, options);
    }
}
