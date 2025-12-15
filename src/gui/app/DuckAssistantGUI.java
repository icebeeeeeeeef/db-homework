package gui;

import ai.AiClient;
import ai.AiClient.Message;
import attendance.ui.AttendanceApp;
import model.duck.DuckCharacter;
import model.duck.customization.DuckOutfit;
import model.duck.behavior.ActionBehavior;
import model.duck.behavior.SoundBehavior;

import javax.swing.*;
import java.awt.*;
import java.util.ArrayList;
import java.util.List;
import java.io.File;

/**
 * 图形化小鸭子助手界面
 */
public class DuckAssistantGUI extends JFrame {
    private static final float BASE_FONT_SIZE = 22f;
    private static final float SMALL_FONT_SIZE = 18f;
    private static final String DONALD_PROMPT = "You are Donald Duck in a game world. " +
            "Speak in Donald Duck's lively, humorous style. Keep responses concise (within 80 characters when possible) " +
            "and reference the player's adventures. If the player speaks Chinese, reply in Chinese but maintain your Donald Duck personality.";

    private JTextArea chatArea;
    private JTextField inputField;
    private JButton sendButton;
    private JButton codeStatsButton;
    private JButton redPacketButton;
    private JButton attendanceButton;
    private JButton aiChatButton;
    private JButton behaviorButton;
    private JButton customizeButton;
    private StagePanel stagePanel;
    private StageCommandHandler stageCommandHandler;
    private JPanel mainPanel;
    private JPanel chatPanel;
    private JPanel buttonPanel;
    private JPanel duckHeader;
    
    private final List<Message> conversation = new ArrayList<>();
    private AiClient aiClient;
    private boolean aiEnabled;

    public DuckAssistantGUI() {
        initializeComponents();
        setupLayout();
        setupEventHandlers();
        setupFrame();
        initAiClient();
    }
    
    private void initializeComponents() {
        // 主面板
        mainPanel = new JPanel(new BorderLayout());
        mainPanel.setBackground(new Color(240, 248, 255)); // 淡蓝色背景
        
        // 聊天区域
        chatArea = new JTextArea(18, 50);
        chatArea.setEditable(false);
        chatArea.setFont(font(BASE_FONT_SIZE, Font.PLAIN));
        chatArea.setBackground(new Color(255, 255, 240)); // 淡黄色背景
        chatArea.setBorder(BorderFactory.createEmptyBorder(10, 10, 10, 10));
        chatArea.setLineWrap(true);
        chatArea.setWrapStyleWord(true);
        
        JScrollPane chatScrollPane = new JScrollPane(chatArea);
        chatScrollPane.setVerticalScrollBarPolicy(JScrollPane.VERTICAL_SCROLLBAR_AS_NEEDED);
        javax.swing.border.TitledBorder border = BorderFactory.createTitledBorder("🦆 Duck Chat");
        border.setTitleFont(font(16f, Font.BOLD));
        chatScrollPane.setBorder(border);
        
        // 输入区域
        JPanel inputPanel = new JPanel(new BorderLayout());
        inputField = new JTextField();
        inputField.setFont(font(BASE_FONT_SIZE, Font.PLAIN));
        inputField.setBorder(BorderFactory.createEmptyBorder(5, 5, 5, 5));
        
        sendButton = new JButton("Send");
        sendButton.setFont(font(SMALL_FONT_SIZE, Font.BOLD));
        sendButton.setBackground(new Color(255, 182, 193)); // 粉色
        sendButton.setForeground(Color.WHITE);
        sendButton.setFocusPainted(false);
        
        inputPanel.add(inputField, BorderLayout.CENTER);
        inputPanel.add(sendButton, BorderLayout.EAST);
        
        // 聊天面板
        chatPanel = new JPanel(new BorderLayout());
        chatPanel.add(chatScrollPane, BorderLayout.CENTER);
        chatPanel.add(inputPanel, BorderLayout.SOUTH);
        
        // 功能按钮面板
        buttonPanel = new JPanel(new GridLayout(4, 2, 16, 16));
        buttonPanel.setBorder(BorderFactory.createEmptyBorder(24, 24, 24, 24));
        buttonPanel.setBackground(new Color(240, 248, 255));
        
        // 代码统计按钮
        codeStatsButton = createStyledButton("📊 Code Stats", new Color(135, 206, 235));
        
        // 抢红包游戏按钮
        redPacketButton = createStyledButton("🎁 Red Packet", new Color(255, 99, 71));

        // 课堂点名按钮
        attendanceButton = createStyledButton("📋 Attendance", new Color(46, 139, 87));
        
        // AI对话按钮
        aiChatButton = createStyledButton("🤖 AI Helper", new Color(147, 112, 219));
        
        // 装扮按钮
        customizeButton = createStyledButton("🧢 Dress Up", new Color(60, 179, 113));
        // 行为按钮
        behaviorButton = createStyledButton("⚡ Actions", new Color(255, 140, 0));

        // 退出按钮
        JButton exitButton = createStyledButton("🚪 Exit", new Color(220, 20, 60));
        
        buttonPanel.add(codeStatsButton);
        buttonPanel.add(redPacketButton);
        buttonPanel.add(attendanceButton);
        buttonPanel.add(customizeButton);
        buttonPanel.add(aiChatButton);
        buttonPanel.add(behaviorButton);
        buttonPanel.add(exitButton);
        
        // 小舞台
        stagePanel = new StagePanel();
        initializeDefaultOutfits();
        stageCommandHandler = new StageCommandHandler(this, stagePanel, this::appendToChat, this::runAiCommand);
        stagePanel.setStageListener(stageCommandHandler::openStageCommandDialog);

        duckHeader = new JPanel(new BorderLayout());
        duckHeader.setOpaque(false);
        JLabel title = new JLabel("🦆 Donald & Ducklings Stage", SwingConstants.CENTER);
        title.setFont(font(BASE_FONT_SIZE + 2, Font.BOLD));
        title.setBorder(BorderFactory.createEmptyBorder(12, 12, 0, 12));
        duckHeader.add(title, BorderLayout.NORTH);
        duckHeader.add(stagePanel, BorderLayout.CENTER);
        
        // 退出按钮事件
        exitButton.addActionListener(e -> {
            showGoodbyeMessage();
            System.exit(0);
        });

        // 装扮按钮事件
        customizeButton.addActionListener(e -> DuckOutfitCustomizer.open(
                this,
                stagePanel,
                this::appendToChat,
                this::createDefaultOutfit));
    }
    
