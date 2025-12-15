# model/duck

Duck-specific models are grouped here:

- `DuckCharacter` – enum for Donald + three ducklings used throughout the GUI
- `behavior/` – action & sound behavior interfaces plus ready-made implementations
- `customization/` – decorator-based appearance system (base duck with hats, scarves, eyes, etc.)

Both the GUI stage and the CLI assistant can mix behaviors and outfits without touching red-packet gameplay classes.
