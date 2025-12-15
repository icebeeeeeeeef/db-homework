package assistant.cli;

import ai.AiClient;
import ai.AiClient.Message;
import app.Main;
import infra.ToolPaths;

import java.io.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

/**
 * 小鸭子助手 - 可爱的交互界面
 */
public class DuckAssistant {
    private static final String DONALD_PROMPT = "You are Donald Duck in a game world. "
            + "Speak in Donald Duck's lively, humorous style. Keep responses concise "
            + "and reference the player's adventures. If the player speaks Chinese, reply in Chinese "
            + "but maintain your Donald Duck personality.";

    private final Scanner scanner;
    private boolean isRunning;
    private AiClient aiClient;
    private boolean aiEnabled;
    private final List<Message> conversation = new ArrayList<>();
    
    public DuckAssistant() {
        this.scanner = new Scanner(System.in);
        this.isRunning = true;
        initAiClient();
    }
    
    public void start() {
        showWelcome();
        
        while (isRunning) {
            showMainMenu();
            handleUserInput();
        }
        
        sayGoodbye();
    }
    
    private void showWelcome() {
        clearScreen();
        printDuck();
        System.out.println("🦆 嘎嘎！我是小鸭子助手，很高兴见到你！");
        sleep(1000);
        System.out.println("🦆 我可以帮你做很多事情哦～");
        sleep(1000);
        System.out.println("🦆 让我们开始吧！");
        sleep(1500);
    }
    
    private void showMainMenu() {
        clearScreen();
        printDuck();
        System.out.println("🦆 嘎嘎！你想让我帮你做什么呢？");
        System.out.println();
        System.out.println("1. 📊 统计当前目录下的代码行数");
        System.out.println("2. 🎁 启动抢红包游戏");
        System.out.println("3. 🤖 询问AI助手");
        System.out.println("4. 🚪 退出");
        System.out.println();
        System.out.print("🦆 请选择 (1-4): ");
    }
    
    private void handleUserInput() {
        try {
            String input = scanner.nextLine().trim();
            
            switch (input) {
                case "1":
                    handleCodeStats();
                    break;
                case "2":
                    handleRedPacketGame();
                    break;
                case "3":
                    handleAIQuery();
                    break;
                case "4":
                    isRunning = false;
                    break;
                default:
                    System.out.println("🦆 嘎嘎！我不太明白，请输入1-4之间的数字哦～");
                    sleep(1500);
            }
        } catch (Exception e) {
            System.out.println("🦆 嘎嘎！出错了，让我重新来～");
            sleep(1500);
        }
    }
    
    private void handleCodeStats() {
        System.out.println("🦆 嘎嘎！让我来帮你统计代码行数～");
        sleep(1000);

        System.out.print("🦆 请输入要统计的目录（直接回车默认当前目录）: ");
        String dirInput = scanner.nextLine().trim();
        if (dirInput.isEmpty()) {
            dirInput = ".";
        }
        dirInput = expandHomeDirectory(dirInput);
        File targetDir = new File(dirInput);

        if (!targetDir.exists()) {
            System.out.println("🦆 嘎嘎！目录不存在: " + targetDir.getPath());
        } else if (!targetDir.isDirectory()) {
            System.out.println("🦆 嘎嘎！看起来这不是一个目录: " + targetDir.getPath());
        } else {
            String dirToAnalyze = targetDir.getAbsolutePath();
            try {
                dirToAnalyze = targetDir.getCanonicalPath();
            } catch (IOException ignored) {
            }
            System.out.println("🦆 嘎嘎！马上开始统计 " + dirToAnalyze + " 下的代码～");
            sleep(800);

            try {
                ProcessBuilder pb = new ProcessBuilder(ToolPaths.codeStatsExecutable(), dirToAnalyze);
                pb.directory(new File(System.getProperty("user.dir")));
                Process process = pb.start();

                try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()))) {
                    String line;
                    System.out.println("\n=== 代码统计结果 ===\n");
                    while ((line = reader.readLine()) != null) {
                        System.out.println(line);
                    }
                }

