"use client";

import { useEffect, useMemo, useRef, useState } from "react";

import "zenn-content-css";
import clsx from "clsx";
import markdownToHtml from "zenn-markdown-html";

interface Props {
  className?: string;
  content?: string;
}

const CopyOutlineIcon = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" className="fill-white" scale="1.2">
      <g data-name="Layer 2">
        <g data-name="copy">
          <rect width="24" height="24" opacity="0"/>
          <path
              d="M18 21h-6a3 3 0 0 1-3-3v-6a3 3 0 0 1 3-3h6a3 3 0 0 1 3 3v6a3 3 0 0 1-3 3zm-6-10a1 1 0 0 0-1 1v6a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1v-6a1 1 0 0 0-1-1z"/>
          <path
              color={'white'}
              d="M9.73 15H5.67A2.68 2.68 0 0 1 3 12.33V5.67A2.68 2.68 0 0 1 5.67 3h6.66A2.68 2.68 0 0 1 15 5.67V9.4h-2V5.67a.67.67 0 0 0-.67-.67H5.67a.67.67 0 0 0-.67.67v6.66a.67.67 0 0 0 .67.67h4.06z"/>
        </g>
      </g>
    </svg>`;

function MarkdownPreview({ className, content }: Props) {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const codeBlockContainers = ref.current?.querySelectorAll(
      "div.code-block-container",
    );
    codeBlockContainers?.forEach((codeBlockContainer) => {
      // eslint-disable-next-line no-param-reassign
      codeBlockContainer.className = "relative";
      const pre = codeBlockContainer.querySelector("pre");
      const code = pre?.querySelector("code");

      const button = document.createElement("button");
      const buttonClassName =
        "btn btn-xs btn-circle btn-ghost fix absolute top-[8px] right-[20px] z-index-10 ";
      button.className = `${buttonClassName} invisible`;

      // CopyOutlineIcon を追加
      const icon = document.createElement("div");
      icon.className = "w-12 h-12 fill-white";
      icon.innerHTML = CopyOutlineIcon;

      button.appendChild(icon);
      codeBlockContainer.appendChild(button);

      button.addEventListener("click", () => {
        const c = code?.innerText;
        if (c) {
          // クリップボードにコピー
          navigator.clipboard.writeText(c);
        }
      });

      codeBlockContainer.addEventListener("mouseover", () => {
        button.className = `${buttonClassName} visible`;
      });
      codeBlockContainer.addEventListener("mouseout", () => {
        button.className = `${buttonClassName} invisible`;
      });
    });
  }, []);

  return useMemo(() => {
    let count = 0;

    return (
      <div
        className={clsx("znc", className)}
        ref={ref}
        /* eslint-disable-next-line react/no-danger */
        dangerouslySetInnerHTML={{
          __html: markdownToHtml(content ?? "", {
            customEmbed: {
              ictscr(str) {
                const lines = str
                  .split("\n")
                  .filter((line) => line.startsWith("- [ ]"));

                const radioButtons = lines.map((line, index) => {
                  const label = line.replace("- [ ]", "").trim();
                  return `
                    <div class="flex items-center mb-4">
                      <input 
                        id="radio-${count}-${index}" 
                        type="radio" 
                        value=${index} 
                        name="radio${count}" 
                        class="w-6 h-4 text-blue-600 bg-gray-100 border-gray-300"
                      >
                      <label for="radio-${count}-${index}" class="ml-2 text-sm font-medium text-gray-900">${label}</label>
                    </div>`;
                });

                count += 1;

                return `
                  <div class="radio-buttons pt-2">
                    <fieldset class="flex flex-col space-y-4">${radioButtons.join(
                      "",
                    )}</fieldset>
                  </div>`;
              },
              ictscc(str) {
                // - [ ] で始まる行を取得
                const lines = str
                  .split("\n")
                  .filter((line) => line.startsWith("- [ ]"));

                // checkbox を追加
                const checkboxes = lines.map((line, index) => {
                  const label = line.replace("- [ ]", "").trim();
                  return `
                        <div class="flex items-center mb-4">
                            <input id="checkbox-${count}-${index}" type="checkbox" value=${index} name="checkbox${count}" class="w-6 h-4 text-blue-600 bg-gray-100 border-gray-30">
                            <label for="checkbox-${count}-${index}" class="ml-2 text-sm font-medium text-gray-900">${label}</label>
                        </div>`;
                });

                count += 1;

                return `
          <div class="checkboxes pt-2">
              <fieldset class="flex flex-col space-y-4">${checkboxes.join(
                "",
              )}</fieldset>
          </div>`;
              },
            },
          }),
        }}
      />
    );
  }, [content, className]);
}

MarkdownPreview.defaultProps = {
  className: "",
  content: "",
};

export default MarkdownPreview;
