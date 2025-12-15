# tools/code-stats

Self-contained C++17 implementation of the `code_stats` binary used by both assistants.

- `code_stats.cpp` – the entire analyzer (language detection, comment counting, TSV output)
- `Makefile` – builds the binary into this folder and exposes `make`, `make clean`, `make test`
- `code_stats` – prebuilt binary (recreate it after changes with `make`)

The Duck Assistant scripts invoke this binary via `infra.ToolPaths`, and the GUI expects TSV output so it can render charts. See `docs/code-stats.md` for full CLI usage details.
