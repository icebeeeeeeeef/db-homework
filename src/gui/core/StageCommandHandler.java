package gui;

import javax.swing.*;
import javax.swing.table.DefaultTableModel;
import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.Font;
import java.awt.GridLayout;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import infra.ToolPaths;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;

/**
 * Handles stage interaction commands such as rain animations and code statistics.
 */
final class StageCommandHandler {
    private final JFrame owner;
    private final StagePanel stagePanel;
    private final Consumer<String> chatAppender;
    private final Consumer<StageCommand.AiOptions> aiHandler;
    private StageCommandDialog stageCommandDialog;

    StageCommandHandler(JFrame owner,
                        StagePanel stagePanel,
                        Consumer<String> chatAppender,
                        Consumer<StageCommand.AiOptions> aiHandler) {
        this.owner = owner;
        this.stagePanel = stagePanel;
        this.chatAppender = chatAppender;
        this.aiHandler = aiHandler;
    }

    void openStageCommandDialog() {
        if (stageCommandDialog == null) {
            stageCommandDialog = new StageCommandDialog(owner);
        }
        StageCommand command = stageCommandDialog.showDialog();
        if (command != null) {
            handleStageCommand(command);
        }
    }

    void runCodeStats(StageCommand.CodeStatsOptions options) {
        StageCommand.CodeStatsOptions fallback =
                options != null ? options :
                new StageCommand.CodeStatsOptions(".", List.of(), true, true, true, false);
        runCodeStatsCommand(fallback);
    }

    private void handleStageCommand(StageCommand command) {
        switch (command.type) {
            case RED_PACKET_RAIN:
                runRainCommand(command.rainOptions);
                break;
            case CODE_STATS:
                runCodeStats(command.codeStatsOptions);
                break;
            case AI_CHAT:
                aiHandler.accept(command.aiOptions);
                break;
        }
    }

    private void runRainCommand(StageCommand.RainOptions options) {
        if (options == null) return;
        if (options.showStageRain) {
            chatAppender.accept("🦆 Duck: Launching a " + options.durationSeconds + " second red packet rain!");
            stagePanel.setSpeech("Rain incoming!");
            int durationMs = Math.max(5, options.durationSeconds) * 1000;
            stagePanel.triggerRain(durationMs, Math.max(1, options.density));
            new javax.swing.Timer(durationMs, e -> stagePanel.setSpeech("Ready for the next mission?")) {{
                setRepeats(false);
            }}.start();
        }
        chatAppender.accept("🦆 Duck: Opening a dedicated red packet game window!");
        RainGameWindow window = new RainGameWindow(options, () -> stagePanel.setSpeech("Game session finished!"));
        window.launch();
    }

