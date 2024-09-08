import Image from "next/image";

import { useForm } from "react-hook-form";

import ICTSCCard from "@/components/card";
import MarkdownPreview from "@/components/markdown-preview";
import useAnswers from "@/hooks/answer";
import useApi from "@/hooks/api";
import useProblems from "@/hooks/problems";
import { Answer } from "@/types/Answer";
import { Problem } from "@/types/Problem";

type AnswerFormProps = {
  problem: Problem;
  answer: Answer;
};

type AnswerFormInputs = {
  point: number;
};

function ScoringAnswerForm({ problem, answer }: AnswerFormProps) {
  const { client } = useApi();
  const { mutate } = useAnswers(problem.id);
  const { mutate: mutateProblem } = useProblems();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<AnswerFormInputs>({
    defaultValues: {
      point: answer.point ?? undefined,
    },
  });

  const onSubmit = async (data: AnswerFormInputs) => {
    await client.patch<Answer>(`problems/${problem.id}/answers/${answer.id}`, {
      problem_id: problem.id,
      answer_id: answer.id,
      // parseInt するとダブルクォートが取り除かれる
      point: parseInt(data.point.toString(), 10),
    });

    await mutate();
    await mutateProblem();
  };

  // yyyy/mm/dd hh:mm:ss
  const createdAt = new Date(Date.parse(answer.created_at)).toLocaleDateString(
    "ja-JP",
    {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    }
  );

  return (
    <ICTSCCard key={answer.id} className="pt-4 mb-4">
      <div className="flex flex-row justify-between pb-4">
        <div className="answer-preview-team-info flex flex-row items-center">
          {answer.point !== null && (
            <div className="pr-2">
              <Image
                src="/assets/svg/check-green.svg"
                height={24}
                width={24}
                alt="checked"
              />
            </div>
          )}
          チーム: {answer.user_group.name}({answer.user_group.organization})
        </div>
        <div className="answer-preview-created-at">{createdAt}</div>
      </div>
      <MarkdownPreview className="answer-preview" content={answer.body} />
      <div className="divider" />
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-row">
        <input
          {...register("point", {
            required: true,
            min: 0,
            max: problem.point,
          })}
          type="text"
          className="input input-bordered input-sm"
        />
        <input
          type="submit"
          className="btn btn-primary btn-sm ml-2"
          value="採点"
        />
      </form>
      <div className="label">
        {errors.point?.type === "required" && (
          <span className="label-text-alt text-error">
            点数を入力して下さい
          </span>
        )}
        {errors.point?.type === "min" && (
          <span className="label-text-alt text-error">
            点数が低すぎます0以上の値を指定して下さい
          </span>
        )}
        {errors.point?.type === "max" && (
          <span className="label-text-alt text-error">
            点数が高すぎます{problem.point}以下の値を指定して下さい
          </span>
        )}
      </div>
    </ICTSCCard>
  );
}

export default ScoringAnswerForm;
