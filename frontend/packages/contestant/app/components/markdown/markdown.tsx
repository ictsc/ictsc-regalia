import { clsx } from "clsx";
import ReactMarkdown, { type Components } from "react-markdown";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import rehypeKatex from "rehype-katex";
import katexCSS from "katex/dist/katex.min.css?url";
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
const markdownComponents: Components = {
  span: (props) => {
    const { node: _node, ...elProps } = props;
    const el = <span {...elProps} />;

    if (elProps.className?.includes("katex")) {
      return (
        <>
          <link rel="stylesheet" precedence="low" href={katexCSS} />
          {el}
        </>
      );
    }
    return el;
  },
};

export function Markdown({ children }: { children?: string }) {
  return (
    <ReactMarkdown
      remarkPlugins={remarkPlugins}
      rehypePlugins={rehypePlugins}
      components={markdownComponents}
    >
      {children}
    </ReactMarkdown>
  );
}
