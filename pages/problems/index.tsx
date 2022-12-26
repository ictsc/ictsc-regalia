import ICTSCNavBar from "../../components/Navbar";
import ProblemCard from "../../components/ProblemCard";
import {useProblems} from "../../hooks/problem";
import {LoadingPage} from "../../components/LoadingPage";


const Problems = () => {
  const {problems, loading} = useProblems()

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
        <h1 className={'title-ictsc text-center py-12'}>問題一覧</h1>
        <ul className={'grid grid-cols-4 gap-8 container-ictsc'}>
          {problems && problems.map((problem) => (
              <li key={problem.id}>
                <ProblemCard problem={problem}/>
              </li>
          ))}
        </ul>
      </>
  )
}

export default Problems