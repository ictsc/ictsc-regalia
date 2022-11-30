import ICTSCNavBar from "../components/Navbar";
import MarkdownPreview from "../components/MarkdownPreview";
import ICTSCCard from "../components/Card";

const Home = () => {
  // TODO(k-shir0): あとで消す
  const content = `
  # 問題タイトル
  
  ## サブタイトル
  
  ### サブサブタイトル
  
  これは本文です。
  
  \`\`\`
  sudo hogehoge
  \`\`\`
  `;

  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'title-ictsc text-center py-12'}>ルール</h1>
        <div className={'container-ictsc'}>
          <ICTSCCard>
            <MarkdownPreview content={content}/>
          </ICTSCCard>
        </div>
      </>
  )
}


export default Home
