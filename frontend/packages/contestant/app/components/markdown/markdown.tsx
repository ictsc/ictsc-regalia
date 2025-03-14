import { use, useMemo } from "react";
import { clsx } from "clsx";
import styles from "./markdown.module.css";
import { markdownToHtml } from "./markdown-utils";

export function Typography(props: {
  className?: string;
  children?: React.ReactNode;
}) {
  return (
    <div className={clsx(styles.content, props.className)}>
      {props.children}
    </div>
  );
}

export function Markdown({ children }: { children?: string }) {
  const htmlPromise = useMemo(() => markdownToHtml(children ?? ""), [children]);
  const html = use(htmlPromise);
  return (
    <div
      className={clsx(styles.content)}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
