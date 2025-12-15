# app/

`app.Main` is the single entry point for the entire project. It inspects CLI flags and routes execution to one of three modes:

- `--duck-gui` launches the Swing based Duck Assistant (`gui.DuckAssistantGUI`)
- `--duck` starts the terminal assistant from `assistant.cli.DuckAssistant`
- default mode spins up the red packet game engine (GUI or console based on `--gui`)

Keeping the entrypoint inside its own package avoids circular dependencies and clarifies where startup logic lives.