                int exitCode = process.waitFor();
                if (exitCode == 0) {
                    System.out.println("\n🦆 嘎嘎！统计完成啦！");
                } else {
                    System.out.println("🦆 嘎嘎！统计时遇到了一些问题...");
                }
            } catch (Exception e) {
                System.out.println("🦆 嘎嘎！统计工具出错了，可能是code_stats文件不存在...");
                System.out.println("请确保已经编译了C++代码统计工具！");
            }
        }

        System.out.println("\n按回车键继续...");
        scanner.nextLine();
    }
    
    private void handleRedPacketGame() {
        System.out.println("🦆 嘎嘎！让我们来玩抢红包游戏吧！");
        sleep(1000);
        System.out.println("🦆 你想要图形界面还是命令行界面呢？");
        System.out.println("1. 🖥️ 图形界面 (推荐)");
        System.out.println("2. 💻 命令行界面");
        System.out.print("🦆 请选择 (1-2): ");
        
        String choice = scanner.nextLine().trim();
        
        if ("1".equals(choice)) {
            System.out.println("🦆 嘎嘎！启动图形界面游戏...");
            sleep(1000);
            try {
                // 启动图形界面游戏
                String[] args = {"--gui", "--count=20", "--duration=15000"};
                Main.main(args);
            } catch (Exception e) {
                System.out.println("🦆 嘎嘎！游戏启动出错了...");
            }
        } else if ("2".equals(choice)) {
            System.out.println("🦆 嘎嘎！启动命令行游戏...");
            sleep(1000);
            try {
                // 启动命令行游戏
                String[] args = {"--count=20", "--duration=10000"};
                Main.main(args);
            } catch (Exception e) {
                System.out.println("🦆 嘎嘎！游戏启动出错了...");
            }
        } else {
            System.out.println("🦆 嘎嘎！无效选择，返回主菜单...");
            sleep(1500);
            return;
        }
        
        System.out.println("\n🦆 嘎嘎！游戏结束，回到主菜单...");
        sleep(2000);
    }
    
    private void handleAIQuery() {
        if (!aiEnabled) {
            System.out.println("🦆 嘎嘎！AI服务未配置，我先用离线大脑陪你聊聊～");
            fallbackChatOnce();
            return;
        }

        System.out.println("🦆 嘎嘎！AI助手已连接，和唐老鸭聊聊吧！（输入 'back' 返回主菜单）");
        while (true) {
            System.out.print("🦆 你: ");
            String query = scanner.nextLine().trim();
            if ("back".equalsIgnoreCase(query)) {
                break;
            }
            if (query.isEmpty()) {
                continue;
            }
            conversation.add(new Message("user", query));
            try {
                String reply = aiClient.chat(conversation);
                if (reply != null && !reply.isBlank()) {
                    reply = reply.trim();
                    conversation.add(new Message("assistant", reply));
                    System.out.println("🦆 唐老鸭: " + reply);
                } else {
                    System.out.println("🦆 唐老鸭: （AI没有返回内容，换个问题吧？）");
                }
            } catch (Exception e) {
                System.out.println("🦆 唐老鸭: (AI error: " + e.getMessage() + ")");
                String fallback = generateDuckResponse(query);
                System.out.println("🦆 唐老鸭: " + fallback);
            }
        }
        System.out.println("\n🦆 嘎嘎！已退出 AI 对话，按回车返回主菜单...");
        scanner.nextLine();
    }

    private void fallbackChatOnce() {
        System.out.println("🦆 你想问我什么呢？(输入 'back' 返回主菜单)");
        System.out.print("🦆 你: ");
        String query = scanner.nextLine().trim();
        if ("back".equalsIgnoreCase(query)) {
            return;
        }
        String response = generateDuckResponse(query);
        System.out.println("🦆 小鸭子: " + response);
        sleep(1500);

        System.out.println("\n🦆 嘎嘎！还想问什么吗？(输入 'back' 返回主菜单)");
        System.out.print("🦆 你: ");
        String followUp = scanner.nextLine().trim();
        if (!"back".equalsIgnoreCase(followUp)) {
            String followUpResponse = generateDuckResponse(followUp);
            System.out.println("🦆 小鸭子: " + followUpResponse);
            sleep(1500);
        }

        System.out.println("\n按回车键返回主菜单...");
        scanner.nextLine();
    }
    
    private String generateDuckResponse(String query) {
        String lowerQuery = query.toLowerCase();
        
        if (lowerQuery.contains("你好") || lowerQuery.contains("hello")) {
            return "嘎嘎！你好！很高兴见到你！";
        } else if (lowerQuery.contains("名字") || lowerQuery.contains("name")) {
            return "嘎嘎！我是小鸭子助手，你可以叫我小鸭鸭！";
        } else if (lowerQuery.contains("功能") || lowerQuery.contains("能做什么")) {
            return "嘎嘎！我可以帮你统计代码、玩游戏，还能和你聊天呢！";
        } else if (lowerQuery.contains("代码") || lowerQuery.contains("编程")) {
            return "嘎嘎！编程很有趣呢！我可以帮你统计代码行数哦！";
        } else if (lowerQuery.contains("游戏") || lowerQuery.contains("红包")) {
            return "嘎嘎！抢红包游戏很好玩呢！要不要试试？";
        } else if (lowerQuery.contains("谢谢") || lowerQuery.contains("thank")) {
            return "嘎嘎！不客气！能帮到你我很开心！";
        } else if (lowerQuery.contains("再见") || lowerQuery.contains("bye")) {
            return "嘎嘎！再见！记得常来找我玩哦！";
        } else {
            return "嘎嘎！这个问题很有趣呢！不过我还不太懂，让我再学习学习！";
        }
    }
    
    private void sayGoodbye() {
        clearScreen();
        printDuck();
        System.out.println("🦆 嘎嘎！再见啦！");
        sleep(1000);
        System.out.println("🦆 记得常来找我玩哦～");
        sleep(1000);
        System.out.println("🦆 嘎嘎嘎嘎！");
    }
    
    private void printDuck() {
        System.out.println("    🦆");
        System.out.println("   /|\\");
        System.out.println("  / | \\");
        System.out.println(" /  |  \\");
        System.out.println("    |");
        System.out.println("   / \\");
        System.out.println("  /   \\");
        System.out.println(" /     \\");
        System.out.println();
    }
    
    private void clearScreen() {
        try {
            if (System.getProperty("os.name").toLowerCase().contains("windows")) {
                new ProcessBuilder("cmd", "/c", "cls").inheritIO().start().waitFor();
            } else {
                new ProcessBuilder("clear").inheritIO().start().waitFor();
            }
        } catch (Exception e) {
            // 如果清屏失败，就打印一些空行
            for (int i = 0; i < 50; i++) {
                System.out.println();
            }
        }
    }
    
    private void sleep(int milliseconds) {
        try {
            Thread.sleep(milliseconds);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    private String expandHomeDirectory(String path) {
        if (path.startsWith("~")) {
            String home = System.getProperty("user.home");
            if (home != null && !home.isBlank()) {
                if (path.length() == 1) {
                    return home;
                }
                char next = path.charAt(1);
                if (next == '/' || next == '\\') {
                    return home + path.substring(1);
                }
            }
        }
        return path;
    }

    private void initAiClient() {
        aiClient = AiClient.fromEnv();
        conversation.clear();
        if (aiClient != null) {
            aiEnabled = true;
            conversation.add(new Message("system", DONALD_PROMPT));
        } else {
            aiEnabled = false;
        }
    }
}
