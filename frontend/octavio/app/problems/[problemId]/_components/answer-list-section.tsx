import { useState } from "react";

import Image from "next/image";

import clsx from "clsx";
import { DateTime } from "luxon";

import ICTSCCard from "@/components/card";
import MarkdownPreview from "@/components/markdown-preview";
import useAnswers from "@/hooks/answer";
import { Problem } from "@/types/Problem";

type AnswerSectionProps = {
  problem: Problem;
};

function AnswerListSection({ problem }: AnswerSectionProps) {
  const [selectedAnswerId, setSelectedAnswerId] = useState<string | null>(null);
  const { answers, getAnswer } = useAnswers(problem?.id as string);
  const selectedAnswer = getAnswer(selectedAnswerId as string);
  const [isPreviewAnswer, setIsPreviewAnswer] = useState(true);

  return (
    <>
      <div className="text-sm pb-2">定期的に自動更新されます</div>
      <div className="overflow-x-auto">
        <table className="table border table-compact w-full">
          <thead>
            <tr>
              <th className="w-[196px]">提出日時</th>
              <th className="w-[100px]">問題コード</th>
              <th>問題</th>
              <th className="w-[100px]">得点</th>
              <th className="w-[100px]">チェック済み</th>
              <th className="w-[50px]" aria-label="投稿内容" />
              <th className="w-[50px]" aria-label="ダウンロード" />
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
                const blob = new Blob([`${getAnswer(answer.id)?.body}`], {
                  type: "text/markdown",
                });

                return (
                  <tr key={answer.id}>
                    <td>{createdAt.toFormat("yyyy-MM-dd HH:mm:ss")}</td>
                    <td>{problem?.code}</td>
                    <td>{problem?.title}</td>
                    <td className="text-right">{answer?.point ?? "--"} pt</td>
                    <td className="text-center">
                      {answer.point != null ? "○" : "採点中"}
                    </td>
                    <td>
                      <a
                        href="#preview"
                        className="link"
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
                        className="link"
                        href={URL.createObjectURL(blob)}
                      >
                        ダウンロード
                      </a>
                    </td>
                  </tr>
                );
              })}
            <tr />
          </tbody>
        </table>
      </div>
      {selectedAnswer != null && (
        <div className="pt-8" id="preview">
          <ICTSCCard>
            <div className="flex flex-row justify-between pb-4">
              <div className="answer-preview-team-info flex flex-row items-center">
                {selectedAnswer.point !== null && (
                  <div className="pr-2">
                    <Image
                      src="/assets/svg/check-green.svg"
                      height={24}
                      width={24}
                      alt="checked"
                    />
                  </div>
                )}
                チーム: {selectedAnswer.user_group.name}(
                {selectedAnswer.user_group.organization})
              </div>
              <div className="answer-preview-created-at">
                {DateTime.fromISO(selectedAnswer.created_at).toFormat(
                  "yyyy-MM-dd HH:mm:ss"
                )}
              </div>
            </div>
            <div className="flex flex-col">
              <div className="tabs">
                <button
                  type="button"
                  onClick={() => setIsPreviewAnswer(false)}
                  className={clsx(
                    `tab tab-lifted`,
                    !isPreviewAnswer && "tab-active"
                  )}
                >
                  Markdown
                </button>
                <button
                  type="button"
                  onClick={() => setIsPreviewAnswer(true)}
                  className={clsx(
                    `tab tab-lifted`,
                    isPreviewAnswer && "tab-active"
                  )}
                >
                  Preview
                </button>
              </div>
              {isPreviewAnswer ? (
                <MarkdownPreview
                  className="answer-preview"
                  content={selectedAnswer.body}
                />
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
}

export default AnswerListSection;
