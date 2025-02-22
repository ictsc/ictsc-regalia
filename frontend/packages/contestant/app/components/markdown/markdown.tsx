import { clsx } from "clsx";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import rehypeKatex from "rehype-katex";
import "katex/dist/katex.min.css";
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

const remarkPlugins = [remarkGfm, remarkMath];
const rehypePlugins = [rehypeKatex];

export function Markdown({ children }: { children?: string }) {
  return (
    <ReactMarkdown remarkPlugins={remarkPlugins} rehypePlugins={rehypePlugins}>
      {children}
    </ReactMarkdown>
  );
}
