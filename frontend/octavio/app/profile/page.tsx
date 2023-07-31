"use client";

import { useEffect, useMemo, useState } from "react";

import Error from "next/error";

import { Controller, useForm } from "react-hook-form";

import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/Alerts";
import ICTSCCard from "@/components/Card";
import LoadingPage from "@/components/LoadingPage";
import useAuth from "@/hooks/auth";

type Inputs = {
  display_name: string;
  self_introduction: string;
  github_id: string;
  twitter_id: string;
  facebook_id: string;
};

function Page() {
  const { user, isLoading, mutate, putProfile } = useAuth();

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null);

  const getCurrentValue = () => ({
    display_name:
      (user?.display_name ?? "") === "" ? user?.name : user?.display_name,
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
      }),
      // eslint-disable-next-line
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
    /* c8 ignore next 404 のところでチェックしているため必ず false になる */
    if (user == null) {
      return;
    }
    const response = await putProfile(user.id, data);

    setStatus(response.code);

    if (response.code === 202) {
      await mutate();
    }
  };

  if (isLoading) {
    return <LoadingPage />;
  }

  if (user === null) {
    return <Error statusCode={404} />;
  }

  return (
    <ICTSCCard>
      {status === 202 && (
        <ICTSCSuccessAlert
          message="プロフィールを更新しました"
          className="mb-8"
        />
      )}
      {status != null && status !== 202 && (
        <ICTSCErrorAlert
          message="プロフィールの更新に失敗しました"
          className="mb-8"
        />
      )}
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="form-control">
          <div className="label">
            <span className="label-text">表示名*</span>
          </div>
          <input
            {...register("display_name", { required: true })}
            type="text"
            placeholder="名前"
            className="input input-bordered"
          />
          {errors.display_name && (
            <p className="label-text-alt text-error mt-1">表示名は必須です</p>
          )}
        </div>
        <div className="form-control pt-4">
          <div className="label">
            <span className="label-text">所属チーム</span>
          </div>
          <div className="teams pl-1">{user?.user_group.name}</div>
        </div>
        <div className="form-control pt-6">
          <div className="label">
            <span className="label-text">自己紹介</span>
          </div>
          <Controller
            name="self_introduction"
            control={control}
            render={({ field }) => (
              <textarea
                {...field}
                className="textarea h-24 textarea-bordered"
                placeholder="自己紹介"
              />
            )}
          />
        </div>
        <div className="form-control pt-4">
          <div className="label">
            <span className="label-text">GitHub ID</span>
          </div>
          <input
            {...register("github_id")}
            type="text"
            placeholder="ユーザー名のみを入力"
            className="input input-bordered"
          />
        </div>
        <div className="form-control pt-4">
          <div className="label">
            <span className="label-text">Twitter ID</span>
          </div>
          <input
            {...register("twitter_id")}
            type="text"
            placeholder="@マークなしで入力"
            className="input input-bordered"
          />
        </div>
        <div className="form-control pt-4">
          <div className="label">
            <span className="label-text">Facebook ID</span>
          </div>
          <input
            {...register("facebook_id")}
            type="text"
            placeholder="ユーザー名のみを入力"
            className="input input-bordered"
          />
        </div>
        <input type="submit" value="更新" className="btn btn-primary mt-4" />
      </form>
    </ICTSCCard>
  );
}

export default Page;
