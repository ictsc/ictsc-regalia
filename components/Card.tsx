import React from "react";

interface Props {
  className?: string;
  children: React.ReactNode;
}

const ICTSCCard = ({className, children}: Props) => {
  return (
      <div className={`border px-8 pt-12 pb-8  rounded-md shadow-sm ${className}`}>
        {children}
      </div>
  )
}

export default ICTSCCard