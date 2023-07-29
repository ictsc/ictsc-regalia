"use client";

import React, { useState } from "react";

import { useRouter } from "next/router";

import { SubmitHandler, useForm } from "react-hook-form";

import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/Alerts";
import useAuth from "@/hooks/auth";
import BaseLayout from "@/layouts/BaseLayout";

type Inputs = {
  name: string;
  password: string;
};

function Page() {
  const router = useRouter();
  const { user_group_id: userGroupId, invitation_code: invitationCode } =
    router.query;

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>();

  const { signUp } = useAuth();

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null);
  // 送信中
  const [submitting, setSubmitting] = useState(false);
  // エラーメッセージ
  const [message, setMessage] = useState<string | null>(null);

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    setSubmitting(true);
    const response = await signUp({
      name: data.name,
      password: data.password,
      user_group_id: userGroupId as string,
      invitation_code: invitationCode as string,
    });

    setSubmitting(false);
    setStatus(response.code);

    if (!(response.code === 201)) {
      const msg = response.data ?? "";
      if (msg.match(/Error 1062: Duplicate entry '\w+' for key 'name'/)) {
        setMessage("ユーザー名が重複しています。");
      }
      if (
        msg.match(
          /Error:Field validation for 'UserGroupID' failed on the 'required' tag/
        )
      ) {
        setMessage("無効なユーザーグループです。");
      }
      if (
        msg.match(
          /Error:Field validation for 'UserGroupID' failed on the 'uuid' tag/
        )
      ) {
        setMessage("無効なユーザーグループです。");
      }
      if (
        msg.match(
          /Error:Field validation for 'InvitationCode' failed on the 'required' tag/
        )
      ) {
        setMessage("無効な招待コードです。");
      }
    }

    if (response.code === 201) {
      await router.push("/login");
    }
  };

  return (
    <BaseLayout title="ユーザー登録">
      <h1 className="title-ictsc text-center py-12">ユーザー登録</h1>
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="form-control flex flex-col container-ictsc items-center"
      >
        {status === 201 && (
          <ICTSCSuccessAlert
            className="mb-8"
            message="ユーザー登録に成功しました！"
          />
        )}
        {status != null && status !== 201 && (
          <ICTSCErrorAlert
            className="mb-8"
            message="エラーが発生しました"
            subMessage={message ?? ""}
          />
        )}
        <input
          {...register("name", { required: true })}
          type="text"
          placeholder="ユーザー名"
          id="username"
          className="input input-bordered max-w-xs min-w-[312px]"
        />
        <div className="label max-w-xs min-w-[312px]">
          {errors.name && (
            <span className="label-text-alt text-error">
              ユーザー名を入力してください
            </span>
          )}
        </div>
        <input
          {...register("password", {
            required: true,
            minLength: 8,
          })}
          type="password"
          placeholder="パスワード"
          id="password"
          className="input input-bordered max-w-xs min-w-[312px] mt-4"
        />
        <div className="label max-w-xs min-w-[312px]">
          {errors.password?.type === "required" && (
            <span className="label-text-alt text-error">
              パスワードを入力して下さい
            </span>
          )}
          {errors.password?.type === "minLength" && (
            <span className="label-text-alt text-error">
              パスワードは8文字以上である必要があります
            </span>
          )}
        </div>
        <button
          type="submit"
          id="signUpBtn"
          className={`btn btn-primary mt-4 max-w-xs min-w-[312px] ${
            submitting && "loading"
          }`}
        >
          登録
        </button>
      </form>
    </BaseLayout>
  );
}

export default Page;
