import { type ReactNode, type HTMLInputAutoCompleteAttribute } from "react";
import { useFormStatus } from "react-dom";
import { clsx } from "clsx";
import { Field, Label, Description, Input } from "@headlessui/react";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";
import type { SignUpResponse } from "../features/viewer/signup";

export function SignUpPage(
  props: {
    submit: (request: {
      invitationCode: string;
      name: string;
      displayName: string;
    }) => void;
    defaultName?: string;
    defaultDisplayName?: string;
  } & SignUpResponse,
) {
  return (
    <>
      <Title>アカウント登録</Title>
      <div className="grid h-full items-center justify-center">
        <form
          className="w-96 flex flex-col rounded-16 p-64 shadow-lg"
          onSubmit={(e) => {
            e.preventDefault();
            const form = e.target as HTMLFormElement;
            const formData = new FormData(form);
            props.submit({
              invitationCode: formData.get("invitation_code") as string,
              name: formData.get("screen_name") as string,
              displayName: formData.get("display_name") as string,
            });
          }}
        >
          <h2 className="text-24 font-bold">アカウント登録</h2>
          <div className="mt-16 grid grid-cols-2 gap-20">
            <TextField
              name="invitation_code"
              className="col-span-2"
              label="招待コード"
              placeholder="(メールをご確認ください)"
              ignoreCompletion
              errorMessage={
                props.invitationCodeError != null
                  ? {
                      required: "未入力です",
                      invalid: "無効なコードです",
                      team_full: "チームが満員です",
                    }[props.invitationCodeError]
                  : undefined
              }
            />
            <TextField
              name="screen_name"
              autoComplete="username"
              label="ID"
              placeholder="ictsc_taro"
              defaultValue={props.defaultName}
              errorMessage={
                props.nameError != null
                  ? {
                      required: "未入力です",
                      invalid: "内容に誤りがあります",
                      duplicate: "既に使われています",
                    }[props.nameError]
                  : undefined
              }
            />
            <TextField
              name="display_name"
              autoComplete="name"
              label="表示名"
              placeholder="ICTSC太郎"
              defaultValue={props.defaultDisplayName}
              errorMessage={
                props.displayNameError != null
                  ? {
                      required: "未入力です",
                      invalid: "内容に誤りがあります",
                    }[props.displayNameError]
                  : undefined
              }
            />
          </div>
          <div className="mt-64 flex items-center justify-end gap-24">
            {props.error && (
              <p className="text-16 font-bold text-primary">
                {
                  {
                    rate_limit: "リクエストが多すぎます",
                    invalid: "正しく入力されていない項目があります",
                    unknown: "エラーが発生しました",
                  }[props.error]
                }
              </p>
            )}
            <Submit />
          </div>
        </form>
      </div>
    </>
  );
}

function TextField(props: {
  name: string;
  className?: string;
  label: ReactNode;
  errorMessage?: ReactNode;
  placeholder?: string;
  defaultValue?: string;
  autoComplete?: HTMLInputAutoCompleteAttribute;
  ignoreCompletion?: boolean;
}) {
  const { pending } = useFormStatus();
  return (
    <Field className={props.className} disabled={pending}>
      <div className="flex items-center gap-16 text-16 font-bold">
        <Label>{props.label}</Label>
        {props.errorMessage && (
          <Description className="text-primary">
            {props.errorMessage}
          </Description>
        )}
      </div>
      <Input
        name={props.name}
        type="text"
        className="mt-8 w-full rounded-12 bg-surface-2 px-12 py-8 transition"
        placeholder={props.placeholder}
        invalid={Boolean(props.errorMessage)}
        defaultValue={props.defaultValue}
        autoComplete={props.autoComplete}
        data-1p-ignore={props.ignoreCompletion}
      />
    </Field>
  );
}

function Submit() {
  const { pending } = useFormStatus();
  return (
    // FIXME:
    // useFormStatus には配下のコンポーネントで状態更新が発生すると状態がリセットされるバグがある(FYI: https://github.com/facebook/react/issues/30368)
    // これを回避するため，headlessui の Button コンポーネントを使わずに button 要素を使う
    <button
      type="submit"
      className={clsx(
        "group rounded-12 bg-surface-2 py-[14px] pl-[36px] pr-[28px] shadow-md transition",
        "hover:bg-surface-2/90 active:shadow-transparent",
      )}
      disabled={pending}
    >
      <div
        className={clsx("flex items-center gap-8 group-disabled:opacity-50")}
      >
        <span className="text-16 font-bold">登録</span>
        <MaterialSymbol size={24} icon="send" />
      </div>
    </button>
  );
}
