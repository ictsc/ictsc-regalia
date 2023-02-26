import 'zenn-content-css';
import {useState} from "react";
import {useRouter} from "next/router";
import Error from "next/error";

import {useForm, Controller, SubmitHandler} from "react-hook-form";
import {DateTime} from "luxon";

import ICTSCNavBar from "../../components/Navbar";
import ICTSCCard from "../../components/Card";
import {ICTSCErrorAlert, ICTSCSuccessAlert} from "../../components/Alerts";
import MarkdownPreview from "../../components/MarkdownPreview";
import LoadingPage from "../../components/LoadingPage";
import HoverCopyText from "../../components/HoverCopyText";
import {useApi} from "../../hooks/api";
import {useAuth} from "../../hooks/auth";
import {useProblems} from "../../hooks/problem";
import {useAnswers} from "../../hooks/answer";
import {Problem} from "../../types/Problem";

type Inputs = {
    answer: string;
}


const ProblemPage = () => {
    const router = useRouter();
    const {problemId} = router.query;

    const {handleSubmit, control, watch, formState: {errors}} = useForm<Inputs>()
    // answer のフォームを監視
    const watchField = watch(['answer'])

    const {apiClient} = useApi()
    const {user} = useAuth()
    const {getProblem, isLoading,} = useProblems();
    const [isPreview, setIsPreview] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [status, setStatus] = useState<number | null>(null);

    const [matter, problem] = getProblem(problemId as string);

    const {mutate} = useAnswers(problem?.id as string);

    // モーダルを表示しバリデーションを行う
    const onModal: SubmitHandler<Inputs> = async () => {
        setIsModalOpen(true)
    }

    const onSubmit: SubmitHandler<Inputs> = async ({answer}) => {
        const response = await apiClient.post(`problems/${problem?.id}/answers`, {
            json: {
                user_group_id: user?.user_group_id,
                problem_id: problem?.id,
                body: answer,
            }
        })

        setStatus(response.status)

        if (response.ok) {
            await mutate()
        }
    }


    if (problem === null) {
        return <Error statusCode={404}/>;
    }


    if (isLoading) {
        return (
            <>
                <ICTSCNavBar/>
                <LoadingPage/>
            </>
        );
    }


    const answerLimit = process.env.NEXT_PUBLIC_ANSWER_LIMIT;

    return (
        <>
            <input type="checkbox" id="my-modal-5" className="modal-toggle"/>
            <div className={`modal ${isModalOpen && 'modal-open'}`}>
                <div className="modal-box container-ictsc">
                    <h3 className="title-ictsc pt-4 pb-8">回答内容確認</h3>
                    <ICTSCCard>
                        <MarkdownPreview content={watchField[0]}/>
                    </ICTSCCard>
                    {answerLimit && <div className={'text-sm pt-2'}>※ 回答は{answerLimit}分に1度のみです</div>}
                    <div className="modal-action">
                        <label onClick={() => setIsModalOpen(false)} className="btn btn-link">閉じる</label>
                        <label onClick={() => {
                            handleSubmit(onSubmit)()
                            setIsModalOpen(false)
                        }} className="btn btn-primary">この内容で提出</label>
                    </div>
                </div>
            </div>

            <ICTSCNavBar/>
            <div className={'container-ictsc'}>
                <div
                    className={`collapse collapse-problem collapse-arrow pt-12 px-0`}>
                    <input type="checkbox" defaultChecked={true}/>
                    <div className={'collapse-title flex flex-row justify-between pl-0 pr-9'}>
                        <div className={'flex flex-row items-end'}>
                            <h1 className={'title-ictsc pr-2'}>{problem.title}</h1>
                            満点
                            {problem.point} pt
                            採点基準
                            {problem.solved_criterion} pt
                        </div>
                        <div className={'text-sm flex flex-row items-end'}>接続情報表示/非表示
                        </div>
                    </div>
                    <div className={'collapse-content pt-2 px-0'}>
                        <div className="overflow-x-auto">
                            <table className="table table-compact w-full">
                                <thead>
                                <tr>
                                    <th>ホスト名</th>
                                    <th></th>
                                    <th>ユーザ</th>
                                    <th>パスワード</th>
                                    <th>ポート</th>
                                    <th>種類</th>
                                </tr>
                                </thead>
                                <tbody>
                                {matter?.connectInfo?.map((info, index) => (
                                    <tr key={index}>
                                        <th>{info.hostname}</th>
                                        <HoverCopyText text={info.command ?? ''}/>
                                        <HoverCopyText text={info.user ?? ''}/>
                                        <HoverCopyText text={info.password ?? ''}/>
                                        <td>{info.port}</td>
                                        <td>{info.type}</td>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                {status === 201 && (
                    <ICTSCSuccessAlert className={'mt-2'} message={'投稿に成功しました'}/>)}
                {status != null && status !== 201 && (
                    <ICTSCErrorAlert className={'mt-2'} message={'投稿に失敗しました'}
                                     subMessage={answerLimit == undefined ? undefined : `回答は${answerLimit}分に1度のみです`}/>)}
                <ICTSCCard className={'mt-8'}>
                    <MarkdownPreview content={problem.body ?? ""}/>
                </ICTSCCard>
                <ICTSCCard className={'mt-8 pt-4'}>
                    <form onSubmit={handleSubmit(onModal)} className={'flex flex-col'}>
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
                                <Controller name={'answer'} control={control}
                                            rules={{required: true}}
                                            render={({field}) => {
                                                return (<textarea {...field}
                                                                  className="textarea textarea-bordered mt-4 px-2 min-h-[300px]"
                                                                  placeholder={`お世話になっております。チーム○○です。
この問題ではxxxxxが原因でトラブルが発生したと考えられました。
そのため、以下のように設定を変更し、○○が正しく動くことを確認いたしました。
確認のほどよろしくお願いします。

### 手順
1. /etc/hoge/hoo.bar の編集`}/>)
                                            }}/>
                            )}
                        <label className="label max-w-xs min-w-[312px]">
                            {errors.answer && <span className="label-text-alt text-error">回答を入力して下さい</span>}
                        </label>
                        <div className={'flex justify-end mt-4'}>
                            <label onClick={handleSubmit(onModal)}
                                   className="btn btn-primary max-w-[312px]">提出確認</label>
                        </div>
                    </form>
                </ICTSCCard>
                {answerLimit && <div className={'text-sm pt-2'}>※ 回答は{answerLimit}分に1度のみです</div>}
                <div className={'divider'}/>
                <AnswerListSection problem={problem}/>
            </div>
        </>
    )
}

type AnswerSectionProps = {
    problem: Problem
}

const AnswerListSection = ({problem}: AnswerSectionProps) => {
    const [selectedAnswerId, setSelectedAnswerId] = useState<string | null>(null);
    const {answers, getAnswer, mutate} = useAnswers(problem?.id as string);
    const selectedAnswer = getAnswer(selectedAnswerId as string);

    return (
        <>
            <div className={'text-sm pb-2'}>
                定期的に自動更新されます
            </div>
            <div className={'overflow-x-auto'}>
                <table className="table border table-compact w-full">
                    <thead>
                    <tr>
                        <th className={'w-[196px]'}>提出日時</th>
                        <th className={'w-[100px]'}>問題コード</th>
                        <th>問題</th>
                        <th className={'w-[100px]'}>得点</th>
                        <th className={'w-[100px]'}>チェック済み</th>
                        <th className={'w-[50px]'}></th>
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
                            const createdAt = DateTime.fromISO(answer.created_at)

                            return (
                                <tr key={answer.id}>
                                    <td>{createdAt.toFormat('yyyy-MM-dd HH:mm:ss')}</td>
                                    <td>{problem?.code}</td>
                                    <td>{problem?.title}</td>
                                    <td className={'text-right'}>{answer?.point ?? '--'} pt</td>
                                    <td className={'text-center'}>{answer.point != null ? '○' : '採点中'}</td>
                                    <td>
                                        <a href={'#preview'} className={'link'}
                                           onClick={() => setSelectedAnswerId(answer.id)}>投稿内容</a>
                                    </td>
                                </tr>
                            );
                        })}
                    <tr>
                    </tr>
                    </tbody>
                </table>
            </div>
            {selectedAnswer != null && (
                <div className={'pt-8'} id={'preview'}>
                    <ICTSCCard>
                        <div className={'flex flex-row justify-between pb-4'}>
                            <div className={'flex flex-row items-center'}>
                                {selectedAnswer.point !== null && (
                                    <div className={'pr-2'}>
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                             strokeWidth={3}
                                             stroke="green" className="w-6 h-6">
                                            <path strokeLinecap="round" strokeLinejoin="round"
                                                  d="M4.5 12.75l6 6 9-13.5"/>
                                        </svg>
                                    </div>
                                )}
                                チーム: {selectedAnswer.user_group.name}({selectedAnswer.user_group.organization})
                            </div>
                            <div>
                                {DateTime.fromISO(selectedAnswer.created_at).toFormat('yyyy-MM-dd HH:mm:ss')}
                            </div>
                        </div>
                        <MarkdownPreview content={selectedAnswer.body}/>
                    </ICTSCCard>
                </div>
            )}
        </>
    )
}

export default ProblemPage
