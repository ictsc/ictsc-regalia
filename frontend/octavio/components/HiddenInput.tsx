import React from "react";

type Props = {
  value: string;
  isHidden: boolean;
  onClick: (e: React.MouseEvent<HTMLInputElement>) => void;
};

function HiddenInput({ value, isHidden, onClick }: Props) {
  return (
    <input
      type={isHidden ? "password" : "readonly"}
      className="input input-bordered max-w-[440px] grow select-none"
      value={value}
      onClick={onClick}
    />
  );
}

export default HiddenInput;
