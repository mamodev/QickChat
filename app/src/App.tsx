import Auth from "./Auth";
import ChatList from "./components/ChatList";
import { ChatPanel } from "./components/ChatPanel";
import { useAuth } from "./context/authContext";
import { addChatMessage, messageSchema } from "./context/chat";
import { safeFetch } from "./hooks/fetch";
import { useOnlineStatus } from "./hooks/useOnline";
import { usePolling } from "./hooks/usePolling";

function App() {
  const authenticated = useAuth();
  const isOnline = useOnlineStatus();

  if (!isOnline) {
    return (
      <div className="app-container v-center h-center">
        <div className="stack v-center spacing-1">
          <p className="title bold">You are offline</p>
          <p className="text">Please check your internet connection</p>
        </div>
      </div>
    );
  }

  if (!authenticated) {
    return <Auth />;
  }

  return <AuthApp />;
}

function AuthApp() {
  usePolling(async () => {
    const messages = await safeFetch(`/api/chats-update`, messageSchema.array());
    if (messages.length === 0) return;

    messages.forEach((message) => {
      addChatMessage(message.chat_id, message);
    });
  }, 1000);

  return (
    <div className="app-container">
      <ChatList />
      <ChatPanel />
    </div>
  );
}

export default App;
