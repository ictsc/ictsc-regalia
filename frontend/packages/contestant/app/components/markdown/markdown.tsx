import { clsx } from "clsx";
import ReactMarkdown, { type Components } from "react-markdown";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import rehypeKatex from "rehype-katex";
import { createHighlighter } from "shiki";
import katexCSS from "katex/dist/katex.min.css?url";
import styles from "./markdown.module.css";
import { use } from "react";

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
const highlighter =  createHighlighter({
  themes: ["dark-plus"],
  langs: ["diff","shellscript","shellsession", "hcl"],
});

export function CodeBlock(props: {
  text: string;
  lang: string;
}){
  const text = props.text;
  const htmlText = use(highlighter).codeToHtml(text, {
    theme: "dark-plus",
    lang: props.lang,
  });
  return (
    <div className="">
    <div dangerouslySetInnerHTML={{ __html: htmlText }} />
    </div>
  );
}

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
  code: (props) => {
    const { node: _node, ...elProps } = props;
    const className = props.className || "";
    const match = /language-(\w+)/.exec(className);
    const lang = match?.[1] ?? "plaintext";
    const text = elProps.children as string | undefined;
    if (typeof text !== "string") {
      return text;
    }
    return (
      <CodeBlock text={text} lang={lang} />
    );
  }
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
