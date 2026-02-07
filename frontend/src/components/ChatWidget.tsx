"use client";

import React, { useMemo, useState } from "react";

type ChatMessage = {
  role: "user" | "assistant";
  content: string;
};

type ChatResponse = {
  answer: string;
  sources?: string[];
};

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE ||
  process.env.NEXT_PUBLIC_RAG_API_BASE ||
  "http://localhost:8080";

export default function ChatWidget() {
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<ChatMessage[]>(() => [
    {
      role: "assistant",
      content:
        "Hi. Ask me about War Tiket Engine, setup steps, or system behavior.",
    },
  ]);

  const canSend = useMemo(() => input.trim().length > 0 && !loading, [input, loading]);

  const sendMessage = async () => {
    const text = input.trim();
    if (!text || loading) return;

    setError("");
    setLoading(true);
    setInput("");

    setMessages((prev) => [...prev, { role: "user", content: text }]);

    try {
      const res = await fetch(`${API_BASE}/api/chat`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: text }),
      });
      const data = (await res.json()) as ChatResponse;
      const answer = data?.answer || "No response from the assistant.";
      setMessages((prev) => [...prev, { role: "assistant", content: answer }]);
    } catch {
      setError("Failed to reach the assistant service.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="chat-widget-root">
      <button
        className="btn btn-primary chat-widget-toggle shadow"
        onClick={() => setOpen((v) => !v)}
        type="button"
      >
        {open ? "Close Assistant" : "Ask Assistant"}
      </button>

      {open && (
        <div className="chat-widget-panel shadow-lg">
          <div className="chat-widget-header">
            <div>
              <div className="fw-bold">War Tiket Assistant</div>
              <div className="small text-muted">RAG on README</div>
            </div>
            <button
              className="btn btn-sm btn-outline-secondary"
              onClick={() => setOpen(false)}
              type="button"
            >
              Close
            </button>
          </div>

          <div className="chat-widget-body">
            {messages.map((msg, idx) => (
              <div
                key={`${msg.role}-${idx}`}
                className={`chat-bubble ${msg.role === "user" ? "chat-user" : "chat-assistant"}`}
              >
                {msg.content}
              </div>
            ))}
            {loading && (
              <div className="chat-bubble chat-assistant">Typing...</div>
            )}
          </div>

          <div className="chat-widget-footer">
            {error && <div className="text-danger small mb-2">{error}</div>}
            <div className="input-group">
              <input
                className="form-control"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="Type your question..."
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    sendMessage();
                  }
                }}
              />
              <button
                className="btn btn-success"
                onClick={sendMessage}
                disabled={!canSend}
                type="button"
              >
                Send
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
