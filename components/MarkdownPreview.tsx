import markdownToHtml from "zenn-markdown-html";
import "zenn-content-css";

interface Props {
  content: string;
}

const MarkdownPreview = ({content}: Props) => {
  // TODO(k-shir0): clipboard に未対応
  return (
      <div
          className={"znc"}
        dangerouslySetInnerHTML={{__html: markdownToHtml(content)}}
      />
  );
};

export default MarkdownPreview;