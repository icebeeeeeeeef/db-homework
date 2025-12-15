package gui;

import java.util.ArrayList;
import java.util.List;

public class CodeStatsResult {
    public static class LanguageStat {
        public final String name;
        public final int files;
        public final int total;
        public final int code;
        public final int comments;
        public final int blanks;

        public LanguageStat(String name, int files, int total, int code, int comments, int blanks) {
            this.name = name;
            this.files = files;
            this.total = total;
            this.code = code;
            this.comments = comments;
            this.blanks = blanks;
        }
    }

    public static class FunctionStat {
        public final String language;
        public final int count;
        public final double average;
        public final int min;
        public final int max;
        public final double median;

        public FunctionStat(String language, int count, double average, int min, int max, double median) {
            this.language = language;
            this.count = count;
            this.average = average;
            this.min = min;
            this.max = max;
            this.median = median;
        }
    }

    private final List<LanguageStat> languages = new ArrayList<>();
    private final List<FunctionStat> functions = new ArrayList<>();
    private LanguageStat total;

    public List<LanguageStat> getLanguages() {
        return languages;
    }

    public List<FunctionStat> getFunctions() {
        return functions;
    }

    public LanguageStat getTotal() {
        return total;
    }

    public static CodeStatsResult parse(List<String> lines) {
        CodeStatsResult result = new CodeStatsResult();
        for (String line : lines) {
            if (line == null || line.isBlank()) continue;
            String[] parts = line.split("\t");
            if (parts.length == 0) continue;
            if ("TOTAL".equalsIgnoreCase(parts[0]) && parts.length >= 6) {
                result.total = new LanguageStat(parts[0],
                        safeInt(parts[1]), safeInt(parts[2]), safeInt(parts[3]),
                        safeInt(parts[4]), safeInt(parts[5]));
                continue;
            } else if ("FUNC".equalsIgnoreCase(parts[0]) && parts.length >= 7) {
                result.functions.add(new FunctionStat(
                        parts[1],
                        safeInt(parts[2]),
                        safeDouble(parts[3]),
                        safeInt(parts[4]),
                        safeInt(parts[5]),
                        safeDouble(parts[6])
                ));
                continue;
            } else if (parts.length >= 6) {
                result.languages.add(new LanguageStat(
                        parts[0],
                        safeInt(parts[1]),
                        safeInt(parts[2]),
                        safeInt(parts[3]),
                        safeInt(parts[4]),
                        safeInt(parts[5])
                ));
            }
        }
        return result;
    }

    private static int safeInt(String text) {
        try {
            return Integer.parseInt(text.trim());
        } catch (Exception e) {
            return 0;
        }
    }

    private static double safeDouble(String text) {
        try {
            return Double.parseDouble(text.trim());
        } catch (Exception e) {
            return 0.0;
        }
    }
}
