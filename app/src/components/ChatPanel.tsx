import React, { useState } from "react";
import InfiniteScroll from "react-infinite-scroll-component";
import { useUserID } from "../context/authContext";
import { Message, addChatMessage, useChatMessages } from "../context/chat";
import { Chat, useSelectedChat } from "../context/chatContext";
import { useMutation } from "../hooks/useMutation";
import { MessageComponent } from "./Message";

export function ChatPanel() {
  const chat = useSelectedChat();

  return (
    <div
      className="flex stack bg-primary-50 "
      style={{
        position: "relative",
        overflow: "hidden",
      }}
    >
      {!chat && (
        <div className="row v-center h-center f-height f-width">
          <p className="title">Seleziona una chat</p>
        </div>
      )}

      {chat && <ChatBody chat={chat} />}
    </div>
  );
}

type ChatBodyProps = {
  chat: Chat;
};

function ChatBody({ chat }: ChatBodyProps) {
  const { messages, fetchMore, hasMore } = useChatMessages(chat.id);

  // infinite scroll
  const onFetchMore = () => {
    fetchMore();
  };

  return (
    <div className="stack flex" style={{ overflow: "hidden", maxHeight: "100%" }}>
      <div className="row v-center spacing-1 bg-primary-950 px-sm">
        <img src={`api/chat/${chat.id}/avatar`} className="rounded p-sm" alt={chat.name} width={50} height={50} />
        <p className="title bold c-white">{chat.name}</p>
      </div>
      <ChatMessages messages={messages} onFetchMore={onFetchMore} hasMore={hasMore} />
      <ChatForm chatId={chat.id} />
    </div>
  );
}

type ChatMessageProps = {
  messages: Message[];
  onFetchMore: () => void;
  hasMore: boolean;
};

function ChatMessages(props: ChatMessageProps) {
  const id = useUserID();
  const { messages, onFetchMore, hasMore } = props;

  return (
    <div
      id="scrollableDiv"
      style={{
        overflow: "auto",
        display: "flex",
        flexDirection: "column-reverse",
        paddingBottom: "5rem",
      }}
    >
      <InfiniteScroll
        dataLength={messages.length}
        next={onFetchMore}
        className="stack flex reverse spacing-1 px-md"
        inverse={true}
        hasMore={hasMore}
        loader={<h4>Loading...</h4>}
        scrollableTarget="scrollableDiv"
      >
        {messages.map((message) => (
          <MessageComponent key={message.id} message={message} userId={id} />
        ))}
      </InfiniteScroll>
    </div>
  );
}

type ChatFormProps = {
  chatId: string;
};

function ChatForm({ chatId }: ChatFormProps) {
  const [message, setMessage] = useState("");
  const { mutate: sendMessage, isLoading: sending } = useMutation("POST", `/api/chat/${chatId}/message`);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!message) return;
    sendMessage({ message }).then((message) => {
      setMessage("");
      addChatMessage(chatId, message);
    });
  };

  return (
    <form
      className="row v-center px-md bg-white m-lg py-sm rounded-md shadow-300"
      style={{ position: "absolute", bottom: 0, left: 0, right: 0 }}
      onSubmit={handleSubmit}
    >
      <input
        className="flex"
        placeholder="Scrivi un messaggio..."
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button className="btn primary" type="submit" disabled={sending}>
        Invia
      </button>
    </form>
  );
}
