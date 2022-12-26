import 'zenn-content-css';
import {useRouter} from "next/router";
import Error from "next/error";

import ICTSCNavBar from "../../components/Navbar";
import MarkdownPreview from "../../components/MarkdownPreview";
import ICTSCCard from "../../components/Card";
import {useProblems} from "../../hooks/problem";
import {LoadingPage} from "../../components/LoadingPage";

const ProblemPage = () => {
  const router = useRouter();

  const {getProblem, loading} = useProblems();

  const {problemId} = router.query;
  const problem = getProblem(problemId as string);

  if (problem === null) {
    return <Error statusCode={404}/>;
  }

  if (loading) {
    return (
        <>
          <ICTSCNavBar/>
          <LoadingPage/>
        </>
    );
  }

  return (
      <>
        <ICTSCNavBar/>
        <div className={'container-ictsc'}>
          <div className={'flex flex-row items-end py-12'}>
            <h1 className={'title-ictsc pr-4'}>{problem?.title}</h1>
            満点
            {problem?.point} pt
            採点基準
            {problem?.solved_criterion} pt
          </div>
          <ICTSCCard>
            <MarkdownPreview content={problem?.body}/>
          </ICTSCCard>
        </div>
      </>
  )
}

export default ProblemPage
