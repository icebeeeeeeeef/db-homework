package gui;

import model.duck.DuckCharacter;
import model.duck.behavior.ActionBehavior;
import model.duck.behavior.BehaviorLibrary;
import model.duck.behavior.DuckBehaviorProfile;
import model.duck.behavior.SoundBehavior;

import javax.swing.*;
import java.awt.*;
import java.util.List;

public class BehaviorDialog extends JDialog {
    public static class Selection {
        public final DuckCharacter character;
        public final DuckBehaviorProfile profile;

        public Selection(DuckCharacter character, DuckBehaviorProfile profile) {
            this.character = character;
            this.profile = profile;
        }
    }

    private final JComboBox<DuckCharacter> duckCombo = new JComboBox<>(DuckCharacter.values());
    private final JComboBox<ActionBehavior> actionCombo;
    private final JComboBox<SoundBehavior> soundCombo;
    private Selection selection;

    public BehaviorDialog(Frame owner) {
        super(owner, "Behavior Playground", true);
        setLayout(new BorderLayout(12, 12));

        List<ActionBehavior> actions = BehaviorLibrary.availableActions();
        List<SoundBehavior> sounds = BehaviorLibrary.availableSounds();
        actionCombo = new JComboBox<>(actions.toArray(new ActionBehavior[0]));
        soundCombo = new JComboBox<>(sounds.toArray(new SoundBehavior[0]));

        Font labelFont = new Font("SansSerif", Font.PLAIN, 16);
        duckCombo.setFont(labelFont);
        actionCombo.setFont(labelFont);
        soundCombo.setFont(labelFont);

        JPanel form = new JPanel(new GridLayout(0, 2, 10, 12));
        form.setBorder(BorderFactory.createEmptyBorder(10, 18, 10, 18));

        form.add(createLabel("Duck:", labelFont));
        form.add(duckCombo);
        form.add(createLabel("Action:", labelFont));
        form.add(actionCombo);
        form.add(createLabel("Sound:", labelFont));
        form.add(soundCombo);

        add(form, BorderLayout.CENTER);

        JPanel buttons = new JPanel(new FlowLayout(FlowLayout.RIGHT, 12, 8));
        JButton apply = new JButton("Apply");
        JButton cancel = new JButton("Cancel");
        apply.setFont(labelFont);
        cancel.setFont(labelFont);
        buttons.add(apply);
        buttons.add(cancel);
        add(buttons, BorderLayout.SOUTH);

        apply.addActionListener(e -> {
            ActionBehavior action = (ActionBehavior) actionCombo.getSelectedItem();
            SoundBehavior sound = (SoundBehavior) soundCombo.getSelectedItem();
            DuckCharacter character = (DuckCharacter) duckCombo.getSelectedItem();
            if (action != null && sound != null && character != null) {
                selection = new Selection(character, new DuckBehaviorProfile(action, sound));
                setVisible(false);
            }
        });
        cancel.addActionListener(e -> {
            selection = null;
            setVisible(false);
        });

        setSize(420, 240);
        setLocationRelativeTo(owner);
    }

    private JLabel createLabel(String text, Font font) {
        JLabel label = new JLabel(text);
        label.setFont(font);
        return label;
    }

    public Selection showDialog() {
        setVisible(true);
        return selection;
    }
}
