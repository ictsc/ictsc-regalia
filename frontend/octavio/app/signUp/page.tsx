"use client";

import React, { useState } from "react";

import { useRouter, useSearchParams } from "next/navigation";

import { valibotResolver } from "@hookform/resolvers/valibot";
import clsx from "clsx";
import { SubmitHandler, useForm } from "react-hook-form";

import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/alerts";
import useAuth from "@/hooks/auth";
import { SignUpSchema, SignUpType } from "@/types/schema/signUp";

function Page() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const userGroupId = searchParams?.get("user_group_id");
  const invitationCode = searchParams?.get("invitation_code");

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignUpType>({
    resolver: valibotResolver(SignUpSchema),
  });

  const { signUp } = useAuth();

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null);
  // 送信中
  const [submitting, setSubmitting] = useState(false);
  // エラーメッセージ
  const [message, setMessage] = useState<string | null>(null);

  const onSubmit: SubmitHandler<SignUpType> = async (data) => {
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
        {...register("name")}
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
        {...register("password")}
        type="password"
        placeholder="パスワード"
        id="password"
        className="input input-bordered max-w-xs min-w-[312px] mt-4"
      />
      <div className="label max-w-xs min-w-[312px]">
        {errors.password && (
          <span className="label-text-alt text-error">
            {errors.password.message}
          </span>
        )}
      </div>
      <button
        type="submit"
        id="signUpBtn"
        className={clsx(
          "btn btn-primary mt-4 max-w-xs min-w-[312px]",
          submitting && "loading"
        )}
      >
        登録
      </button>
    </form>
  );
}

export default Page;
