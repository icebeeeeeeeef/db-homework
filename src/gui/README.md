# gui/

Swing UI code is grouped here and organized by responsibility:

- `app/` – hosts `DuckAssistantGUI`, the main frame wiring chat, AI, and stage widgets.
- `core/` – stage orchestration utilities such as `StageCommandHandler`, `StageCommand`, and `GameOverListener`.
- `game/` – rendering-focused components (`GuiGame`, `StagePanel`, `RainGameWindow`) that drive the red packet experience.
- `dialogs/` – interaction surfaces including `BehaviorDialog`, `StageCommandDialog`, `StartDialog`, and `ResultDialog`.
- `analytics/` – statistics models and charts (`CodeStatsChartPanel`, `FunctionStatsPanel`, `CodeStatsResult`).
- `customization/` – duck appearance tooling like `DuckOutfitCustomizer`, `DuckAvatarPanel`, and `DuckUiTheme`.

Every helper still lives in the `gui` package so components can collaborate without exposing Swing details to other modules.