    private JButton createStyledButton(String text, Color color) {
        JButton button = new JButton(text);
        button.setFont(font(BASE_FONT_SIZE, Font.BOLD));
        button.setBackground(color);
        button.setForeground(Color.WHITE);
        button.setFocusPainted(false);
        button.setBorderPainted(false);
        button.setPreferredSize(new Dimension(280, 96));
        
        // 添加悬停效果
        button.addMouseListener(new java.awt.event.MouseAdapter() {
            public void mouseEntered(java.awt.event.MouseEvent evt) {
                button.setBackground(color.darker());
            }
            public void mouseExited(java.awt.event.MouseEvent evt) {
                button.setBackground(color);
            }
        });
        
        return button;
    }
    
    private void setupLayout() {
        // 左侧聊天区域
        JPanel leftPanel = new JPanel(new BorderLayout());
        leftPanel.add(chatPanel, BorderLayout.CENTER);
        leftPanel.setPreferredSize(new Dimension(700, 560));
        
        // 右侧功能区域
        JPanel rightPanel = new JPanel(new BorderLayout());
        rightPanel.add(duckHeader, BorderLayout.NORTH);
        rightPanel.add(buttonPanel, BorderLayout.CENTER);
        rightPanel.setPreferredSize(new Dimension(380, 560));
        
        // 主布局
        mainPanel.add(leftPanel, BorderLayout.WEST);
        mainPanel.add(rightPanel, BorderLayout.EAST);
        
        add(mainPanel);
    }
    
    private void setupEventHandlers() {
        // 发送按钮事件
        sendButton.addActionListener(e -> handleUserInput());
        
        // 输入框回车事件
        inputField.addActionListener(e -> handleUserInput());
        
        // 代码统计按钮
        codeStatsButton.addActionListener(e -> handleCodeStats());
        
        // 抢红包游戏按钮
        redPacketButton.addActionListener(e -> handleRedPacketGame());

        // 课堂点名按钮
        attendanceButton.addActionListener(e -> openAttendanceApp());
        
        // AI对话按钮
        aiChatButton.addActionListener(e -> handleAIChat());
        // 行为按钮
        behaviorButton.addActionListener(e -> openBehaviorDialog());
    }
    
    private void setupFrame() {
        setTitle("🦆 Duck Assistant");
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        setSize(1120, 660);
        setLocationRelativeTo(null);
        setResizable(true);
        
        // 设置窗口图标（如果有的话）
        try {
            setIconImage(Toolkit.getDefaultToolkit().getImage(getClass().getResource("/duck.png")));
        } catch (Exception e) {
            // 忽略图标加载错误
        }
    }
    
