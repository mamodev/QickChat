import { useEffect, useState } from "react";
import { z } from "zod";

export const chatSchema = z
  .object({
    id: z.string(),
    name: z.string(),
  })
  .array();

export type Chat = z.infer<typeof chatSchema>[number];

let selectedChat: Chat | null = null;

const listeners = new Set<(selected: Chat) => void>();

export function addChatListener(listener: (selected: Chat) => void) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function setSelectedChat(chat: Chat) {
  selectedChat = chat;
  listeners.forEach((listener) => listener(chat));
}

export function useSelectedChat() {
  const [chat, setChat] = useState(selectedChat);

  useEffect(() => {
    const cleanUp = addChatListener(setChat);
    return () => {
      cleanUp();
    };
  }, []);

  return chat;
}

export function getSelectedChat() {
  return selectedChat;
}
