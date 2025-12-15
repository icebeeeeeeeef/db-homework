# game/core

Red packet gameplay is broken into four focused classes:

- `GameConfig` – parses CLI flags and stores tunable runtime values (board size, duration, player radius, etc.)
- `GameEngine` – owns the update loop, collision checks, and summary printing
- `Renderer` – draws the ASCII arena according to the player's position and each `model.redpacket` entity
- `InputController` – listens for keyboard input on a background thread and exposes the latest direction vector

This package deliberately stays free of Swing/GUI code so both CLI and GUI layers can reuse the same mechanics.
