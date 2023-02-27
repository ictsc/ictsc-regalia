import ICTSCNavBar from "../components/Navbar";
import MarkdownPreview from "../components/MarkdownPreview";
import ICTSCCard from "../components/Card";
import Head from "next/head";
import { site } from "../components/_const";

const Home = () => {
  // TODO(k-shir0): あとで消す
  const content = `
  # 問題タイトル
  
  ## サブタイトル
  
  ### サブサブタイトル
  
  これは本文です。
  
  \`\`\`mermaid
  graph TB
      A[Hard edge] -->|Link text| B(Round edge)
      B --> C{Decision}
      C -->|One| D[Result one]
      C -->|Two| E[Result two]
  \`\`\`
  `;

  return (
    <>
      <Head>
        <title>ルール - {site}</title>
      </Head>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>ルール</h1>
      <div className={"container-ictsc"}>
        <ICTSCCard>
          <MarkdownPreview content={content} />
        </ICTSCCard>
      </div>
    </>
  );
};

export default Home;
