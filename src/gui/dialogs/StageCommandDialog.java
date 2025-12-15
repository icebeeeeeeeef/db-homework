package gui;

import javax.swing.*;
import java.awt.*;
import java.io.File;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Command dialog triggered by clicking Donald Duck.
 */
public class StageCommandDialog extends JDialog {
    private final JTabbedPane tabs = new JTabbedPane();

    // Rain controls
    private final JSpinner rainDuration = new JSpinner(new SpinnerNumberModel(15, 5, 120, 5));
    private final JSlider rainDensity = new JSlider(1, 5, 3);
    private final JCheckBox showStageRain = new JCheckBox("Show stage rain animation", true);
    private final JSpinner worldWidth = new JSpinner(new SpinnerNumberModel(80, 40, 400, 10));
    private final JSpinner worldHeight = new JSpinner(new SpinnerNumberModel(40, 20, 200, 5));
    private final JSpinner packetCount = new JSpinner(new SpinnerNumberModel(30, 10, 200, 5));
    private final JSpinner playerRadius = new JSpinner(new SpinnerNumberModel(1.5, 0.5, 5.0, 0.1));
    private final JSpinner fpsSpinner = new JSpinner(new SpinnerNumberModel(40, 10, 120, 5));
    private final JSpinner gameDuration = new JSpinner(new SpinnerNumberModel(30, 5, 180, 5)); // seconds

    // Code stats controls
    private final JTextField statsDirField = new JTextField(".");
    private final JTextField statsLangField = new JTextField("java,cpp,python");
    private final JCheckBox includeBlank = new JCheckBox("Show blank lines", true);
    private final JCheckBox includeComments = new JCheckBox("Show comment lines", true);
    private final JCheckBox includeFunctionStats = new JCheckBox("Include function length stats (C/Python)");
    private final JCheckBox pieChart = new JCheckBox("Use pie chart");

    // AI controls
    private final JTextArea aiPrompt = new JTextArea(5, 30);

    private StageCommand result;

    public StageCommandDialog(Frame owner) {
        super(owner, "Donald Command Center", true);
        setLayout(new BorderLayout(12, 12));
        tabs.setFont(new Font("SansSerif", Font.BOLD, 16));
        add(tabs, BorderLayout.CENTER);
        tabs.addTab("Red Packet Rain", createRainPanel());
        tabs.addTab("Code Stats", createStatsPanel());
        tabs.addTab("AI Chat", createAiPanel());

        JPanel buttons = new JPanel(new FlowLayout(FlowLayout.RIGHT, 16, 12));
        JButton ok = new JButton("Apply");
        JButton cancel = new JButton("Close");
        Font buttonFont = new Font("SansSerif", Font.BOLD, 16);
        ok.setFont(buttonFont);
        cancel.setFont(buttonFont);
        Dimension btnSize = new Dimension(140, 44);
        ok.setPreferredSize(btnSize);
        cancel.setPreferredSize(btnSize);
        buttons.add(ok);
        buttons.add(cancel);
        add(buttons, BorderLayout.SOUTH);

        ok.addActionListener(e -> onSubmit());
        cancel.addActionListener(e -> {
            result = null;
            setVisible(false);
        });

        setDefaultCloseOperation(DISPOSE_ON_CLOSE);
        pack();
        setSize(640, 520);
        setLocationRelativeTo(owner);
    }

