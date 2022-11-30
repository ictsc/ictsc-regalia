import ICTSCNavBar from "../../components/Navbar";
import MarkdownPreview from "../../components/MarkdownPreview";
import 'zenn-content-css';
import ICTSCCard from "../../components/Card";

const Problem = () => {
  // TODO(k-shir0): problemId を取得するコード
  // const router = useRouter();
  // const {problemId} = router.query;

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
        <div className={'container-ictsc'}>
          <div className={'flex flex-row items-end py-12'}>
            <h1 className={'title-ictsc pr-4'}>問題タイトル</h1>
            満点
            100 pt 採点基準
            100 pt
          </div>
          <ICTSCCard>
            <MarkdownPreview content={content}/>
          </ICTSCCard>
        </div>
      </>
  )
}

export default Problem
