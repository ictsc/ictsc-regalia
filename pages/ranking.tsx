import Error from "next/error";

import ICTSCNavBar from "../components/Navbar";
import LoadingPage from "../components/LoadingPage";

import {useRanking} from "../hooks/ranking";

const Ranking = () => {
  const {ranking, loading} = useRanking()


  if (loading) {
    return (
        <>
          <ICTSCNavBar/>
          <LoadingPage/>
        </>
    );
  }

  if (ranking === null) {
    return <Error statusCode={404}/>;
  }


  return (
      <>
        <ICTSCNavBar/>
        <h1 className={'title-ictsc text-center py-12'}>スコアボード</h1>
        <div className={'container-ictsc'}>
          <table className="table border w-full">
            <thead>
            <tr>
              <th>#</th>
              <th>チーム名</th>
              <th>所属</th>
              <th className={'text-right'}>得点</th>
            </tr>
            </thead>
            <tbody>
            {ranking?.map((rank) => (
                <tr key={rank.user_group.id}>
                  <td className={'w-[64px]'}>{rank.rank}</td>
                  <td>{rank.user_group.name}</td>
                  <td>{rank.user_group.organization}</td>
                  <td className={'w-[124px] text-right'}>{rank.point}pt</td>
                </tr>
            ))}
            </tbody>
          </table>
        </div>
      </>
  )
}

export default Ranking