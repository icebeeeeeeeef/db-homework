package infra;

import java.nio.file.Path;

/**
 * Centralizes file-system locations for helper tools.
 */
public final class ToolPaths {
    public static final Path CODE_STATS = Path.of("tools", "code-stats", "code_stats");

    private ToolPaths() {}

    public static String codeStatsExecutable() {
        return CODE_STATS.toString();
    }
}
