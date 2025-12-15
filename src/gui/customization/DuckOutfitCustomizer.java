package gui;

import model.duck.DuckCharacter;
import model.duck.customization.DuckOutfit;

import javax.swing.*;
import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Component;
import java.awt.Dimension;
import java.awt.FlowLayout;
import java.awt.Font;
import java.awt.GridLayout;
import java.util.EnumMap;
import java.util.function.Consumer;
import java.util.function.Function;

/**
 * Provides the duck outfit customization dialog logic.
 */
final class DuckOutfitCustomizer {
    private DuckOutfitCustomizer() {}

    static void open(JFrame owner,
                     StagePanel stagePanel,
                     Consumer<String> chatAppender,
                     Function<DuckCharacter, DuckOutfit> defaultFactory) {
        JDialog dlg = new JDialog(owner, "Duck Outfit Studio", true);
        dlg.setLayout(new BorderLayout(12, 12));

        EnumMap<DuckCharacter, DuckOutfit> working = new EnumMap<>(DuckCharacter.class);
        for (DuckCharacter character : DuckCharacter.values()) {
            working.put(character, stagePanel.getOutfit(character).clone());
        }

        JComboBox<DuckCharacter> targetSelector = new JComboBox<>(DuckCharacter.values());
        targetSelector.setFont(font(16f, Font.BOLD));

        JCheckBox hat = new JCheckBox("Hat");
        JCheckBox scarf = new JCheckBox("Scarf");
        JCheckBox eyes = new JCheckBox("Glasses");
        JCheckBox tie = new JCheckBox("Necktie");
        JCheckBox cane = new JCheckBox("Cane");
        Font optionFont = font(16f, Font.PLAIN);
        hat.setFont(optionFont);
        scarf.setFont(optionFont);
        eyes.setFont(optionFont);
        tie.setFont(optionFont);
        cane.setFont(optionFont);

        JButton hatColorBtn = new JButton("Hat Color");
        JButton scarfColorBtn = new JButton("Scarf Color");
        JButton eyeColorBtn = new JButton("Glasses Color");
        JButton tieColorBtn = new JButton("Tie Color");
        JButton caneColorBtn = new JButton("Cane Color");
        hatColorBtn.setFont(optionFont);
        scarfColorBtn.setFont(optionFont);
        eyeColorBtn.setFont(optionFont);
        tieColorBtn.setFont(optionFont);
        caneColorBtn.setFont(optionFont);

        JPanel grid = new JPanel(new GridLayout(0, 2, 8, 8));
        grid.setBorder(BorderFactory.createEmptyBorder(10, 10, 10, 10));
        JLabel characterLabel = new JLabel("Character:");
        characterLabel.setFont(font(16f, Font.BOLD));
        grid.add(characterLabel);
        grid.add(targetSelector);
        grid.add(hat);
        grid.add(hatColorBtn);
        grid.add(scarf);
        grid.add(scarfColorBtn);
        grid.add(eyes);
        grid.add(eyeColorBtn);
        grid.add(tie);
        grid.add(tieColorBtn);
        grid.add(cane);
        grid.add(caneColorBtn);

        dlg.add(grid, BorderLayout.CENTER);

        JPanel buttons = new JPanel(new FlowLayout(FlowLayout.RIGHT, 12, 8));
        JButton apply = new JButton("Apply to Selected");
        JButton applyAll = new JButton("Apply to All Ducklings");
        JButton resetSelected = new JButton("Reset to Default");
        JButton close = new JButton("Close");
        apply.setFont(font(18f, Font.BOLD));
        applyAll.setFont(optionFont);
        resetSelected.setFont(optionFont);
        close.setFont(optionFont);
        apply.setPreferredSize(new Dimension(200, 44));
        applyAll.setPreferredSize(new Dimension(240, 44));
        resetSelected.setPreferredSize(new Dimension(200, 44));
        close.setPreferredSize(new Dimension(120, 44));
        buttons.add(apply);
        buttons.add(applyAll);
        buttons.add(resetSelected);
        buttons.add(close);
        dlg.add(buttons, BorderLayout.SOUTH);
        dlg.getRootPane().setDefaultButton(apply);

        Runnable refresh = () -> {
            DuckCharacter target = (DuckCharacter) targetSelector.getSelectedItem();
            DuckOutfit outfit = working.get(target);
            hat.setSelected(outfit.hasHat());
            scarf.setSelected(outfit.hasScarf());
            eyes.setSelected(outfit.hasEyes());
            tie.setSelected(outfit.hasTie());
            cane.setSelected(outfit.hasCane());
            setButtonColor(hatColorBtn, outfit.getHatColor());
            setButtonColor(scarfColorBtn, outfit.getScarfColor());
            setButtonColor(eyeColorBtn, outfit.getEyeFrameColor());
            setButtonColor(tieColorBtn, outfit.getTieColor());
            setButtonColor(caneColorBtn, outfit.getCaneColor());
        };

        targetSelector.addActionListener(e -> refresh.run());

        hatColorBtn.addActionListener(e -> pickColor(dlg, hatColorBtn, working, targetSelector, DuckColorTarget.HAT));
        scarfColorBtn.addActionListener(e -> pickColor(dlg, scarfColorBtn, working, targetSelector, DuckColorTarget.SCARF));
        eyeColorBtn.addActionListener(e -> pickColor(dlg, eyeColorBtn, working, targetSelector, DuckColorTarget.EYES));
        tieColorBtn.addActionListener(e -> pickColor(dlg, tieColorBtn, working, targetSelector, DuckColorTarget.TIE));
        caneColorBtn.addActionListener(e -> pickColor(dlg, caneColorBtn, working, targetSelector, DuckColorTarget.CANE));

        hat.addActionListener(e -> working.get(targetSelector.getItemAt(targetSelector.getSelectedIndex())).setHat(hat.isSelected()));
        scarf.addActionListener(e -> working.get(targetSelector.getItemAt(targetSelector.getSelectedIndex())).setScarf(scarf.isSelected()));
        eyes.addActionListener(e -> working.get(targetSelector.getItemAt(targetSelector.getSelectedIndex())).setEyes(eyes.isSelected()));
        tie.addActionListener(e -> working.get(targetSelector.getItemAt(targetSelector.getSelectedIndex())).setTie(tie.isSelected()));
        cane.addActionListener(e -> working.get(targetSelector.getItemAt(targetSelector.getSelectedIndex())).setCane(cane.isSelected()));

        apply.addActionListener(e -> {
            DuckCharacter target = (DuckCharacter) targetSelector.getSelectedItem();
            stagePanel.setOutfit(target, working.get(target));
            chatAppender.accept("🦆 Duck: Outfit updated for " + target.getDisplayName());
            refresh.run();
        });

        applyAll.addActionListener(e -> {
            DuckOutfit template = working.get(DuckCharacter.DONALD).clone();
            for (DuckCharacter c : DuckCharacter.values()) {
                if (c == DuckCharacter.DONALD) continue;
                stagePanel.setOutfit(c, template);
                working.put(c, template.clone());
            }
            chatAppender.accept("🦆 Duck: Ducklings updated with the same style!");
            refresh.run();
        });

        resetSelected.addActionListener(e -> {
            DuckCharacter target = (DuckCharacter) targetSelector.getSelectedItem();
            DuckOutfit preset = defaultFactory.apply(target);
            working.put(target, preset.clone());
            refresh.run();
            chatAppender.accept("🦆 Duck: Default style loaded for " + target.getDisplayName() + ". Press Apply to Selected to confirm.");
        });

        close.addActionListener(e -> dlg.dispose());

        refresh.run();
        dlg.pack();
        dlg.setSize(420, 360);
        dlg.setLocationRelativeTo(owner);
        dlg.setVisible(true);
    }

