import { useEffect, useMemo, useState } from "react";
import Error from "next/error";

import { useForm, Controller } from "react-hook-form";

import BaseLayout from "@/layouts/BaseLayout";
import ICTSCCard from "@/components/Card";
import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/Alerts";
import LoadingPage from "../components/LoadingPage";
import { useApi } from "@/hooks/api";
import { useAuth } from "@/hooks/auth";

type Inputs = {
  display_name: string;
  self_introduction: string;
  twitter_id: string;
  github_id: string;
  facebook_id: string;
};

const Profile = () => {
  const { apiClient } = useApi();
  const { user, isLoading, mutate } = useAuth();

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null);

  const getCurrentValue = () => ({
    display_name:
      (user?.display_name ?? "") == "" ? user?.name : user?.display_name,
    self_introduction: user?.user_profile?.self_introduction ?? "",
    github_id: user?.user_profile?.github_id ?? "",
    twitter_id: user?.user_profile?.twitter_id ?? "",
    facebook_id: user?.user_profile?.facebook_id ?? "",
  });

  const {
    register,
    control,
    reset,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>({
    defaultValues: useMemo(
      () => ({
        ...getCurrentValue(),
        // eslint-disable-next-line
      }),
      [user]
    ),
  });

  useEffect(() => {
    // ユーザーをフェッチしたときに反映する
    if (user != null) {
      reset({
        ...getCurrentValue(),
      });
    }
    // eslint-disable-next-line
  }, [user]);

  const onSubmit = async (data: Inputs) => {
    const response = await apiClient.put(`users/${user?.id}`, {
      json: data,
    });

    setStatus(response.status);

    if (response.ok) {
      await mutate();
    }
  };

  if (isLoading) {
    return (
      <BaseLayout title={"プロフィール"}>
        <LoadingPage />
      </BaseLayout>
    );
  }

  if (user === null) {
    return <Error statusCode={404} />;
  }

  return (
    <BaseLayout title={"プロフィール"}>
      <div className={"container-ictsc pt-8"}>
        <ICTSCCard>
          <h1 className={"title-ictsc"}>プロフィール</h1>
          <div className={"divider"} />
          {status === 202 && (
            <ICTSCSuccessAlert
              message={"プロフィールを更新しました"}
              className={"mb-8"}
            />
          )}
          {status != null && status !== 202 && (
            <ICTSCErrorAlert
              message={"プロフィールの更新に失敗しました"}
              className={"mb-8"}
            />
          )}
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className={"form-control"}>
              <label className={"label"}>
                <span className={"label-text"}>表示名*</span>
              </label>
              <input
                {...register("display_name", { required: true })}
                type="text"
                placeholder="名前"
                className={"input input-bordered"}
              />
              {errors.display_name && (
                <p className={"label-text-alt text-error mt-1"}>
                  表示名は必須です
                </p>
              )}
            </div>
            <div className={"form-control pt-4"}>
              <label className={"label"}>
                <span className={"label-text"}>所属チーム</span>
              </label>
              <div className={"pl-1"}>{user?.user_group.name}</div>
            </div>
            <div className={"form-control pt-6"}>
              <label className={"label"}>
                <span className={"label-text"}>自己紹介</span>
              </label>
              <Controller
                name={"self_introduction"}
                control={control}
                render={({ field }) => (
                  <textarea
                    {...field}
                    className={"textarea h-24 textarea-bordered"}
                    placeholder="自己紹介"
                  />
                )}
              />
            </div>
            <div className={"form-control pt-4"}>
              <label className={"label"}>
                <span className={"label-text"}>GitHub ID</span>
              </label>
              <input
                {...register("github_id")}
                type="text"
                placeholder="ユーザー名のみを入力"
                className={"input input-bordered"}
              />
            </div>
            <div className={"form-control pt-4"}>
              <label className={"label"}>
                <span className={"label-text"}>Twitter ID</span>
              </label>
              <input
                {...register("twitter_id")}
                type="text"
                placeholder="@マークなしで入力"
                className={"input input-bordered"}
              />
            </div>
            <div className={"form-control pt-4"}>
              <label className={"label"}>
                <span className={"label-text"}>Facebook ID</span>
              </label>
              <input
                {...register("facebook_id")}
                type="text"
                placeholder="ユーザー名のみを入力"
                className={"input input-bordered"}
              />
            </div>
            <input
              type={"submit"}
              value={"更新"}
              className={"btn btn-primary mt-4"}
            />
          </form>
        </ICTSCCard>
      </div>
    </BaseLayout>
  );
};

export default Profile;
