"use client";

import React, { useState } from "react";

import ICTSCCard from "@/components/card";
import HiddenInput from "@/components/hidden-input";
import Eye from "@/components/icons/eye";
import EyeSlash from "@/components/icons/eye-slash";
import LoadingPage from "@/components/loading-page";
import useAuth from "@/hooks/auth";
import useRanking from "@/hooks/ranking";

function TeamInfo() {
  const { user, isLoading } = useAuth();
  const { ranking, loading: isRankingLoading } = useRanking();
  const [isSSHHidden, setIsSSHHidden] = useState(false);
  const [isPasswordHidden, setIsPasswordHidden] = useState(true);

  if (isLoading || isRankingLoading) {
    return <LoadingPage />;
  }

  const rank = ranking?.find((r) => r.user_group_id === user?.user_group.id);
  const bastionUser = user?.user_group.bastion?.bastion_user ?? "";
  const bastionPort = user?.user_group.bastion?.bastion_port ?? "";
  const bastionIp = user?.user_group.bastion?.bastion_host ?? "";
  const bastionPassword = user?.user_group.bastion?.bastion_password ?? "";
  const ssh = `ssh ${bastionUser}@${bastionIp} -p ${bastionPort}`;

  return (
    <ICTSCCard>
      <p className="group-name-and-organization font-extrabold text-4xl">
        {user?.user_group.name}@{user?.user_group.organization}
      </p>
      {rank && (
        <ul className="grid grid-cols-3 gap-8 pt-4">
          <div className="border rounded-md p-4">
            <p className="text-sm">ランキング</p>
            <p className="ranking pt-1">
              <span className="font-extrabold text-3xl">{rank.rank}</span>
              <span className="text-lg">/{ranking?.length} teams</span>
            </p>
          </div>
          <div className="border rounded-md p-4">
            <p className="text-sm">得点</p>
            <p className="score pt-1">
              <span className="font-extrabold text-3xl">{rank.point}</span>
              <span className="text-lg pl-2">pt</span>
            </p>
          </div>
        </ul>
      )}
      <div className="label pt-12">
        <span className="label-text">踏み台サーバへの接続情報</span>
      </div>
      <div className="flex flex-row">
        <HiddenInput
          className="ssh-info"
          value={ssh}
          isHidden={isSSHHidden}
          onClick={(e) => {
            // @ts-ignore
            e.target.select();
          }}
        />
        <button
          type="button"
          className="ml-[-36px]"
          onClick={() => {
            setIsSSHHidden(!isSSHHidden);
          }}
        >
          {isSSHHidden ? <Eye /> : <EyeSlash />}
        </button>
        <button
          type="button"
          className="link link-hover pl-6"
          onClick={(e) => {
            e.stopPropagation();
            navigator.clipboard.writeText(ssh);
          }}
        >
          Copy
        </button>
      </div>
      <div className="label">
        <span className="label-text">パスワード</span>
      </div>
      <div className="flex flex-row">
        <HiddenInput
          className="password-info"
          value={bastionPassword}
          isHidden={isPasswordHidden}
          onClick={(e) => {
            // @ts-ignore
            e.target.select();
          }}
        />
        <button
          type="button"
          className="ml-[-36px]"
          onClick={() => {
            setIsPasswordHidden(!isPasswordHidden);
          }}
        >
          {isPasswordHidden ? <Eye /> : <EyeSlash />}
        </button>
        <button
          type="button"
          className="link link-hover pl-6"
          onClick={(e) => {
            e.stopPropagation();
            navigator.clipboard.writeText(bastionPassword);
          }}
        >
          Copy
        </button>
      </div>
    </ICTSCCard>
  );
}

export default TeamInfo;
