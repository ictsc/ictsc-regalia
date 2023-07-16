import React, { useState } from "react";

import ICTSCCard from "@/components/Card";
import HiddenInput from "@/components/HiddenInput";
import LoadingPage from "@/components/LoadingPage";
import Eye from "@/components/icons/Eye";
import EyeSlash from "@/components/icons/EyeSlash";
import useAuth from "@/hooks/auth";
import useRanking from "@/hooks/ranking";
import CommonLayout from "@/layouts/CommonLayout";

function TeamInfo() {
  const { user, isLoading } = useAuth();
  const { ranking, loading: isRankingLoading } = useRanking();
  const [isSSHHidden, setIsSSHHidden] = useState(false);
  const [isPasswordHidden, setIsPasswordHidden] = useState(true);

  if (isLoading || isRankingLoading) {
    return (
      <CommonLayout title="チーム情報">
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
    <CommonLayout title="チーム情報">
      <div className="container-ictsc">
        <ICTSCCard>
          <p className="font-extrabold text-4xl">
            {user?.user_group.name}@{user?.user_group.organization}
          </p>
          {rank && (
            <ul className="grid grid-cols-3 gap-8 pt-4">
              <div className="border rounded-md p-4">
                <p className="text-sm">ランキング</p>
                <p className="pt-1">
                  <span className="font-extrabold text-3xl">{rank.rank}</span>
                  <span className="text-lg">/{ranking?.length} teams</span>
                </p>
              </div>
              <div className="border rounded-md p-4">
                <p className="text-sm">得点</p>
                <p className="pt-1">
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
      </div>
    </CommonLayout>
  );
}

export default TeamInfo;
