import ICTSCNavBar from "../../components/Navbar";
import MarkdownPreview from "../../components/MarkdownPreview";
import 'zenn-content-css';
import ICTSCCard from "../../components/Card";

// @ts-ignore
import metadataParser from 'markdown-yaml-metadata-parser';

const ProblemPage = () => {
  // TODO(k-shir0): problemId を取得するコード
  // const router = useRouter();
  // const {problemId} = router.query;

  // TODO(k-shir0): あとで消す
  const source = `---
id: 1
title: 問題タイトル
point:
  max: 100
  solvedCriterion: 100
---
  
  # 問題タイトル
  
  ## サブタイトル
  
  ### サブサブタイトル
  
  これは本文です。
  
  \`\`\`
  sudo hogehoge
  \`\`\`
  `;

  const result = metadataParser(source);
  const {title, point} = result.metadata;

  return (
      <>
        <ICTSCNavBar/>
        <div className={'container-ictsc'}>
          <div className={'flex flex-row items-end py-12'}>
            <h1 className={'title-ictsc pr-4'}>{title}</h1>
            満点
            {point.max} pt
            採点基準
            {point.solvedCriterion} pt
          </div>
          <ICTSCCard>
            <MarkdownPreview content={result.content}/>
          </ICTSCCard>
        </div>
      </>
  )
}

export default ProblemPage
