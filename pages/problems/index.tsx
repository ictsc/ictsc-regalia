import ICTSCNavBar from "../../components/navbar";
import ProblemCard from "../../components/ProblemCard";

const Problems = () => {
  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'title-ictsc text-center py-12'}>問題一覧</h1>
        <ul className={'grid grid-cols-4 gap-8 container-ictsc'}>
          <li>
            <ProblemCard/>
          </li>
          <li>
            <ProblemCard/>
          </li>
          <li>
            <ProblemCard/>
          </li>
          <li>
            <ProblemCard/>
          </li>
          <li>
            <ProblemCard/>
          </li>
        </ul>
      </>
  )
}

export default Problems