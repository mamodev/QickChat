import { forwardRef } from "react";

type ListItemProps = {
  children: React.ReactNode;
  onClick?: () => void;
  selected?: boolean;
};

export const ListItem = forwardRef<HTMLDivElement, ListItemProps>(({ selected, children, onClick }, ref) => {
  let className = "item";

  if (onClick) {
    className += " clickable";
  }

  if (selected) {
    className += " selected";
  }

  return (
    <div className={className} onClick={onClick} ref={ref}>
      {children}
    </div>
  );
});
