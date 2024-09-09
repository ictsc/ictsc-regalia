"use client";

import React, { useState } from "react";

import { useRouter, useSearchParams } from "next/navigation";

import clsx from "clsx";
import { SubmitHandler, useForm } from "react-hook-form";

import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/alerts";
import useAuth from "@/hooks/auth";

type Inputs = {
  name: string;
  password: string;
};

function Page() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const userGroupId = searchParams?.get("user_group_id");
  const invitationCode = searchParams?.get("invitation_code");

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
    signUp({
      name: data.name,
      password: data.password,
      user_group_id: userGroupId as string,
      invitation_code: invitationCode as string,
    })
      .then((res) => {
        setStatus(res.code);

        if (res.code === 201) {
          router.push("/login");
        }
      })
      .catch((e) => {
        setStatus(e.code);

        const msg = e.response.data.error ?? "";
        if (msg.match(/Error 1062: Duplicate entry '\w+' for key 'name'/)) {
          setMessage("ユーザー名が重複しています。");
        }
        if (
          msg.match(
            /Error:Field validation for 'UserGroupID' failed on the 'required' tag/,
          )
        ) {
          setMessage("無効なユーザーグループです。");
        }
        if (
          msg.match(
            /Error:Field validation for 'UserGroupID' failed on the 'uuid' tag/,
          )
        ) {
          setMessage("無効なユーザーグループです。");
        }
        if (
          msg.match(
            /Error:Field validation for 'InvitationCode' failed on the 'required' tag/,
          )
        ) {
          setMessage("無効な招待コードです。");
        }
      })
      .finally(() => {
        setSubmitting(false);
      });
  };

  return (
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
        {...register("name", {
          required: "ユーザー名を入力してください",
          pattern: {
            value: /^[A-Za-z0-9_]{3,32}$/,
            message:
              "ユーザー名は3〜32文字のアルファベット、数字、アンダースコアのみが使用できます",
          },
        })}
        type="text"
        placeholder="ユーザー名"
        id="username"
        className="input input-bordered max-w-xs min-w-[312px]"
      />
      <div className="label max-w-xs min-w-[312px]">
        {errors.name && (
          <span className="label-text-alt text-error">
            {errors.name.message}
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
        className={clsx(
          "btn mt-4 max-w-xs min-w-[312px]",
          !submitting ? "btn-primary" : "btn-disabled",
        )}
        disabled={submitting}
      >
        {submitting && <span className="loading loading-spinner" />}
        登録
      </button>
    </form>
  );
}

export default Page;
