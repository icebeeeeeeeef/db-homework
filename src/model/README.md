# model/

Domain models are split into duck-centric and red-packet-centric namespaces:

- `duck/` – characters, behaviors, and appearance decorators used by the stage and dress-up UI
- `redpacket/` – data objects that power the game engine (player, packet shapes/sizes/statistics)

This keeps gameplay data models separate from cosmetic duck logic, which makes the code easier to reason about.
