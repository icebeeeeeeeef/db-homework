# scripts/

Wrapper scripts for local development:

- `start_duck.sh` – builds the native `code_stats` helper, compiles all Java sources into `build/classes`, then lets you choose GUI or CLI mode
- `start_duck_gui.sh` – same as above but goes straight into the GUI assistant

Both scripts assume they are run from the repository root and keep the workspace clean by writing build outputs under `build/`.
