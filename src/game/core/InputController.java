package game.core;

import java.io.IOException;

public class InputController implements Runnable {
    private volatile int dx = 0;
    private volatile int dy = 0;
    private volatile boolean quit = false;
    private volatile boolean running = true;

    public int getDx() { return dx; }
    public int getDy() { return dy; }
    public boolean isQuit() { return quit; }
    public void stop() { running = false; }

    @Override
    public void run() {
        try {
            while (running) {
                int c = System.in.read();
                if (c == -1) {
                    // 输入结束
                    break;
                }
                // 方向键为 ESC [ A/B/C/D
                if (c == 27) { // ESC
                    System.in.mark(2);
                    int c1 = System.in.read();
                    int c2 = System.in.read();
                    if (c1 == 91) { // [
                        if (c2 == 'A') { // Up
                            dx = 0; dy = -1;
                        } else if (c2 == 'B') { // Down
                            dx = 0; dy = 1;
                        } else if (c2 == 'C') { // Right
                            dx = 1; dy = 0;
                        } else if (c2 == 'D') { // Left
                            dx = -1; dy = 0;
                        }
                    } else {
                        // 非期望序列，回退
                        try { System.in.reset(); } catch (IOException ignored) {}
                    }
                } else {
                    char ch = (char) c;
                    switch (Character.toLowerCase(ch)) {
                        case 'w': dx = 0; dy = -1; break;
                        case 's': dx = 0; dy = 1;  break;
                        case 'a': dx = -1; dy = 0; break;
                        case 'd': dx = 1;  dy = 0; break;
                        case 'q': quit = true; break;
                        default: break;
                    }
                }
            }
        } catch (IOException ignored) {
        }
    }
}

