import Image from "next/image";

import MarkdownPreview from "@/components/MarkdownPreview";
import { Notice } from "@/types/Notice";

interface Props {
  notice: Notice;
  onDismiss?: () => void;
}

function NotificationCard({ notice, onDismiss }: Props) {
  return (
    <div className="container-ictsc">
      <div className="bg-gray-200 p-4 rounded rounded-lg  shadow-lg grow">
        <div>
          <div className="flex flex-col">
            <div className="flex flex-row justify-between justify-items-center">
              <div className="flex flex-row">
                <span className="font-bold">{notice.title}</span>
              </div>
              {onDismiss && (
                <button type="button" onClick={onDismiss}>
                  <Image
                    src="assets/svg/x-mark.svg"
                    width={24}
                    height={24}
                    alt="dismiss"
                  />
                </button>
              )}
            </div>
            <MarkdownPreview content={notice.body ?? ""} />
          </div>
        </div>
      </div>
    </div>
  );
}

NotificationCard.defaultProps = {
  onDismiss: undefined,
};

export default NotificationCard;
