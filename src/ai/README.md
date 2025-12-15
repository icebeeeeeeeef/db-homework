# ai/

`AiClient` encapsulates the OpenAI-compatible HTTP calls used by the GUI. It reads API credentials from environment variables and exposes a simple `chat(List<Message>)` method so Swing code can stay asynchronous and testable.
