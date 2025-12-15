# model/duck/customization

Outfits use the decorator pattern so we can layer accessories without branching everywhere:

- `DuckAppearance` – contract for painting a duck given a graphics context
- `BaseDuck` – draws the base sprite
- `AccessoryDecorator` + concrete decorators (`HatDecorator`, `ScarfDecorator`, `EyeDecorator`, `TieDecorator`, `CaneDecorator`) – toggleable accessories
- `DuckOutfit` – mutable configuration object that knows which decorators to apply and stores per-accessory colors

The GUI editor and `StagePanel` both stick to `DuckOutfit`, while rendering is handled by the appearance/decorator graph.
