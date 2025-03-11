import { useReducer } from "react";
import { MaterialSymbol } from "./material-symbol";
import { clsx } from "clsx";

type NotificationState = boolean;
type NotificationAction = "show" | "hide";

function reduceNotification(
  _state: NotificationState,
  action: NotificationAction,
): NotificationState {
  switch (action) {
    case "hide":
      return false;
    default:
      return _state;
  }
}

export function NotificationBanner(props: { readonly message: string }) {
  const [isVisible, dispatch] = useReducer(reduceNotification, true);

  if (!isVisible) return null;

  return (
    <div
      className={clsx(
        "flex flex-row items-center justify-between rounded-12 bg-surface-2 p-8",
      )}
    >
      <div className="flex items-center gap-4">
        <MaterialSymbol icon="error" fill size={24} />
        <span className="text-16 font-bold text-text">{props.message}</span>
      </div>
      <button
        onClick={() => {
          dispatch("hide");
        }}
        className="flex justify-center rounded-full hover:opacity-80"
      >
        <MaterialSymbol icon="close" size={24} />
      </button>
    </div>
  );
}
