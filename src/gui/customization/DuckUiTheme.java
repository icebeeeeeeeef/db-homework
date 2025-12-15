package gui;

import java.awt.Font;
import java.awt.GraphicsEnvironment;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

/**
 * Provides consistent font handling for the duck assistant experience.
 */
final class DuckUiTheme {
    private static final String[] FONT_CANDIDATES = {
            "Noto Sans CJK SC", "Noto Sans SC", "WenQuanYi Micro Hei",
            "Source Han Sans CN", "PingFang SC", "Microsoft YaHei",
            "SimHei", "Sarasa Gothic SC", "STHeiti", "Heiti SC",
            "LXGW WenKai", "Arial Unicode MS", "DejaVu Sans", "Dialog", Font.SANS_SERIF
    };
    private static final Set<String> AVAILABLE_FONTS = new HashSet<>();

    static {
        try {
            String[] names = GraphicsEnvironment.getLocalGraphicsEnvironment().getAvailableFontFamilyNames();
            AVAILABLE_FONTS.addAll(Arrays.asList(names));
        } catch (Exception ignored) {
        }
    }

    private DuckUiTheme() {}

    static Font font(float size, int style) {
        for (String name : FONT_CANDIDATES) {
            if (AVAILABLE_FONTS.contains(name)) {
                return new Font(name, style, Math.round(size));
            }
        }
        return new Font(Font.SANS_SERIF, style, Math.round(size));
    }
}
