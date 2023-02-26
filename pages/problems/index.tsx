import ICTSCNavBar from "../../components/Navbar";
import ProblemCard from "../../components/ProblemCard";
import LoadingPage from "../../components/LoadingPage";
import { useProblems } from "../../hooks/problem";

const Problems = () => {
  const { problems, isLoading } = useProblems();

  if (isLoading) {
    return (
      <>
        <ICTSCNavBar />
        <LoadingPage />
      </>
    );
  }

  return (
    <>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>問題一覧</h1>
      <ul className={"grid grid-cols-4 gap-8 container-ictsc"}>
        {problems &&
          problems.map((problem, index) => (
            <li key={problem.id}>
              <ProblemCard index={index + 1} problem={problem} />
            </li>
          ))}
      </ul>
    </>
  );
};

export default Problems;
