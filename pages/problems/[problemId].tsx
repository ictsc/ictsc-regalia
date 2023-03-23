import "zenn-content-css";

import { useState } from "react";

import Error from "next/error";
import Image from "next/image";
import { useRouter } from "next/router";

import { DateTime } from "luxon";
import { useForm, Controller, SubmitHandler } from "react-hook-form";
import toast, { Toaster } from "react-hot-toast";

import { ICTSCErrorAlert, ICTSCSuccessAlert } from "@/components/Alerts";
import ICTSCCard from "@/components/Card";
import LoadingPage from "@/components/LoadingPage";
import MarkdownPreview from "@/components/MarkdownPreview";
import ProblemConnectionInfo from "@/components/ProblemConnectionInfo";
import ProblemMeta from "@/components/ProblemMeta";
import ProblemTitle from "@/components/ProblemTitle";
import { answerLimit, recreateRule } from "@/components/_const";
import { useAnswers } from "@/hooks/answer";
import { useApi } from "@/hooks/api";
import { useAuth } from "@/hooks/auth";
import { useProblem } from "@/hooks/problem";
import { useReCreateInfo } from "@/hooks/reCreateInfo";
import BaseLayout from "@/layouts/BaseLayout";
import { Problem } from "@/types/Problem";

type Inputs = {
  answer: string;
};

const AnswerForm = ({ code }: { code: string | null }) => {
  const {
    handleSubmit,
    control,
    watch,
    formState: { errors },
  } = useForm<Inputs>();
  // answer のフォームを監視
  const watchField = watch(["answer"]);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isPreview, setIsPreview] = useState(false);

  const { apiClient } = useApi();
  const { user } = useAuth();
  const { problem } = useProblem(code);
  const { mutate } = useAnswers(problem?.id ?? null);

  const onSubmit: SubmitHandler<Inputs> = async ({ answer }) => {
    const response = await apiClient.post(`problems/${problem?.id}/answers`, {
      json: {
        user_group_id: user?.user_group_id,
        problem_id: problem?.id,
        body: answer,
      },
    });

    if (response.ok) {
      successNotify();

      await mutate();
    } else {
      errorNotify();
    }
  };

  // モーダルを表示しバリデーションを行う
  const onModal: SubmitHandler<Inputs> = async () => {
    setIsModalOpen(true);
  };

  const successNotify = () =>
    toast.custom((t) => (
      <ICTSCSuccessAlert
        className={`mt-2 ${t.visible ? "animate-enter" : "animate-leave"}`}
        message={"投稿に成功しました"}
      />
    ));

  const errorNotify = () =>
    toast.custom((t) => (
      <ICTSCErrorAlert
        className={`mt-2 ${t.visible ? "animate-enter" : "animate-leave"}`}
        message={"投稿に失敗しました"}
        subMessage={
          answerLimit == undefined
            ? undefined
            : `回答は${answerLimit}分に1度のみです`
        }
      />
    ));

  return (
    <ICTSCCard className={"mt-8 pt-4"}>
      <div className={`modal ${isModalOpen && "modal-open"}`}>
        <div className="modal-box container-ictsc">
          <h3 className="title-ictsc pt-4 pb-8">回答内容確認</h3>
          <ICTSCCard>
            <MarkdownPreview content={watchField[0]} />
          </ICTSCCard>
          {answerLimit && (
            <div className={"text-sm pt-2"}>
              ※ 回答は{answerLimit}分に1度のみです
            </div>
          )}
          <div className="modal-action">
            <label
              onClick={() => setIsModalOpen(false)}
              className="btn btn-link"
            >
              閉じる
            </label>
            <label
              onClick={() => {
                handleSubmit(onSubmit)();
                setIsModalOpen(false);
              }}
              className="btn btn-primary"
            >
              この内容で提出
            </label>
          </div>
        </div>
      </div>
      <form onSubmit={handleSubmit(onModal)} className={"flex flex-col"}>
        <div className="tabs">
          <a
            onClick={() => setIsPreview(false)}
            className={`tab tab-lifted ${!isPreview && "tab-active"}`}
          >
            Markdown
          </a>
          <a
            onClick={() => setIsPreview(true)}
            className={`tab tab-lifted ${isPreview && "tab-active"}`}
          >
            Preview
          </a>
        </div>
        {isPreview ? (
          <>
            <MarkdownPreview className={"pt-4"} content={watchField[0]} />
            <div className="divider mt-0" />
          </>
        ) : (
          <Controller
            name={"answer"}
            control={control}
            rules={{ required: true }}
            render={({ field }) => {
              return (
                <textarea
                  {...field}
                  className="textarea textarea-bordered mt-4 px-2 min-h-[300px]"
                  placeholder={`お世話になっております。チーム○○です。
この問題ではxxxxxが原因でトラブルが発生したと考えられました。
そのため、以下のように設定を変更し、○○が正しく動くことを確認いたしました。
確認のほどよろしくお願いします。

### 手順
1. /etc/hoge/hoo.bar の編集`}
                />
              );
            }}
          />
        )}
        <label className="label max-w-xs min-w-[312px]">
          {errors.answer && (
            <span className="label-text-alt text-error">
              回答を入力して下さい
            </span>
          )}
        </label>
        <div className={"flex justify-end mt-4"}>
          <label
            onClick={handleSubmit(onModal)}
            className="btn btn-primary max-w-[312px]"
          >
            提出確認
          </label>
        </div>
      </form>
    </ICTSCCard>
  );
};

