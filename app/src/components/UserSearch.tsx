import { useState } from "react";
import { z } from "zod";
import { ListItem } from "./ui/ListItem";
import { useMutation } from "../hooks/useMutation";
import { usePaginatedQuery } from "../hooks/usePaginatedQuery";
import { UserListItem } from "./ui/UserListItem";
import { invalidateQuery } from "../context/queryContext";
import InfiniteScroll from "react-infinite-scroll-component";
import { setSelectedChat } from "../context/chatContext";
import { useDialog } from "./ui/Dialog";

export function UserSearch() {
  const [name, setName] = useState("");

  return (
    <div className="stack" style={{ overflow: "hidden", paddingTop: 10 }}>
      <div className="f-width row px-sm">
        <input
          className="flex input rounded-sm"
          value={name}
          placeholder="Username..."
          onChange={(e) => setName(e.target.value)}
        />
      </div>
      <UserList username={name} />
    </div>
  );
}

type UserListProps = {
  username: string;
};

const userSchema = z.object({
  id: z.string(),
  username: z.string(),
  email: z.string().email(),
  profile_picture: z.string(),
  friend_request_sent: z.boolean(),
  is_friend: z.boolean(),
  friend_request_id: z.string().nullable(),
});

function UserList({ username }: UserListProps) {
  const { data, hasMore, fetchMore, error } = usePaginatedQuery(`/api/user?user=${username}`, userSchema);

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  const users = data ?? [];

  return (
    <div
      id="user-list"
      style={{
        overflowY: "scroll",
        paddingBottom: "5rem",
        maxHeight: "100%",
      }}
    >
      <InfiniteScroll
        dataLength={data.length}
        next={fetchMore}
        className="stack flex spacing-1 px-md"
        hasMore={hasMore}
        loader={<h4>Loading...</h4>}
        scrollableTarget="user-list"
      >
        {users.map((user) => (
          <ListItem key={user.id}>
            <UserItem user={user} />
          </ListItem>
        ))}
      </InfiniteScroll>
    </div>
  );
}

type UserItemProps = {
  user: z.infer<typeof userSchema>;
};

const options = {
  // onSuccess: () => invalidateQuery("/api/user"),
};

function UserItem({ user }: UserItemProps) {
  const { onClose } = useDialog();
  const { mutate: create, isError: creationError } = useMutation("POST", "/api/friend-request", options);
  const { mutate: remove, isError: deletionError } = useMutation("DELETE", "/api/friend-request/:id", options);
  const { mutate: createChat, isError: acceptError } = useMutation("POST", "/api/chat", options);

  const handleAction = () => {
    if (user.friend_request_id) {
      user.friend_request_sent = false;
      remove(null, { id: user.friend_request_id })
        .then(() => {
          user.friend_request_id = null;
        })
        .catch(() => {
          user.friend_request_sent = true;
        });
    }

    if (!user.is_friend && !user.friend_request_id) {
      user.friend_request_sent = true;
      create({ user_id: user.id })
        .then((res) => {
          if (res.request_id) user.friend_request_id = res.request_id;
          if (res.friended) user.is_friend = true;
        })
        .catch(() => {
          user.friend_request_sent = false;
        });
    }

    if (user.is_friend)
      createChat({ users: [user.id], name: user.username }).then((res) => {
        invalidateQuery(`/api/chat`);
        if (res.id) setSelectedChat({ id: res.id, name: user.username });
        onClose();
      });
  };

  const error = creationError || deletionError || acceptError;

  return (
    <UserListItem
      username={user.username}
      email={user.email}
      profile_picture={user.profile_picture}
      action={
        error ? (
          <div>Error</div>
        ) : (
          <button onClick={handleAction}>
            {user.friend_request_id && "Cancel Request"}
            {!user.friend_request_id && !user.is_friend && "Add Friend"}
            {user.is_friend && "Send"}
          </button>
        )
      }
    />
  );
}
