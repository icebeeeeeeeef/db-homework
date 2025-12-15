package app;

import assistant.cli.DuckAssistant;
import game.core.GameConfig;
import game.core.GameEngine;
import gui.GuiGame;
import gui.DuckAssistantGUI;
import attendance.ui.AttendanceApp;

import java.util.*;

public class Main {
    public static void main(String[] args) {
        // 检查是否启动小鸭子助手
        boolean duckMode = Arrays.stream(args).anyMatch(a -> a.equals("--duck") || a.equals("-d"));
        boolean duckGUIMode = Arrays.stream(args).anyMatch(a -> a.equals("--duck-gui") || a.equals("-dg"));
        boolean attendanceMode = Arrays.stream(args).anyMatch(a -> a.equals("--attendance") || a.equals("-a"));
        
        if (attendanceMode) {
            AttendanceApp.launch();
        } else if (duckGUIMode) {
            // 启动图形化小鸭子助手
            DuckAssistantGUI.launch();
        } else if (duckMode) {
            // 启动命令行小鸭子助手
            DuckAssistant duck = new DuckAssistant();
            duck.start();
        } else {
            // 原有的游戏逻辑
            GameConfig config = GameConfig.fromArgs(args);
            boolean gui = Arrays.stream(args).anyMatch(a -> a.equals("--gui"));
            if (gui) {
                GuiGame.launchInteractive();
            } else {
                GameEngine engine = new GameEngine(config);
                engine.start();
                engine.printSummary();
            }
        }
    }
}
