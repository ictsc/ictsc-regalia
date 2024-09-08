import React from "react";

import clsx from "clsx";

function Card({
  children,
  className = "",
  hidden = false,
}: {
  children: React.ReactNode;
  className?: string;
  hidden?: boolean;
}) {
  return (
    <div
      className={clsx("border md:border-t rounded-md shadow-sm", className)}
      hidden={hidden}
    >
      {children}
    </div>
  );
}

export default Card;