    private JPanel createRainPanel() {
        JPanel panel = new JPanel(new GridBagLayout());
        GridBagConstraints gbc = new GridBagConstraints();
        gbc.insets = new Insets(8, 8, 8, 8);
        gbc.anchor = GridBagConstraints.WEST;
        gbc.fill = GridBagConstraints.HORIZONTAL;
        gbc.gridx = 0;
        gbc.gridy = 0;
        Font labelFont = new Font("SansSerif", Font.PLAIN, 16);

        JLabel durationLabel = new JLabel("Duration (seconds):");
        durationLabel.setFont(labelFont);
        panel.add(durationLabel, gbc);
        gbc.gridx = 1;
        rainDuration.setFont(labelFont);
        ((JSpinner.DefaultEditor) rainDuration.getEditor()).getTextField().setColumns(6);
        panel.add(rainDuration, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        JLabel densityLabel = new JLabel("Density (1-5):");
        densityLabel.setFont(labelFont);
        panel.add(densityLabel, gbc);
        gbc.gridx = 1;
        rainDensity.setMajorTickSpacing(1);
        rainDensity.setPaintLabels(true);
        rainDensity.setPaintTicks(true);
        rainDensity.setPreferredSize(new Dimension(260, 60));
        rainDensity.setFont(labelFont);
        panel.add(rainDensity, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        gbc.gridwidth = 2;
        showStageRain.setFont(labelFont);
        panel.add(showStageRain, gbc);

        gbc.gridy++;
        gbc.gridwidth = 1;
        panel.add(createSubLabel("Board Width:", labelFont), gbc);
        gbc.gridx = 1;
        configureSpinner(worldWidth);
        panel.add(worldWidth, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        panel.add(createSubLabel("Board Height:", labelFont), gbc);
        gbc.gridx = 1;
        configureSpinner(worldHeight);
        panel.add(worldHeight, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        panel.add(createSubLabel("Red Packet Count:", labelFont), gbc);
        gbc.gridx = 1;
        configureSpinner(packetCount);
        panel.add(packetCount, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        panel.add(createSubLabel("Player Radius:", labelFont), gbc);
        gbc.gridx = 1;
        configureSpinner(playerRadius);
        panel.add(playerRadius, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        panel.add(createSubLabel("Game Duration (seconds):", labelFont), gbc);
        gbc.gridx = 1;
        configureSpinner(gameDuration);
        panel.add(gameDuration, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        panel.add(createSubLabel("Game FPS:", labelFont), gbc);
        gbc.gridx = 1;
        configureSpinner(fpsSpinner);
        panel.add(fpsSpinner, gbc);
        return panel;
    }

    private JLabel createSubLabel(String text, Font font) {
        JLabel label = new JLabel(text);
        label.setFont(font);
        return label;
    }

    private void configureSpinner(JSpinner spinner) {
        spinner.setFont(new Font("SansSerif", Font.PLAIN, 16));
        if (spinner.getEditor() instanceof JSpinner.NumberEditor) {
            ((JSpinner.NumberEditor) spinner.getEditor()).getTextField().setColumns(6);
        }
    }

    private JPanel createStatsPanel() {
        JPanel panel = new JPanel(new GridBagLayout());
        GridBagConstraints gbc = new GridBagConstraints();
        gbc.insets = new Insets(8, 8, 8, 8);
        gbc.anchor = GridBagConstraints.WEST;
        gbc.fill = GridBagConstraints.HORIZONTAL;
        gbc.gridx = 0;
        gbc.gridy = 0;

        Font labelFont = new Font("SansSerif", Font.PLAIN, 16);
        JLabel dirLabel = new JLabel("Directory:");
        dirLabel.setFont(labelFont);
        panel.add(dirLabel, gbc);
        gbc.gridx = 1;
        statsDirField.setFont(labelFont);
        statsDirField.setColumns(18);
        panel.add(statsDirField, gbc);
        JButton browse = new JButton("Browse…");
        browse.setFont(labelFont);
        gbc.gridx = 2;
        panel.add(browse, gbc);
        browse.addActionListener(e -> onBrowse());

        gbc.gridx = 0;
        gbc.gridy++;
        JLabel langLabel = new JLabel("Languages (comma separated):");
        langLabel.setFont(labelFont);
        panel.add(langLabel, gbc);
        gbc.gridx = 1;
        gbc.gridwidth = 2;
        statsLangField.setFont(labelFont);
        statsLangField.setColumns(18);
        panel.add(statsLangField, gbc);

        gbc.gridx = 0;
        gbc.gridy++;
        gbc.gridwidth = 3;
        includeBlank.setFont(labelFont);
        panel.add(includeBlank, gbc);
        gbc.gridy++;
        includeComments.setFont(labelFont);
        panel.add(includeComments, gbc);
        gbc.gridy++;
        includeFunctionStats.setFont(labelFont);
        panel.add(includeFunctionStats, gbc);
        gbc.gridy++;
        pieChart.setFont(labelFont);
        panel.add(pieChart, gbc);

        return panel;
    }

    private JPanel createAiPanel() {
        JPanel panel = new JPanel(new BorderLayout(8, 8));
        aiPrompt.setLineWrap(true);
        aiPrompt.setWrapStyleWord(true);
        aiPrompt.setFont(new Font("SansSerif", Font.PLAIN, 16));
        JLabel promptLabel = new JLabel("Prompt:");
        promptLabel.setFont(new Font("SansSerif", Font.PLAIN, 16));
        panel.add(promptLabel, BorderLayout.NORTH);
        panel.add(new JScrollPane(aiPrompt), BorderLayout.CENTER);
        return panel;
    }

    private void onBrowse() {
        JFileChooser chooser = new JFileChooser(statsDirField.getText());
        chooser.setFileSelectionMode(JFileChooser.DIRECTORIES_ONLY);
        if (chooser.showOpenDialog(this) == JFileChooser.APPROVE_OPTION) {
            File selected = chooser.getSelectedFile();
            if (selected != null) {
                statsDirField.setText(selected.getAbsolutePath());
            }
        }
    }

    private void onSubmit() {
        int tab = tabs.getSelectedIndex();
        if (tab == 0) {
            int duration = ((Number) rainDuration.getValue()).intValue();
            int density = rainDensity.getValue();
            int width = ((Number) worldWidth.getValue()).intValue();
            int height = ((Number) worldHeight.getValue()).intValue();
            int count = ((Number) packetCount.getValue()).intValue();
            double radius = ((Number) playerRadius.getValue()).doubleValue();
            long gameDurationMillis = ((Number) gameDuration.getValue()).longValue() * 1000L;
            int fps = ((Number) fpsSpinner.getValue()).intValue();
            result = StageCommand.rain(new StageCommand.RainOptions(
                    duration,
                    density,
                    showStageRain.isSelected(),
                    width,
                    height,
                    count,
                    gameDurationMillis,
                    fps,
                    radius
            ));
        } else if (tab == 1) {
            String dir = statsDirField.getText().trim();
            if (dir.isEmpty()) {
                JOptionPane.showMessageDialog(this, "Please choose a directory to analyze.", "Info", JOptionPane.INFORMATION_MESSAGE);
                return;
            }
            List<String> langs = Arrays.stream(statsLangField.getText().split(","))
                    .map(String::trim)
                    .filter(s -> !s.isEmpty())
                    .collect(Collectors.toList());
            result = StageCommand.codeStats(
                    new StageCommand.CodeStatsOptions(dir, langs, includeBlank.isSelected(), includeComments.isSelected(),
                            includeFunctionStats.isSelected(), pieChart.isSelected()));
        } else {
            String prompt = aiPrompt.getText().trim();
            if (prompt.isEmpty()) {
                JOptionPane.showMessageDialog(this, "Please enter a prompt.", "Info", JOptionPane.INFORMATION_MESSAGE);
                return;
            }
            result = StageCommand.ai(new StageCommand.AiOptions(prompt));
        }
        setVisible(false);
    }

    public StageCommand showDialog() {
        setVisible(true);
        return result;
    }
}