    private void handleUserInput() {
        String input = inputField.getText().trim();
        if (input.isEmpty()) return;
        
        // 显示用户输入
        appendToChat("You: " + input);
        inputField.setText("");
        
        if (aiEnabled) {
            conversation.add(new Message("user", input));
            requestAiResponse(input);
        } else {
            String response = generateDuckResponse(input);
            appendToChat("🦆 Duck: " + response);
        }
    }
    
    private void requestAiResponse(String userInput) {
        appendToChat("🦆 Duck: ... (thinking)");
        SwingWorker<String, Void> worker = new SwingWorker<>() {
            @Override
            protected String doInBackground() throws Exception {
                return aiClient.chat(conversation);
            }
 
            @Override
            protected void done() {
                try {
                    String reply = get();
                    if (reply != null && !reply.isBlank()) {
                        // 替换刚才提示行
                        replaceLastThinkingLine();
                        appendToChat("🦆 Duck: " + reply.trim());
                        conversation.add(new Message("assistant", reply.trim()));
                        stagePanel.setSpeech("AI response ready!");
                        return;
                    }
                } catch (Exception e) {
                    replaceLastThinkingLine();
                    appendToChat("🦆 Duck: (AI error: " + e.getMessage() + ")");
                }
                String fallback = generateDuckResponse(userInput);
                appendToChat("🦆 Duck: " + fallback);
                stagePanel.setSpeech("Fallback response ready");
            }
        };
        worker.execute();
    }
 
    private void replaceLastThinkingLine() {
        SwingUtilities.invokeLater(() -> {
            String text = chatArea.getText();
            int idx = text.lastIndexOf("🦆 Duck: ... (thinking)");
            if (idx >= 0) {
                chatArea.replaceRange("", idx, text.length());
            }
        });
    }
    
    private void handleCodeStats() {
        JFileChooser chooser = new JFileChooser(new File(System.getProperty("user.dir")));
        chooser.setFileSelectionMode(JFileChooser.DIRECTORIES_ONLY);
        chooser.setDialogTitle("Select directory for code statistics");
        int result = chooser.showOpenDialog(this);
        if (result != JFileChooser.APPROVE_OPTION) {
            appendToChat("🦆 Duck: Code stats canceled (no directory selected).");
            return;
        }
        File selected = chooser.getSelectedFile();
        if (selected == null) {
            appendToChat("🦆 Duck: Hmm, I couldn't read that directory.");
            return;
        }
        String directory = selected.getAbsolutePath();
        appendToChat("🦆 Duck: Preparing stats for " + directory);
        StageCommand.CodeStatsOptions opts =
                new StageCommand.CodeStatsOptions(directory, List.of(), true, true, true, false);
        stageCommandHandler.runCodeStats(opts);
    }
    
    private void handleRedPacketGame() {
        appendToChat("🦆 Duck: Let's play the red packet game!");
        
        // 创建游戏选择对话框
        String[] options = {"🖥️ GUI", "💻 Console"};
        int choice = JOptionPane.showOptionDialog(
            this,
            "🦆 Which mode do you prefer?",
            "Select Game Mode",
            JOptionPane.YES_NO_OPTION,
            JOptionPane.QUESTION_MESSAGE,
            null,
            options,
            options[0]
        );
        
        if (choice == 0) {
            appendToChat("🦆 Duck: Launching GUI game...");
            SwingUtilities.invokeLater(() -> {
                try {
                    GuiGame.launchInteractive();
                } catch (Exception e) {
                    appendToChat("🦆 Duck: Failed to start GUI game...");
                }
            });
        } else if (choice == 1) {
            appendToChat("🦆 Duck: Launching console game (runs in terminal)...");
            new Thread(() -> {
                try {
                    String classPath = System.getProperty("java.class.path");
                    ProcessBuilder pb = new ProcessBuilder(
                            "java", "-cp", classPath, "app.Main", "--count=20", "--duration=10000");
                    pb.directory(new java.io.File(System.getProperty("user.dir")));
                    pb.inheritIO();
                    pb.start();
                } catch (Exception e) {
                    SwingUtilities.invokeLater(() -> appendToChat("🦆 Duck: Failed to start console game..."));
                }
            }).start();
        }
    }

    private void openAttendanceApp() {
        appendToChat("🦆 Duck: Launching classroom attendance system...");
        SwingUtilities.invokeLater(() -> {
            try {
                AttendanceApp.launch();
            } catch (Exception ex) {
                appendToChat("🦆 Duck: Failed to open attendance UI: " + ex.getMessage());
            }
        });
    }
    
    private void handleAIChat() {
        appendToChat("🦆 Duck: AI helper is still under construction...");
        appendToChat("🦆 Duck: But I can chat with you right now!");
        appendToChat("🦆 Duck: Type anything on the left and hit Enter.");
    }

