import toast from "react-hot-toast";
import {ICTSCErrorAlert, ICTSCSuccessAlert} from "@/components/alerts";
import {answerLimit} from "@/components/_const";
import clsx from "clsx";

export const successNotify = () =>
    toast.custom((t) => (
        <ICTSCSuccessAlert
            className={clsx("mt-2", t.visible ? "animate-enter" : "animate-leave")}
            message="投稿に成功しました"
        />
    ));

export const errorNotify = () =>
    toast.custom((t) => (
        <ICTSCErrorAlert
            className={clsx("mt-2", t.visible ? "animate-enter" : "animate-leave")}
            message="投稿に失敗しました"
            subMessage={
                answerLimit === undefined
                    ? undefined
                    : `回答は${answerLimit}分に1度のみです`
            }
        />
    ));