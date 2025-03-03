import { Link} from "@tanstack/react-router";

export function NotificationList({ notifications }) {
    return (
      <div className="flex flex-col items-center h-full justify-center">
        {notifications.map((notification) => (
          <Link
              to="/problems"
              title="アナウンス"
              className="bg-surface-1 border-text rounded-8 w-full pl-40 py-4 mb-16 text-16 font-bold hover:bg-surface-2 transition-colors"
          >
              <div>
                  {notification.text}
              </div>
          </Link>
        ))}
      </div>
    );
  }
