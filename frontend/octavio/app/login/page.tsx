"use client";

import { useState } from "react";

import { useRouter } from "next/navigation";

import clsx from "clsx";
import { SubmitHandler, useForm } from "react-hook-form";

import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/Alerts";
import useAuth from "@/hooks/auth";

type Inputs = {
  name: string;
  password: string;
};

function Page() {
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<Inputs>();

  const { signIn, mutate } = useAuth();

  // ステータスコード
  const [status, setStatus] = useState<number | null>(null);
  // 送信中
  const [submitting, setSubmitting] = useState(false);

  const onSubmit: SubmitHandler<Inputs> = async (data) => {
    setSubmitting(true);
    const response = await signIn({
      name: data.name,
      password: data.password,
    });

    setSubmitting(false);
    setStatus(response.code);

    if (response.code === 200) {
      await mutate();
      await router.push("/");
    }
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="form-control flex flex-col container-ictsc items-center"
    >
      {status === 200 && (
        <ICTSCSuccessAlert className="mb-8" message="ログインに成功しました" />
      )}
      {status != null && status !== 200 && (
        <ICTSCErrorAlert className="mb-8" message="ログインに失敗しました" />
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
        {...register("password", { required: true })}
        type="password"
        placeholder="パスワード"
        id="password"
        className="input input-bordered max-w-xs min-w-[312px] mt-4"
      />
      <div className="label max-w-xs min-w-[312px]">
        {errors.password && (
          <span className="label-text-alt text-error">
            パスワードを入力して下さい
          </span>
        )}
      </div>
      <button
        type="submit"
        id="loginBtn"
        className={clsx(
          "btn btn-primary mt-4 max-w-xs min-w-[312px]",
          submitting && "loading"
        )}
      >
        ログイン
      </button>
    </form>
  );
}

export default Page;
