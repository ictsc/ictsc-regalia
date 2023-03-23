import Image from "next/image";

type Props = {
  className?: string;
  message: string;
  subMessage?: string;
};
export const ICTSCSuccessAlert = ({
  className,
  message,
  subMessage,
}: Props) => {
  return (
    <div
      className={`alert alert-success shadow-lg max-w-xs min-w-[312ppx] ${className}`}
    >
      <div>
        <Image
          src={"/assets/svg/check-circle.svg"}
          width={24}
          height={24}
          alt={"success"}
        />
        <div>
          <h3>{message}</h3>
          {subMessage && <span className="text-xs">{subMessage}</span>}
        </div>
      </div>
    </div>
  );
};

export const ICTSCErrorAlert = ({ className, message, subMessage }: Props) => {
  return (
    <div
      className={`alert alert-error shadow-lg max-w-xs min-w-[312ppx] ${className}`}
    >
      <div>
        <Image
          src={"/assets/svg/x-circle.svg"}
          width={24}
          height={24}
          alt={"success"}
        />
        <div>
          <h3>{message}</h3>
          {subMessage && <span className="text-xs">{subMessage}</span>}
        </div>
      </div>
    </div>
  );
};
