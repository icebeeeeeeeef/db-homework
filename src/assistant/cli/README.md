# assistant/cli

Contains `DuckAssistant`, a whimsical terminal assistant that:

- invokes the shared `tools/code-stats/code_stats` binary (via `infra.ToolPaths`) to print quick code metrics
- launches either the GUI or console version of the red packet game
- proxies to the Swing GUI by calling `app.Main` so both versions share the same configuration parsing

The class sticks to CLI responsibilities and delegates gameplay/AI work to the GUI and engine packages.
