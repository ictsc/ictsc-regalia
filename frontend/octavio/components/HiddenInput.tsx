import React from "react";

type Props = {
  // eslint-disable-next-line react/require-default-props
  className?: string;
  value: string;
  isHidden: boolean;
  onClick: (e: React.MouseEvent<HTMLInputElement>) => void;
};

function HiddenInput({ className, value, isHidden, onClick }: Props) {
  return (
    <input
      type={isHidden ? "password" : "readonly"}
      className={`input input-bordered max-w-[440px] grow select-none ${className}`}
      value={value}
      onClick={onClick}
    />
  );
}

export default HiddenInput;
