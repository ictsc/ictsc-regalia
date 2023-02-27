import React, { useState } from "react";

import ICTSCNavBar from "../components/Navbar";
import ICTSCCard from "../components/Card";
import LoadingPage from "../components/LoadingPage";
import { useAuth } from "../hooks/auth";
import { useRanking } from "../hooks/ranking";

const TeamInfo = () => {
  const { user, isLoading } = useAuth();
  const { ranking, loading: isRankingLoading } = useRanking();
  const [isHidden, setIsHidden] = useState(true);

  if (isLoading || isRankingLoading) {
    return (
      <>
        <ICTSCNavBar />
        <LoadingPage />
      </>
    );
  }

  const rank = ranking?.find((r) => r.user_group_id === user?.user_group.id);
  const bastionUser = user?.user_group.bastion?.bastion_user ?? "";
  const bastionPort = user?.user_group.bastion?.bastion_port ?? "";
  const bastionIp = user?.user_group.bastion?.bastion_host ?? "";
  const bastionPassword = user?.user_group.bastion?.bastion_password ?? "";
  const ssh = `ssh ${bastionUser}@${bastionIp} -p ${bastionPort}`;

  return (
    <>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>チーム情報</h1>
      <div className={"container-ictsc"}>
        <ICTSCCard>
          <p className={"font-extrabold text-4xl"}>
            {user?.user_group.name}@{user?.user_group.organization}
          </p>
          {rank && (
            <>
              <ul className={"grid grid-cols-3 gap-8 pt-4"}>
                <div className={"border rounded-md p-4"}>
                  <p className={"text-sm"}>ランキング</p>
                  <p className={"pt-1"}>
                    <span className={"font-extrabold text-3xl"}>
                      {rank.rank}
                    </span>
                    <span className={"text-lg"}>/{ranking?.length} teams</span>
                  </p>
                </div>
                <div className={"border rounded-md p-4"}>
                  <p className={"text-sm"}>得点</p>
                  <p className={"pt-1"}>
                    <span className={"font-extrabold text-3xl"}>
                      {rank.point}
                    </span>
                    <span className={"text-lg pl-2"}>pt</span>
                  </p>
                </div>
                <div className={"border rounded-md p-4"}>
                  <p className={"text-sm"}>メンバー</p>
                  <ul className={"pt-1 "}>
                    {rank.user_group.members?.map((user) => (
                      <div key={user.id} className={"badge text-black"}>
                        {user.display_name}
                      </div>
                    ))}
                  </ul>
                </div>
              </ul>
            </>
          )}
          <label className={"label pt-12"}>
            <span className={"label-text"}>接続情報</span>
          </label>
          <div
            className={
              "flex justify-between input input-bordered items-center max-w-[336px]"
            }
            onClick={() => {
              setIsHidden(!isHidden);
            }}
          >
            <p className={`${!isHidden && "invisible"}`}>{ssh}</p>
            <button
              className={`link link-hover`}
              onClick={() => {
                navigator.clipboard.writeText(ssh);
              }}
            >
              Copy
            </button>
          </div>
          <label className={"label"}>
            <span className={"label-text"}>パスワード</span>
          </label>
          <div
            className={
              "flex justify-between input input-bordered items-center max-w-[336px]"
            }
          >
            ********
            <button
              className={`link link-hover`}
              onClick={() => {
                navigator.clipboard.writeText(bastionPassword);
              }}
            >
              Copy
            </button>
          </div>
        </ICTSCCard>
      </div>
    </>
  );
};

export default TeamInfo;
