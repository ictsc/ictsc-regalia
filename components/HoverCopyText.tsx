import React, { useState } from "react";

import MarkdownPreview from "./MarkdownPreview";

interface Props {
  text: string;
}

function HoverCopyText({ text }: Props) {
  const [isHover, setIsHover] = useState(false);

  return (
    <td
      onMouseEnter={() => {
        setIsHover(true);
      }}
      onMouseLeave={() => {
        setIsHover(false);
      }}
    >
      <div className="flex items-center">
        <MarkdownPreview content={`\`${text}\``} />
        <button
          type="button"
          className={`link link-hover pl-2 ${!isHover && "invisible"}`}
          onClick={() => {
            navigator.clipboard.writeText(text);
          }}
        >
          Copy
        </button>
      </div>
    </td>
  );
}

export default HoverCopyText;
