import React from "react";

import clsx from "clsx";

type Props = {
  className?: string;
  value: string;
  isHidden: boolean;
  onClick: (e: React.MouseEvent<HTMLInputElement>) => void;
};

function HiddenInput({ className, value, isHidden, onClick }: Props) {
  return (
    <input
      type={isHidden ? "password" : "readonly"}
      className={clsx(
        "input input-bordered max-w-[440px] grow select-none",
        className
      )}
      value={value}
      onClick={onClick}
    />
  );
}

HiddenInput.defaultProps = {
  className: "",
};

export default HiddenInput;
