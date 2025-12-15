package ai;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * 简单的 OpenAI 兼容 Chat Completions 客户端实现。
 */
public class AiClient {
    public static class Message {
        public final String role;
        public final String content;

        public Message(String role, String content) {
            this.role = role;
            this.content = content;
        }
    }

    private final HttpClient http;
    private final String baseUrl;
    private final String model;
    private final String apiKey;

    private AiClient(HttpClient http, String baseUrl, String model, String apiKey) {
        this.http = http;
        this.baseUrl = baseUrl;
        this.model = model;
        this.apiKey = apiKey;
    }

    public static AiClient fromEnv() {
        String key = System.getenv("OPENAI_API_KEY");
        if (key == null || key.isBlank()) {
            return null;
        }
        String base = System.getenv("OPENAI_BASE_URL");
        if (base == null || base.isBlank()) {
            base = "https://api.openai.com";
        } else if (base.endsWith("/")) {
            base = base.substring(0, base.length() - 1);
        }
        String model = System.getenv("OPENAI_MODEL");
        if (model == null || model.isBlank()) {
            model = "gpt-3.5-turbo";
        }
        HttpClient client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(10))
                .build();
        return new AiClient(client, base, model, key);
    }

    public String chat(List<Message> messages) throws IOException, InterruptedException {
        Objects.requireNonNull(messages, "messages");
        if (messages.isEmpty()) {
            throw new IllegalArgumentException("messages must not be empty");
        }
        String body = buildRequestBody(messages);
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/v1/chat/completions"))
                .timeout(Duration.ofSeconds(30))
                .header("Authorization", "Bearer " + apiKey)
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(body, StandardCharsets.UTF_8))
                .build();

        HttpResponse<String> response = http.send(request, HttpResponse.BodyHandlers.ofString(StandardCharsets.UTF_8));
        if (response.statusCode() < 200 || response.statusCode() >= 300) {
            throw new IOException("AI request failed: status=" + response.statusCode() + " body=" + response.body());
        }
        return parseContent(response.body());
    }

    private String buildRequestBody(List<Message> messages) {
        StringBuilder builder = new StringBuilder();
        builder.append("{\"model\":\"")
                .append(escapeJson(model))
                .append("\",\"messages\":[");
        List<String> payload = new ArrayList<>();
        for (Message m : messages) {
            String msg = String.format("{\"role\":\"%s\",\"content\":\"%s\"}",
                    escapeJson(m.role), escapeJson(m.content));
            payload.add(msg);
        }
        builder.append(String.join(",", payload));
        builder.append("]}");
        return builder.toString();
    }

    private String parseContent(String responseBody) throws IOException {
        Pattern pattern = Pattern.compile("\\\"content\\\"\\s*:\\s*\\\"(.*?)\\\"", Pattern.DOTALL);
        Matcher matcher = pattern.matcher(responseBody);
        if (matcher.find()) {
            return unescapeJson(matcher.group(1)).trim();
        }
        throw new IOException("AI response missing content: " + responseBody);
    }

    private String escapeJson(String text) {
        if (text == null) return "";
        return text.replace("\\", "\\\\")
                .replace("\"", "\\\"")
                .replace("\n", "\\n")
                .replace("\r", "\\r");
    }

    private String unescapeJson(String text) {
        if (text == null) return "";
        return text
                .replace("\\r", "\r")
                .replace("\\n", "\n")
                .replace("\\\"", "\"")
                .replace("\\\\", "\\");
    }
}
