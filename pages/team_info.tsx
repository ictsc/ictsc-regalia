import ICTSCNavBar from "../components/Navbar";
import ICTSCCard from "../components/Card";
import React from "react";
import { useUserGroups } from "../hooks/userGroups";
import { useAuth } from "../hooks/auth";
import LoadingPage from "../components/LoadingPage";
import users from "./users";

const TeamInfo = () => {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <>
        <ICTSCNavBar />
        <LoadingPage />
      </>
    );
  }

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
          <label className={"label"}>
            <span className={"label-text"}>接続情報</span>
          </label>
          <div
            className={
              "flex justify-between input input-bordered items-center max-w-[336px]"
            }
          >
            {ssh}
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
