import markdownToHtml from "zenn-markdown-html";
import "zenn-content-css";

interface Props {
  className?: string;
  content: string;
}

const MarkdownPreview = ({className, content}: Props) => {
  // TODO(k-shir0): clipboard に未対応
  return (
      <div
          className={`znc ${className}`}
        dangerouslySetInnerHTML={{__html: markdownToHtml(content)}}
      />
  );
};

export default MarkdownPreview;