type AnswerSectionProps = {
  problem: Problem;
};

const AnswerListSection = ({ problem }: AnswerSectionProps) => {
  const [selectedAnswerId, setSelectedAnswerId] = useState<string | null>(null);
  const { answers, getAnswer } = useAnswers(problem?.id as string);
  const selectedAnswer = getAnswer(selectedAnswerId as string);
  const [isPreviewAnswer, setIsPreviewAnswer] = useState(true);

  return (
    <>
      <div className={"text-sm pb-2"}>定期的に自動更新されます</div>
      <div className={"overflow-x-auto"}>
        <table className="table border table-compact w-full">
          <thead>
            <tr>
              <th className={"w-[196px]"}>提出日時</th>
              <th className={"w-[100px]"}>問題コード</th>
              <th>問題</th>
              <th className={"w-[100px]"}>得点</th>
              <th className={"w-[100px]"}>チェック済み</th>
              <th className={"w-[50px]"}></th>
              <th className={"w-[50px]"}></th>
            </tr>
          </thead>
          <tbody>
            {answers
              .sort((a, b) => {
                if (a.created_at < b.created_at) {
                  return 1;
                }
                if (a.created_at > b.created_at) {
                  return -1;
                }
                return 0;
              })
              .map((answer) => {
                const createdAt = DateTime.fromISO(answer.created_at);
                let blob = new Blob(["" + getAnswer(answer.id)?.body], {
                  type: "text/markdown",
                });

                return (
                  <tr key={answer.id}>
                    <td>{createdAt.toFormat("yyyy-MM-dd HH:mm:ss")}</td>
                    <td>{problem?.code}</td>
                    <td>{problem?.title}</td>
                    <td className={"text-right"}>{answer?.point ?? "--"} pt</td>
                    <td className={"text-center"}>
                      {answer.point != null ? "○" : "採点中"}
                    </td>
                    <td>
                      <a
                        href={"#preview"}
                        className={"link"}
                        onClick={() => setSelectedAnswerId(answer.id)}
                      >
                        投稿内容
                      </a>
                    </td>
                    <td>
                      <a
                        download={`ictsc-${
                          problem?.code
                        }-${createdAt.toUnixInteger()}.md`}
                        className={"link"}
                        href={URL.createObjectURL(blob)}
                      >
                        ダウンロード
                      </a>
                    </td>
                  </tr>
                );
              })}
            <tr></tr>
          </tbody>
        </table>
      </div>
      {selectedAnswer != null && (
        <div className={"pt-8"} id={"preview"}>
          <ICTSCCard>
            <div className={"flex flex-row justify-between pb-4"}>
              <div className={"flex flex-row items-center"}>
                {selectedAnswer.point !== null && (
                  <div className={"pr-2"}>
                    <Image
                      src={"/assets/svg/check-green.svg"}
                      height={24}
                      width={24}
                      alt={"checked"}
                    />
                  </div>
                )}
                チーム: {selectedAnswer.user_group.name}(
                {selectedAnswer.user_group.organization})
              </div>
              <div>
                {DateTime.fromISO(selectedAnswer.created_at).toFormat(
                  "yyyy-MM-dd HH:mm:ss"
                )}
              </div>
            </div>
            <div className="flex flex-col">
              <div className="tabs">
                <a
                  onClick={() => setIsPreviewAnswer(false)}
                  className={`tab tab-lifted ${
                    !isPreviewAnswer && "tab-active"
                  }`}
                >
                  Markdown
                </a>
                <a
                  onClick={() => setIsPreviewAnswer(true)}
                  className={`tab tab-lifted ${
                    isPreviewAnswer && "tab-active"
                  }`}
                >
                  Preview
                </a>
              </div>
              {isPreviewAnswer ? (
                <MarkdownPreview content={selectedAnswer.body} />
              ) : (
                <textarea
                  readOnly
                  className="textarea textarea-bordered mt-4 px-2 min-h-[300px]"
                >
                  {selectedAnswer.body}
                </textarea>
              )}
            </div>
          </ICTSCCard>
        </div>
      )}
    </>
  );
};