    private void runCodeStatsCommand(StageCommand.CodeStatsOptions opts) {
        chatAppender.accept("🦆 Duck: Crunching code stats for " + opts.directory + "...");
        stagePanel.setSpeech("Analyzing...");

        SwingWorker<CodeStatsResult, Void> worker = new SwingWorker<>() {
            @Override
            protected CodeStatsResult doInBackground() throws Exception {
                List<String> cmd = new ArrayList<>();
                cmd.add(ToolPaths.codeStatsExecutable());
                cmd.add("--tsv");
                cmd.add("--dir=" + opts.directory);
                if (opts.includeFunctionStats) {
                    cmd.add("--functions");
                }
                if (opts.languages != null && !opts.languages.isEmpty()) {
                    cmd.add("--languages=" + String.join(",", opts.languages));
                }

                ProcessBuilder pb = new ProcessBuilder(cmd);
                pb.directory(new java.io.File(System.getProperty("user.dir")));
                Process process = pb.start();
                List<String> lines = new ArrayList<>();
                try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()))) {
                    String ln;
                    while ((ln = reader.readLine()) != null) {
                        lines.add(ln);
                    }
                }
                int exitCode = process.waitFor();
                if (exitCode != 0) {
                    throw new IllegalStateException("code_stats exited with " + exitCode);
                }
                return CodeStatsResult.parse(lines);
            }

            @Override
            protected void done() {
                try {
                    CodeStatsResult result = get();
                    if (result == null || result.getLanguages().isEmpty()) {
                        chatAppender.accept("🦆 Duck: No results to display.");
                    } else {
                        showCodeStatsDialog(result, opts);
                        chatAppender.accept("🦆 Duck: Stats ready! Charts displayed.");
                    }
                } catch (Exception e) {
                    chatAppender.accept("🦆 Duck: Stat run failed - " + e.getMessage());
                } finally {
                    stagePanel.setSpeech("Analysis complete!");
                }
            }
        };
        worker.execute();
    }

    private void showCodeStatsDialog(CodeStatsResult result, StageCommand.CodeStatsOptions options) {
        JDialog dialog = new JDialog(owner, "Code Statistics", true);
        dialog.setLayout(new BorderLayout(12, 12));

        JTabbedPane tabs = new JTabbedPane();
        tabs.addTab("Summary", buildSummaryPanel(result, options));

        JTabbedPane charts = new JTabbedPane();
        List<CodeStatsChartPanel.Item> chartItems = toChartItems(result);
        charts.addTab("Bar Chart", new CodeStatsChartPanel(chartItems, CodeStatsChartPanel.ChartMode.BAR));
        charts.addTab("Pie Chart", new CodeStatsChartPanel(chartItems, CodeStatsChartPanel.ChartMode.PIE));
        if (options != null && options.pieChart) {
            charts.setSelectedIndex(1);
        }
        tabs.addTab("Charts", charts);

        if (result.getFunctions() != null && !result.getFunctions().isEmpty()) {
            tabs.addTab("Function Stats", new FunctionStatsPanel(result.getFunctions()));
        }

        dialog.add(tabs, BorderLayout.CENTER);
        dialog.setSize(960, 600);
        dialog.setLocationRelativeTo(owner);
        dialog.setVisible(true);
    }

    private JPanel buildSummaryPanel(CodeStatsResult result, StageCommand.CodeStatsOptions options) {
        JPanel panel = new JPanel(new BorderLayout(8, 8));
        List<String> columns = new ArrayList<>();
        columns.add("Language");
        columns.add("Files");
        columns.add("Total Lines");
        columns.add("Code Lines");
        boolean showComments = options == null || options.includeComments;
        boolean showBlanks = options == null || options.includeBlank;
        if (showComments) columns.add("Comment Lines");
        if (showBlanks) columns.add("Blank Lines");

        DefaultTableModel model = new DefaultTableModel(columns.toArray(), 0) {
            @Override
            public boolean isCellEditable(int row, int column) {
                return false;
            }
        };
        JTable table = new JTable(model);
        table.setRowHeight(26);
        table.setFont(font(14f, Font.PLAIN));
        table.getTableHeader().setFont(font(14f, Font.BOLD));

        for (CodeStatsResult.LanguageStat stat : result.getLanguages()) {
            List<Object> row = new ArrayList<>();
            row.add(stat.name);
            row.add(stat.files);
            row.add(stat.total);
            row.add(stat.code);
            if (showComments) row.add(stat.comments);
            if (showBlanks) row.add(stat.blanks);
            model.addRow(row.toArray());
        }

        panel.add(new JScrollPane(table), BorderLayout.CENTER);

        if (result.getTotal() != null) {
            JPanel summary = new JPanel(new GridLayout(0, 1));
            summary.setBorder(BorderFactory.createTitledBorder("Totals"));
            summary.add(new JLabel("Files: " + result.getTotal().files));
            summary.add(new JLabel("Lines: " + result.getTotal().total));
            summary.add(new JLabel("Code: " + result.getTotal().code));
            if (showComments) summary.add(new JLabel("Comments: " + result.getTotal().comments));
            if (showBlanks) summary.add(new JLabel("Blanks: " + result.getTotal().blanks));

            summary.setFont(font(14f, Font.PLAIN));
            panel.add(summary, BorderLayout.EAST);
        }
        return panel;
    }

    private List<CodeStatsChartPanel.Item> toChartItems(CodeStatsResult result) {
        List<CodeStatsChartPanel.Item> items = new ArrayList<>();
        Color[] palette = {
                new Color(255, 99, 71), new Color(135, 206, 235),
                new Color(60, 179, 113), new Color(255, 165, 0),
                new Color(147, 112, 219), new Color(70, 130, 180),
                new Color(220, 20, 60), new Color(46, 139, 87)
        };
        int idx = 0;
        for (CodeStatsResult.LanguageStat stat : result.getLanguages()) {
            items.add(new CodeStatsChartPanel.Item(stat.name, stat.total, palette[idx % palette.length]));
            idx++;
        }
        return items;
    }

    private Font font(float size, int style) {
        return DuckUiTheme.font(size, style);
    }
}
