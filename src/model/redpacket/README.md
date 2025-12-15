# model/redpacket

Data structures that the game engine uses:

- `Player` – stores current position/radius and exposes collision checks
- `RedPacket` – procedural generation of packet amount, velocity, and shape/size
- `RedPacketShape` / `RedPacketSize` – enums that describe rendering + collision properties
- `RedPacketStatistics` – collects summary metrics for end-of-game dialogs

Separating these from duck-specific models keeps the engine focused on gameplay math.