const ProblemPage = () => {
  const router = useRouter();
  const { problemId } = router.query;

  const { apiClient } = useApi();
  const { user } = useAuth();
  const { matter, problem, isLoading } = useProblem(problemId as string | null);
  const [isReCreateModalOpen, setIsReCreateModalOpen] = useState(false);

  const { recreateInfo, mutate: recreateMutate } = useReCreateInfo(
    problem?.code ?? null
  );

  const isReadOnly = user?.is_read_only ?? false;
  const onReCreateSubmit = async () => {
    const response = await apiClient.post(`recreate/${problem?.code}`);

    if (response.ok) {
      await recreateMutate();
    }
  };

  if (isLoading) {
    return (
      <BaseLayout title={"Loading..."}>
        <LoadingPage />
      </BaseLayout>
    );
  }

  if (problem === null) {
    return <Error statusCode={404} />;
  }

  return (
    <>
      <Toaster />
      <input type="checkbox" id="my-modal-5" className="modal-toggle" />
      <div className={`modal ${isReCreateModalOpen && "modal-open"}`}>
        <div className="modal-box container-ictsc">
          <h3 className="title-ictsc pt-4 pb-8">
            問題の再展開を行います。よろしいですか？
          </h3>
          <MarkdownPreview
            content={recreateRule?.replace(/\\n/g, "\n") ?? ""}
          />
          <div className="modal-action">
            <label
              onClick={() => setIsReCreateModalOpen(false)}
              className="btn btn-link"
            >
              閉じる
            </label>
            <label
              onClick={() => {
                onReCreateSubmit();
                setIsReCreateModalOpen(false);
              }}
              className="btn btn-primary"
            >
              問題の再展開を行う
            </label>
          </div>
        </div>
      </div>
      <BaseLayout title={`${problem.code} ${problem.title} 問題`}>
        <div className={"container-ictsc"}>
          <div
            className={
              "flex flex-row justify-between pt-12 justify-items-center"
            }
          >
            <ProblemTitle title={problem.title} />
            {!isReadOnly && (
              <button
                className="btn text-red-500 btn-sm"
                onClick={() => {
                  setIsReCreateModalOpen(true);
                }}
                disabled={
                  recreateInfo?.available != null &&
                  !(recreateInfo?.available ?? false)
                }
              >
                再展開を行う
              </button>
            )}
          </div>
          <ProblemMeta problem={problem} />
          <ProblemConnectionInfo matter={matter} />
          {recreateInfo?.available != null &&
            !(recreateInfo?.available ?? false) && (
              <div className={`alert alert-info shadow-lg grow`}>
                <div>
                  <div className={"animate-spin"}>
                    <Image
                      src={"/assets/svg/arrow-path.svg"}
                      height={24}
                      width={24}
                      alt={"recreate"}
                    />
                  </div>
                  <div className={"flex flex-col"}>
                    <span>問題を再展開中です</span>
                  </div>
                </div>
              </div>
            )}
          <ICTSCCard className={"mt-8"}>
            <MarkdownPreview content={problem.body ?? ""} />
          </ICTSCCard>

          {!isReadOnly && <AnswerForm code={problemId as string | null} />}
          {answerLimit && (
            <div className={"text-sm pt-2"}>
              ※ 回答は{answerLimit}分に1度のみです
            </div>
          )}
          <div className={"divider"} />
          <AnswerListSection problem={problem} />
        </div>
      </BaseLayout>
    </>
  );
};

export default ProblemPage;
