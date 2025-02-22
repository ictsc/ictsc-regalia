import { clsx } from "clsx";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import styles from "./markdown.module.css";

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

const remarkPlugins = [remarkGfm];

export function Markdown({ children }: { children?: string }) {
  return (
    <ReactMarkdown remarkPlugins={remarkPlugins}>{children}</ReactMarkdown>
  );
}
