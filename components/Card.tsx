import MarkdownPreview from "./MarkdownPreview";
import React from "react";

interface Props {
  children: React.ReactNode;
}

const ICTSCCard = ({children}: Props) => {
  return (
      <div className={'border px-8 pt-12 pb-8  rounded-md shadow-sm'}>
        {children}
      </div>
  )
}

export default ICTSCCard