import Image from "next/image";

import clsx from "clsx";

type Props = {
  className?: string;
  message: string;
  subMessage?: string;
};

export function ICTSCSuccessAlert({ className, message, subMessage }: Props) {
  return (
    <div
      className={clsx(
        "alert alert-success shadow-lg max-w-xs min-w-[312ppx]",
        className,
      )}
    >
      <div className="flex flex-row">
        <Image
          src="/assets/svg/check-circle.svg"
          width={24}
          height={24}
          alt="success"
        />
        <div className="pl-2">
          <h3>{message}</h3>
          {subMessage && <span className="text-xs">{subMessage}</span>}
        </div>
      </div>
    </div>
  );
}

export function ICTSCErrorAlert({ className, message, subMessage }: Props) {
  return (
    <div
      className={clsx(
        "alert alert-error shadow-lg max-w-xs min-w-[312ppx]",
        className,
      )}
    >
      <div className="flex flex-row">
        <Image
          src="/assets/svg/x-circle.svg"
          width={24}
          height={24}
          alt="success"
        />
        <div className="pl-2">
          <h3>{message}</h3>
          {subMessage && <span className="text-xs">{subMessage}</span>}
        </div>
      </div>
    </div>
  );
}

ICTSCSuccessAlert.defaultProps = {
  className: "",
  subMessage: "",
};

ICTSCErrorAlert.defaultProps = {
  className: "",
  subMessage: "",
};
