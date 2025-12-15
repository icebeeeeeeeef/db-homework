# model/duck/behavior

Defines how each duck moves and sounds:

- `ActionBehavior` / `SoundBehavior` – small strategy interfaces for actions and quacks
- `FlyAction`, `RunAction`, `SwimAction` – sample action implementations
- `QuackSound`, `ChirpSound`, `WhistleSound` – sound implementations
- `DuckBehaviorProfile` – pairs an action + sound for a specific character
- `BehaviorLibrary` – exposes curated lists for the GUI dialogs

The GUI dialog simply picks behaviors from these lists and stores them in `StagePanel` without duplicating logic.
