import ICTSCNavBar from "../../components/Navbar";
import ProblemCard from "../../components/ProblemCard";
import LoadingPage from "../../components/LoadingPage";
import { useProblems } from "../../hooks/problem";
import Head from "next/head";
import { rule, shortRule, site } from "../../components/_const";
import ICTSCCard from "../../components/Card";
import MarkdownPreview from "../../components/MarkdownPreview";

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
      <Head>
        <title>問題一覧 - {site}</title>
      </Head>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>問題一覧</h1>
      {shortRule != "" && (
        <div className={"container-ictsc"}>
          <ICTSCCard>
            <MarkdownPreview content={shortRule?.replace(/\\n/g, "\n") ?? ""} />
          </ICTSCCard>
        </div>
      )}
      <ul
        className={
          "grid grid-cols-2 md:grid-cols-4 lg:grid-cols-4 gap-8 container-ictsc"
        }
      >
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
