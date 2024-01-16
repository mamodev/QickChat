import { z } from "zod";
import { safeFetch } from "../hooks/fetch";
import { useEffect, useState } from "react";
import { getSelectedChat, setSelectedChat } from "./chatContext";

export const messageSchema = z.object({
  id: z.string(),
  chat_id: z.string(),
  message: z.string(),
  timestamp: z.string(),
  sender_id: z.string(),
  sender_username: z.string(),
});

export type Message = z.infer<typeof messageSchema>;

type ChatMessageListener = (messages: Message[], hasMore: boolean) => void;

function removeDuplicateMessages(messages: Message[]) {
  const seen = new Set<string>();
  const uniqueMessages = [];

  for (let i = messages.length - 1; i >= 0; i--) {
    const message = messages[i];
    if (!seen.has(message.id)) {
      uniqueMessages.push(message);
      seen.add(message.id);
    }
  }

  return uniqueMessages;
}

function cleanUpMessages(messages: Message[]) {
  return removeDuplicateMessages(
    messages.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
  );
}

class ChatMessages {
  private messages: Message[] = [];
  private listeners: ChatMessageListener[] = [];

  private chat_id: string;

  private fetching = false;
  private latestTimestamp: Date | null = null;
  private hasMore = true;
  private updates = 0;

  constructor(chat_id: string) {
    this.chat_id = chat_id;
  }

  public fetchMessages = async () => {
    if (this.fetching) return;
    this.fetching = true;

    let params = "";
    if (this.latestTimestamp) {
      params = `?before_time=${this.latestTimestamp.toISOString()}`;
    }

    const messages = await safeFetch(`/api/chat/${this.chat_id}/message${params}`, messageSchema.array());

    this.messages = cleanUpMessages([...this.messages, ...messages]);

    if (messages.length === 0) {
      this.hasMore = false;
      this.fetching = false;
      this.notifyListeners();
      return;
    }

    this.latestTimestamp = new Date(messages[messages.length - 1].timestamp);
    this.fetching = false;
    this.notifyListeners();
  };

  public pushMessage(message: Message) {
    this.updates++;

    if (getSelectedChat()?.id !== this.chat_id) {
      const notification = new Notification(message.sender_username, {
        body: message.message,
      });

      notification.onclick = () => {
        setSelectedChat({
          id: this.chat_id,
          name: message.sender_username,
        });
        notification.close();
      };

      notification.onclose = () => {
        notification.close();
      };
    }

    this.messages = cleanUpMessages([...this.messages, message]);
    this.notifyListeners();
  }

  public addListener = (listener: ChatMessageListener) => {
    this.listeners.push(listener);
    listener(this.messages, this.hasMore);
    return () => {
      this.listeners = this.listeners.filter((l) => l !== listener);
    };
  };

  public notifyListeners = () => {
    this.listeners.forEach((listener) => listener(this.messages, this.hasMore));
  };
}

const chatPool = new Map<string, ChatMessages>();

const getChatMessages = (chat_id: string) => {
  if (!chatPool.has(chat_id)) {
    chatPool.set(chat_id, new ChatMessages(chat_id));
  }

  return chatPool.get(chat_id)!;
};

export function addChatMessage(chat_id: string, message: Message) {
  if (chatPool.has(chat_id)) {
    getChatMessages(chat_id).pushMessage(message);
  }
}

export function useChatMessages(chat_id: string) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    const cm = getChatMessages(chat_id);

    const listener = (messages: Message[], hasMore: boolean) => {
      setMessages(messages);
      setHasMore(hasMore);
      if (hasMore && messages.length === 0) cm.fetchMessages();
    };

    const removeListener = cm.addListener(listener);

    return () => {
      removeListener();
    };
  }, [chat_id]);

  return {
    messages,
    hasMore,
    fetchMore: () => getChatMessages(chat_id).fetchMessages(),
  };
}