    private static void pickColor(Component owner,
                                  JButton button,
                                  EnumMap<DuckCharacter, DuckOutfit> working,
                                  JComboBox<DuckCharacter> targetSelector,
                                  DuckColorTarget target) {
        DuckCharacter selected = (DuckCharacter) targetSelector.getSelectedItem();
        DuckOutfit outfit = working.get(selected);
        Color current = getColorForTarget(outfit, target);
        Color color = JColorChooser.showDialog(owner, "Pick a color", current);
        if (color != null) {
            setColorForTarget(outfit, target, color);
            setButtonColor(button, color);
        }
    }

    private static Color getColorForTarget(DuckOutfit outfit, DuckColorTarget target) {
        switch (target) {
            case HAT: return outfit.getHatColor();
            case SCARF: return outfit.getScarfColor();
            case EYES: return outfit.getEyeFrameColor();
            case TIE: return outfit.getTieColor();
            case CANE: return outfit.getCaneColor();
            default: return Color.WHITE;
        }
    }

    private static void setColorForTarget(DuckOutfit outfit, DuckColorTarget target, Color color) {
        switch (target) {
            case HAT: outfit.setHatColor(color); break;
            case SCARF: outfit.setScarfColor(color); break;
            case EYES: outfit.setEyeFrameColor(color); break;
            case TIE: outfit.setTieColor(color); break;
            case CANE: outfit.setCaneColor(color); break;
        }
    }

    private static void setButtonColor(JButton button, Color color) {
        button.setBackground(color);
        button.setForeground(color.getRed() + color.getGreen() + color.getBlue() > 382 ? Color.BLACK : Color.WHITE);
    }

    private static Font font(float size, int style) {
        return DuckUiTheme.font(size, style);
    }

    private enum DuckColorTarget { HAT, SCARF, EYES, TIE, CANE }
}