    private void runAiCommand(StageCommand.AiOptions options) {
        if (options == null || options.prompt == null || options.prompt.isBlank()) return;
        String prompt = options.prompt.trim();
        appendToChat("You (stage): " + prompt);
        stagePanel.setSpeech("Thinking...");
        if (aiEnabled) {
            conversation.add(new Message("user", prompt));
            requestAiResponse(prompt);
        } else {
            appendToChat("🦆 Duck: " + generateDuckResponse(prompt));
            stagePanel.setSpeech("Offline reply ready!");
        }
    }
    
    private String generateDuckResponse(String input) {
        String lowerInput = input.toLowerCase();
        
        if (lowerInput.contains("你好") || lowerInput.contains("hello")) {
            return "Quack! Nice to meet you!";
        } else if (lowerInput.contains("名字") || lowerInput.contains("name")) {
            return "You can call me Duck Assistant!";
        } else if (lowerInput.contains("功能") || lowerInput.contains("能做什么") || lowerInput.contains("help")) {
            return "I can count code, launch the game, and keep you company!";
        } else if (lowerInput.contains("代码") || lowerInput.contains("编程") || lowerInput.contains("code")) {
            return "Coding is fun! Want me to run the code stats report?";
        } else if (lowerInput.contains("游戏") || lowerInput.contains("红包") || lowerInput.contains("game")) {
            return "The red packet game is ready whenever you are!";
        } else if (lowerInput.contains("谢谢") || lowerInput.contains("thank")) {
            return "You're welcome! Happy to help.";
        } else if (lowerInput.contains("再见") || lowerInput.contains("bye")) {
            return "See you soon! Come back and chat with me again.";
        } else {
            return "Quack! That's interesting. I'm still learning, but I'm listening!";
        }
    }
    
    private void appendToChat(String message) {
        SwingUtilities.invokeLater(() -> {
            chatArea.append(message + "\n");
            chatArea.setCaretPosition(chatArea.getDocument().getLength());
        });
    }
    
    private void showGoodbyeMessage() {
        JOptionPane.showMessageDialog(
            this,
            "🦆 Quack! See you next time!\n🦆 Come back and play with me soon!",
            "Duck Assistant",
            JOptionPane.INFORMATION_MESSAGE
        );
    }

    private void initializeDefaultOutfits() {
        for (DuckCharacter character : DuckCharacter.values()) {
            stagePanel.setOutfit(character, createDefaultOutfit(character));
        }
    }

    private DuckOutfit createDefaultOutfit(DuckCharacter character) {
        DuckOutfit outfit = new DuckOutfit();
        switch (character) {
            case DONALD:
                outfit.setHat(true);
                outfit.setTie(true);
                break;
            case DUCKLING_ONE:
                outfit.setScarf(true);
                break;
            case DUCKLING_TWO:
                outfit.setHat(true);
                outfit.setEyes(true);
                break;
            case DUCKLING_THREE:
                outfit.setTie(true);
                outfit.setCane(true);
                break;
            default:
                break;
        }
        return outfit;
    }

    private Font font(float size, int style) {
        return DuckUiTheme.font(size, style);
    }

    private void openBehaviorDialog() {
        BehaviorDialog dialog = new BehaviorDialog(this);
        BehaviorDialog.Selection selection = dialog.showDialog();
        if (selection == null) return;
        stagePanel.setBehavior(selection.character, selection.profile);

        ActionBehavior action = selection.profile.getActionBehavior();
        SoundBehavior sound = selection.profile.getSoundBehavior();
        String actionText = action != null ? action.perform() : "stands still.";
        String soundText = sound != null ? sound.makeSound() : "(silence)";
        String message = String.format("%s %s and says %s",
                selection.character.getDisplayName(), actionText, soundText);
        appendToChat("🦆 Duck: " + message);
        String speech = selection.character.getDisplayName() + (action != null ? " " + action.getName() : "");
        stagePanel.setSpeech(speech.trim());
    }
    
    public static void launch() {
        SwingUtilities.invokeLater(() -> {
            try {
                UIManager.setLookAndFeel(UIManager.getSystemLookAndFeelClassName());
            } catch (Exception e) {
                // 忽略主题设置错误
            }
            
            new DuckAssistantGUI().setVisible(true);
        });
    }

    private void initAiClient() {
        aiClient = AiClient.fromEnv();
        if (aiClient != null) {
            aiEnabled = true;
            conversation.clear();
            conversation.add(new Message("system", DONALD_PROMPT));
            appendToChat("🦆 Duck: AI link ready! Let's chat.");
        } else {
            aiEnabled = false;
            appendToChat("🦆 Duck: AI service not configured. Using offline brain.");
        }
    }
}
