import { forwardRef } from "react";
import { Message } from "../context/chat";

type MessageProps = {
  message: Message;
  userId: string;
};

const getContainerClass = (primary: boolean) => {
  if (primary) {
    return "row v-center reverse message primary spacing-1";
  } else {
    return "row v-center message secondary spacing-1";
  }
};

export const MessageComponent = forwardRef<HTMLDivElement, MessageProps>(function (props: MessageProps, ref) {
  const { message, userId } = props;

  const containerClass = getContainerClass(message.sender_id === userId);

  const bubbleClass =
    "rounded-md px-md py-sm " + (message.sender_id === userId ? "c-white bg-primary-950" : "c-primary-950 bg-white");

  return (
    <div ref={ref} className={containerClass} key={message.id}>
      <img className="rounded" src={`api/user/${message.sender_id}/avatar`} alt="avatar" width={50} height={50} />
      <div className={bubbleClass}>
        <p className="text bold">{message.sender_username}</p>
        <p>{message.message}</p>
      </div>
    </div>
  );
});
