import React, { useState } from "react";

import ICTSCCard from "@/components/Card";
import LoadingPage from "@/components/LoadingPage";
import { useAuth } from "@/hooks/auth";
import { useRanking } from "@/hooks/ranking";
import CommonLayout from "@/layouts/CommonLayout";

const TeamInfo = () => {
  const { user, isLoading } = useAuth();
  const { ranking, loading: isRankingLoading } = useRanking();
  const [isSSHHidden, setIsSSHHidden] = useState(false);
  const [isPasswordHidden, setIsPasswordHidden] = useState(true);

  if (isLoading || isRankingLoading) {
    return (
      <CommonLayout title={"チーム情報"}>
        <LoadingPage />
      </CommonLayout>
    );
  }

  const rank = ranking?.find((r) => r.user_group_id === user?.user_group.id);
  const bastionUser = user?.user_group.bastion?.bastion_user ?? "";
  const bastionPort = user?.user_group.bastion?.bastion_port ?? "";
  const bastionIp = user?.user_group.bastion?.bastion_host ?? "";
  const bastionPassword = user?.user_group.bastion?.bastion_password ?? "";
  const ssh = `ssh ${bastionUser}@${bastionIp} -p ${bastionPort}`;

  return (
    <CommonLayout title={"チーム情報"}>
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
              </ul>
            </>
          )}
          <label className={"label pt-12"}>
            <span className={"label-text"}>踏み台サーバへの接続情報</span>
          </label>
          <div className={"flex flex-row"}>
            <input
              type={isSSHHidden ? "password" : "readonly"}
              className={"input input-bordered max-w-[440px] grow select-none"}
              value={ssh}
              onClick={(e) => {
                // @ts-ignore
                e.target.select();
              }}
            />
            <button
              className={"ml-[-36px]"}
              onClick={() => {
                setIsSSHHidden(!isSSHHidden);
              }}
            >
              {isSSHHidden ? <Eye /> : <EyeSlash />}
            </button>
            <button
              className={`link link-hover pl-6`}
              onClick={(e) => {
                e.stopPropagation();
                navigator.clipboard.writeText(ssh);
              }}
            >
              Copy
            </button>
          </div>
          <label className={"label"}>
            <span className={"label-text"}>パスワード</span>
          </label>
          <div className={"flex flex-row"}>
            <input
              type={isPasswordHidden ? "password" : "readonly"}
              className={"input input-bordered max-w-[440px] grow select-none"}
              value={bastionPassword}
              onClick={(e) => {
                // @ts-ignore
                e.target.select();
              }}
            />
            <button
              className={"ml-[-36px]"}
              onClick={() => {
                setIsPasswordHidden(!isPasswordHidden);
              }}
            >
              {isPasswordHidden ? <Eye /> : <EyeSlash />}
            </button>
            <button
              className={`link link-hover pl-6`}
              onClick={(e) => {
                e.stopPropagation();
                navigator.clipboard.writeText(bastionPassword);
              }}
            >
              Copy
            </button>
          </div>
        </ICTSCCard>
      </div>
    </CommonLayout>
  );
};

const Eye = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    strokeWidth={1.5}
    stroke="currentColor"
    className="w-6 h-6"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
    />
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
    />
  </svg>
);

const EyeSlash = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    strokeWidth={1.5}
    stroke="currentColor"
    className="w-6 h-6"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"
    />
  </svg>
);

export default TeamInfo;
