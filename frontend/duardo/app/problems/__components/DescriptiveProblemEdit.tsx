"use client";

import React, { useState } from "react";

import { useRouter } from "next/navigation";

import { zodResolver } from "@hookform/resolvers/zod";
import * as Toast from "@radix-ui/react-toast";
import MarkdownPreview from "@repo/ui/markdown-preview";
import clsx from "clsx";
import { useFieldArray, useForm } from "react-hook-form";

import ProblemConnectionInfo from "@/app/problems/edit/[id]/__components/ConnectionInfo";
import MarkdownEdit from "@/app/problems/edit/[id]/__components/MarkdownEdit";
import {
  formSchema,
  FormSchema,
} from "@/app/problems/edit/[id]/__type/MarkdownEditFormSchema";
import Card from "@/components/card";
import { ConnectionInfo, Problem } from "@/proto/admin/v1/problem_pb";

type ToastType = "info" | "error" | "success" | "warning" | undefined;
export type Props = {
  problem?: Problem;
  saveProblemData: (data: FormSchema) => Promise<boolean>;
};

function DescriptiveProblemEdit({ problem, saveProblemData }: Props) {
  const router = useRouter();

  const [text, setText] = useState("");

  const timerRef = React.useRef(0);

  const [toastText, setToastText] = useState("");
  const [toastType, setToastType] = useState<ToastType>();
  const [open, setOpen] = useState(false);
  const [preview, setPreview] = useState(false);

  React.useEffect(() => () => clearTimeout(timerRef.current), []);

  const defaultValues =
    problem?.body.case === "descriptive"
      ? {
          title: problem.title,
          code: problem.code,
          point: problem.point,
          body: problem.body.value.body,
          connectionInfos: problem.body.value.connectionInfos,
        }
      : {};

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    watch,
  } = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues,
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: "connectionInfos",
  });

  const onSubmit = handleSubmit(async (data) => {
    const success = await saveProblemData(data);

    setOpen(false);
    timerRef.current = window.setTimeout(() => {
      if (success) {
        setToastType("success");
        setToastText("問題を保存しました");
        router.push("/problems");
      } else {
        setToastType("error");
        setToastText("不明なエラーが発生しました");
      }
      setOpen(true);
    }, 100);
  });

  return (
    <form className="container-xl mx-auto px-2 flex-grow" onSubmit={onSubmit}>
      <div className="mt-16 mb-6 flex flex-col">
        <input
          type="text"
          placeholder="タイトル"
          className="text-2xl font-bold input pl-0 focus:border-0 border-0 focus:outline-0 min-w-full"
          {...register("title")}
        />
        <div className="flex flex-row">
          <div>
            <span className="text-gray-500 text-xs">問題コード</span>
            <input
              type="text"
              placeholder="ABC"
              className="text-xl font-bold input pl-0 focus:border-0 border-0 focus:outline-0 w-[80px]"
              {...register("code")}
            />
          </div>
          <div>
            <input
              type="number"
              placeholder="100"
              className="text-xl font-bold input px-0 focus:border-0 border-0 focus:outline-0 w-[60px] text-right"
              {...register("point", {
                valueAsNumber: true,
              })}
            />
            <span className="text-gray-500 text-xs">ポイント</span>
          </div>
        </div>
        <div hidden={preview}>
          {fields.map((field, index) => (
            <div key={field.id} className="flex flex-row items-center pt-2">
              <button
                type="button"
                aria-label="削除"
                className="btn btn-circle btn-ghost btn-xs"
                onClick={() => remove(index)}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
              <input
                type="text"
                placeholder="ホスト名"
                className="input input-bordered input-sm ml-2"
                {...register(`connectionInfos.${index}.hostname` as const)}
              />
              <input
                type="text"
                placeholder="コマンド"
                className="input input-bordered input-sm ml-2 w-1/4 min-w-[200px]"
                {...register(`connectionInfos.${index}.command` as const)}
              />
              <input
                type="text"
                placeholder="パスワード"
                className="input input-bordered input-sm ml-2"
                {...register(`connectionInfos.${index}.password` as const)}
              />
              <input
                type="text"
                placeholder="タイプ"
                className="input input-bordered input-sm ml-2"
                {...register(`connectionInfos.${index}.type` as const)}
              />
            </div>
          ))}
          <button
            type="button"
            className="btn btn-sm my-4"
            onClick={() => append(new ConnectionInfo())}
          >
            接続情報を追加
          </button>
        </div>
      </div>
      <div className="pb-12" hidden={!preview}>
        <ProblemConnectionInfo connectionInfos={watch("connectionInfos")} />
      </div>

      <Card className="px-8 pt-12 pb-8" hidden={preview}>
        <MarkdownEdit register={register} onChange={setText} />
      </Card>
      <Card className="px-8 pt-12 pb-8" hidden={!preview}>
        <div>
          <MarkdownPreview content={text} />
        </div>
      </Card>
      <div className="flex justify-between items-center mt-4">
        <Toast.Provider>
          <button type="submit" className="btn btn-primary">
            保存
          </button>

          <Toast.Root
            className="toast data-[state=open]:animate-slideIn data-[state=closed]:animate-hide data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)] data-[swipe=cancel]:translate-x-0 data-[swipe=cancel]:transition-[transform_200ms_ease-out] data-[swipe=end]:animate-swipeOut"
            open={open}
            onOpenChange={setOpen}
          >
            <div
              className={clsx(
                "alert",
                toastType === "success" && "alert-success",
                toastType === "error" && "alert-error",
              )}
            >
              <Toast.Title>{toastText}</Toast.Title>
            </div>
          </Toast.Root>
          <Toast.Viewport className="ToastViewport" />
        </Toast.Provider>
        <button
          type="button"
          className="btn btn-secondary"
          onClick={() => setPreview(!preview)}
        >
          {preview ? "プレビューを隠す" : "プレビュー"}
        </button>
      </div>
      <div className="text-red-500 pt-4">
        {errors.title && <p>{errors.title.message}</p>}
        {errors.code && <p>{errors.code.message}</p>}
        {errors.point && <p>{errors.point.message}</p>}
        {errors.body && <p>{errors.body.message}</p>}
      </div>
    </form>
  );
}

export default DescriptiveProblemEdit;
