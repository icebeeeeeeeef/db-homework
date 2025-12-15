# src/

All Java sources live under `src/` and are split by responsibility instead of keeping everything in the root. Each immediate child folder has its own README with more detail, but the high-level grouping is:

- `app/` – the packaged entry point (`app.Main`) that wires CLI, GUI, and the game engine together
- `assistant/` – interactive helpers (currently a CLI companion under `assistant/cli`)
- `ai/` – lightweight OpenAI-compatible client so the GUI can talk to LLMs without mixing UI with networking code
- `game/` – the real-time red packet gameplay core (engine, renderer, and input loop)
- `gui/` – all Swing UI, now split into multiple focused classes (main window, dialogs, charts, handlers, customizers)
- `model/` – shared data objects, separated into duck customization/behavior vs. red-packet gameplay models
- `geom/` – math helpers such as `Vector2`
- `infra/` – shared infrastructure utilities (currently common tool paths)
