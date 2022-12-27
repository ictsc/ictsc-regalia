import 'zenn-content-css';
import {useState} from "react";
import {useRouter} from "next/router";
import Error from "next/error";

import {useForm, Controller, SubmitHandler} from "react-hook-form";

import ICTSCNavBar from "../../components/Navbar";
import ICTSCCard from "../../components/Card";
import MarkdownPreview from "../../components/MarkdownPreview";
import LoadingPage from "../../components/LoadingPage";
import {useApi} from "../../hooks/api";
import {useAuth} from "../../hooks/auth";
import {useProblems} from "../../hooks/problem";
import {AnswerResult, Result} from "../../types/_api";

type Inputs = {
  answer: string;
}

const ProblemPage = () => {
  const router = useRouter();

  const {handleSubmit, control, watch} = useForm<Inputs>()
  // answer のフォームを監視
  const watchField = watch(['answer'])

  const {apiClient} = useApi()
  const [isPreview, setIsPreview] = useState(false);
  const {user} = useAuth()
  const {getProblem, loading} = useProblems();


  const {problemId} = router.query;
  const problem = getProblem(problemId as string);

  const onSubmit: SubmitHandler<Inputs> = async ({answer}) => {
    const response = await apiClient.post(`problems/${problem?.id}/answers`, {
      json: {
        user_group_id: user?.user_group_id,
        problem_id: problem?.id,
        body: answer,
      }
    }).json<Result<AnswerResult>>()

    // TODO(k-shir0): 投稿に成功しましたアラートを表示
    // TODO(k-shir0): 投稿に失敗しましたアラートを表示
    console.log(response)
  }


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
        <div className={'container-ictsc'}>
          <div className={'flex flex-row items-end py-12'}>
            <h1 className={'title-ictsc pr-4'}>{problem.title}</h1>
            満点
            {problem.point} pt
            採点基準
            {problem.solved_criterion} pt
          </div>
          <ICTSCCard>
            <MarkdownPreview content={problem.body}/>
          </ICTSCCard>
          <ICTSCCard className={'mt-8 pt-4'}>
            <form onSubmit={handleSubmit(onSubmit)} className={'flex flex-col'}>
              <div className="tabs">
                <a onClick={() => setIsPreview(false)}
                   className={`tab tab-lifted ${!isPreview && 'tab-active'}`}>Write</a>
                <a onClick={() => setIsPreview(true)}
                   className={`tab tab-lifted ${isPreview && 'tab-active'}`}>Preview</a>
              </div>
              {isPreview
                  ? (
                      <>
                        <MarkdownPreview className={'pt-4'} content={watchField[0]}/>
                        <div className="divider mt-0"/>
                      </>
                  )
                  : (
                      <Controller name={'answer'} control={control} render={({field}) => {
                        return (<textarea {...field} className="textarea textarea-bordered my-4 px-2 min-h-[300px]"
                                          placeholder={`お世話になっております。チーム○○です。
この問題ではxxxxxが原因でトラブルが発生したと考えられました。
そのため、以下のように設定を変更し、○○が正しく動くことを確認いたしました。
確認のほどよろしくお願いします。
### 手順
1. /etc/hoge/hoo.bar の編集`}/>)
                      }}/>
                  )}
              <div className={'flex justify-end'}>
                <input type={"submit"} className={'btn btn-primary max-w-[312px]'} value={"提出"}/>
              </div>
            </form>
          </ICTSCCard>
          <div className={'text-sm pt-2'}>※ 回答は15分に1度のみです</div>
        </div>
      </>
  )
}

export default ProblemPage
