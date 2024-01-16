import { createContext, useContext } from "react";

type DialogProps = {
  children: React.ReactNode;
  open: boolean;
  title?: string;
  onClose: () => void;
};

const DialogContext = createContext<DialogProps | null>(null);

export function Dialog(props: DialogProps) {
  const { children, open, onClose, title } = props;
  if (!open) return <></>;

  return (
    <DialogContext.Provider value={{ children, open, onClose, title }}>
      <dialog open={true} className="dialog stack spacing-0">
        <div className="row v-center p-sm shadow-300" style={{ justifyContent: "space-between" }}>
          <p className="title">{title ?? "Title"}</p>
          <button className="btn primary" onClick={onClose}>
            Close
          </button>
        </div>
        {children}
      </dialog>
    </DialogContext.Provider>
  );
}

export function useDialog() {
  const ctx = useContext(DialogContext);
  if (!ctx) throw new Error("useDialog must be used within a DialogProvider");
  return ctx;
}
