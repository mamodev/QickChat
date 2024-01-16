import { z } from "zod";
import { usePaginatedQuery } from "../hooks/usePaginatedQuery";
import { UserListItem } from "./ui/UserListItem";
import { useMutation } from "../hooks/useMutation";
import { invalidateQuery } from "../context/queryContext";
import InfiniteScroll from "react-infinite-scroll-component";
import { ListItem } from "./ui/ListItem";

const requestSchema = z.object({
  id: z.string(),
  sender_id: z.string(),
  sender_name: z.string(),
  sender_email: z.string(),
  sender_profile_picture: z.string(),
});

export function Request() {
  const { data, hasMore, fetchMore } = usePaginatedQuery("/api/friend-request", requestSchema);

  const { mutate } = useMutation("POST", "/api/friend-request/:id/respond", {
    onSuccess: () => invalidateQuery("/api/friend-request"),
  });

  const requests = data ?? [];

  const handleResponse = (id: string, response: boolean) => {
    mutate({ accepted: response }, { id });
  };

  return (
    <div
      id="requets-list"
      style={{
        overflow: "auto",
        paddingBottom: "5rem",
      }}
    >
      <InfiniteScroll
        dataLength={data.length}
        next={fetchMore}
        className="stack flex spacing-1 px-md"
        hasMore={hasMore}
        loader={<h4>Loading...</h4>}
        scrollableTarget="requets-list"
      >
        {requests.map((request) => (
          <ListItem key={request.id}>
            <UserListItem
              username={request.sender_name}
              email={request.sender_email}
              profile_picture={request.sender_profile_picture}
              action={
                <div>
                  <button onClick={() => handleResponse(request.id, true)}>Accept</button>
                  <button onClick={() => handleResponse(request.id, false)}>Decline</button>
                </div>
              }
            />
          </ListItem>
        ))}
      </InfiniteScroll>
    </div>
  );
}
