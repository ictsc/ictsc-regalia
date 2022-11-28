import ICTSCNavBar from "../components/navbar";
import ProblemCard from "../components/ProblemCard";

const Home = () => {
  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'text-2xl font-bold py-12 text-center'}>問題一覧</h1>
        <ul className={'grid grid-cols-4 gap-8 ictsc-container'}>
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


export default Home