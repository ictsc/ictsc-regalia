import ICTSCNavBar from "../../components/Navbar";
import ProblemCard from "../../components/ProblemCard";
import {GetServerSideProps} from "next";
import {Problem} from "../../types/Problem";

interface Props {
  problems: Problem[];
}

export const getServerSideProps: GetServerSideProps = async () => {
  const problemA = {id: "abc", content: "test"};
  const problemB = {id: "def", content: "test"};
  return {
    props: {
      problems: [problemA, problemB]
    }
  }
}

const Problems = ({problems}: Props) => {
  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'title-ictsc text-center py-12'}>問題一覧</h1>
        <ul className={'grid grid-cols-4 gap-8 container-ictsc'}>
          {problems.map((problem) => (
              <li key={problem.id}>
                <ProblemCard/>
              </li>
          ))}
        </ul>
      </>
  )
}

export default Problems