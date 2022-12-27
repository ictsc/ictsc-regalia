import {useState} from "react";
import {useRouter} from "next/router";
import Error from "next/error";

import {useForm} from "react-hook-form";

import ICTSCNavBar from "../../components/Navbar";
import ICTSCCard from "../../components/Card";
import LoadingPage from "../../components/LoadingPage";
import MarkdownPreview from "../../components/MarkdownPreview";
import {useApi} from "../../hooks/api";
import {useProblems} from "../../hooks/problem";
import {useAnswers} from "../../hooks/answer";
import {Answer} from "../../types/Answer";
import {Problem} from "../../types/Problem";
import {Result} from "../../types/_api";

const ScoringProblem = () => {
  const router = useRouter()

  const {getProblem, loading} = useProblems();

  const {code} = router.query
  const problem = getProblem(code as string)
  const {answers} = useAnswers(problem?.id ?? "")
  const [showProblem, setShowProblem] = useState(true)


  if (loading) {
    return (
        <>
          <ICTSCNavBar/>
          <LoadingPage/>
        </>
    );
  }

  if (problem === null) {
    return <Error statusCode={404}/>;
  }


  return (
      <>
        <ICTSCNavBar/>
        <div className="container-ictsc">
          <div className="collapse collapse-arrow" onClick={() => setShowProblem(!showProblem)}>
            <input type="checkbox" checked={showProblem}/>
            <div className={'collapse-title flex flex-row items-end px-0 mt-8'}>
              <h1 className={'title-ictsc pr-4'}>{problem.title}</h1>
              満点
              {problem.point} pt
              採点基準
              {problem.solved_criterion} pt
            </div>
            <div className={'collapse-content px-0'}>
              <ICTSCCard className={'ml-0'}>
                <MarkdownPreview content={problem.body}/>
              </ICTSCCard>
            </div>
          </div>
          <div className={'divider'}/>
          {answers.map((answer) => (
              <AnswerForm key={answer.id} problem={problem} answer={answer}/>
          ))}
        </div>
      </>
  )
}

type AnswerFormProps = {
  problem: Problem
  answer: Answer
}

type AnswerFormInputs = {
  point: number
}

const AnswerForm = ({problem, answer}: AnswerFormProps) => {
  const {apiClient} = useApi();
  const {mutate} = useAnswers(problem.id)

  const {register, handleSubmit, formState: {errors}} = useForm<AnswerFormInputs>({
    defaultValues: {
      point: answer.point ?? undefined
    }
  })

  const onSubmit = async (data: AnswerFormInputs) => {


    await apiClient.patch(`problems/${problem.id}/answers/${answer.id}`, {
      json: {
        problem_id: problem.id,
        answer_id: answer.id,
        // parseInt するとダブルクォートが取り除かれる
        point: parseInt(data.point.toString()),
      }
    }).json<Result<Answer>>()

    await mutate()
  }


  return (
      <form onSubmit={handleSubmit(onSubmit)}>
        <ICTSCCard key={answer.id} className={'pt-4 mb-4'}>
          <div className={'flex flex-row justify-between pb-4'}>
            <div className={'flex flex-row items-center'}>
              {answer.point !== null && (
                  <div className={'pr-2'}>
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={3}
                         stroke="green" className="w-6 h-6">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5"/>
                    </svg>
                  </div>
              )}
              チーム: {answer.user_group.name}({answer.user_group.organization})
            </div>
            <div>
              {answer.created_at.toString()}
            </div>
          </div>
          <MarkdownPreview content={problem.body}/>
          <div className={"divider"}/>
          <div className={'flex flex-row'}>
            <input {...register("point", {
              required: true,
              min: 0,
              max: problem.point,
            })} type={"text"} className={"input input-bordered input-sm"}/>
            <input type={"submit"} className={"btn btn-primary btn-sm ml-2"} value={"採点"}/>
          </div>
          <label className="label">
            {errors.point?.type === 'required' &&
              <span className="label-text-alt text-error">点数を入力して下さい</span>}
            {errors.point?.type === 'min' &&
              <span className="label-text-alt text-error">点数が低すぎます0以上の値を指定して下さい</span>}
            {errors.point?.type === 'max' &&
              <span className="label-text-alt text-error">点数が高すぎます{problem.point}以下の値を指定して下さい</span>}
          </label>
        </ICTSCCard>
      </form>
  )
}

export default ScoringProblem