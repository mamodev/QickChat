import { chatSchema, setSelectedChat, useSelectedChat } from "../context/chatContext";
import { useQuery } from "../hooks/useQuery";
import { Profile } from "./Profile";
import { ListItem } from "./ui/ListItem";

export default function ChatList() {
  const selected = useSelectedChat();

  const { data: chats, error } = useQuery("/api/chat", chatSchema);

  if (error) {
    return <p>{error.message}</p>;
  }

  return (
    <div className="stack v-center" style={{ justifyContent: "space-between" }}>
      <div>
        <div className="row v-center px-sm">
          <img src="/logo.png" alt="logo" width={80} height={80} />
          <h1 style={{ paddingRight: 10 }}>QuickChat</h1>
        </div>

        <div className="stack">
          {chats?.map((chat) => (
            <ListItem selected={selected?.id === chat.id} key={chat.id} onClick={() => setSelectedChat(chat)}>
              <div className="row v-center spacing-1 p-sm">
                <img
                  src={`api/chat/${chat.id}/avatar`}
                  alt={chat.name}
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <p>{chat.name}</p>
              </div>
            </ListItem>
          ))}
        </div>
      </div>

      <Profile />
    </div>
  );
}
