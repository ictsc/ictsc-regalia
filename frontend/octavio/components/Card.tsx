import React from "react";

interface Props {
  className?: string;
  children: React.ReactNode;
}

function ICTSCCard({ className, children }: Props) {
  return (
    <div
      className={`border px-8 pt-12 pb-8  rounded-md shadow-sm ${className}`}
    >
      {children}
    </div>
  );
}

ICTSCCard.defaultProps = {
  className: "",
};

export default ICTSCCard;
