package gui;

import javax.swing.*;
import java.awt.*;

/**
 * Dedicated window for running a configured red packet game.
 */
public class RainGameWindow extends JFrame {
    private final StageCommand.RainOptions options;
    private final Runnable onClose;

    public RainGameWindow(StageCommand.RainOptions options, Runnable onClose) {
        super("Red Packet Rain");
        this.options = options;
        this.onClose = onClose;
        setDefaultCloseOperation(JFrame.DISPOSE_ON_CLOSE);
        setLocationRelativeTo(null);
        setResizable(false);
    }

    public void launch() {
        SwingUtilities.invokeLater(this::startGame);
    }

    private void startGame() {
        getContentPane().removeAll();
        GuiGame game = new GuiGame(
                options.worldWidth,
                options.worldHeight,
                options.redPacketCount,
                options.playerRadius,
                options.gameDurationMillis,
                options.fps
        );
        game.setGameOverListener((count, amount) -> {
            boolean again = ResultDialog.show(this, count, amount, game.getStatistics());
            if (again) {
                startGame();
            } else {
                dispose();
                if (onClose != null) {
                    onClose.run();
                }
            }
        });
        setContentPane(game);
        pack();
        setLocationRelativeTo(null);
        setVisible(true);
        game.start();
    }
}
