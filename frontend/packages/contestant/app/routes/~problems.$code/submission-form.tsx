import { useId, useEffect, useState, useActionState, Suspense } from "react";
import { Button, Field, Label, Textarea } from "@headlessui/react";
import { MaterialSymbol } from "../../components/material-symbol";
import { ConfirmModal } from "./confirmModal";
import { Markdown, Typography } from "@app/components/markdown";

interface AnswerableState {
  remainingSeconds: number;
}

function formatRemainingTimeParts(remainingSeconds: number) {
  const minutes = Math.floor(remainingSeconds / 60);
  const seconds = remainingSeconds % 60;
  return { minutes, seconds };
}

function useAnswerable(
  lastSubmittedAt?: string,
  submitInterval?: number,
): AnswerableState {
  const [state, setState] = useState<AnswerableState>({ remainingSeconds: 0 });

  useEffect(() => {
    if (!lastSubmittedAt || !submitInterval) {
      setState({ remainingSeconds: 0 });
      return;
    }

    const checkAnswerable = () => {
      const now = new Date();
      const lastSubmit = new Date(lastSubmittedAt);
      const nextSubmitTime = new Date(
        lastSubmit.getTime() + submitInterval * 1000,
      );

      const diffMs = nextSubmitTime.getTime() - now.getTime();
      const diffSec = Math.ceil(diffMs / 1000);

      const remainingSeconds = diffSec > 0 ? diffSec : 0;

      setState({ remainingSeconds });
    };

    checkAnswerable();
    const interval = setInterval(checkAnswerable, 1000);
    return () => clearInterval(interval);
  }, [lastSubmittedAt, submitInterval]);

  return state;
}

export function SubmissionForm(props: {
  readonly action: (answer: string) => Promise<"success" | "failure">;
  readonly submitInterval?: number;
  readonly lastSubmittedAt?: string;
  readonly storageKey?: string;
}) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [lastResult, action, isPending] = useActionState(
    async (_prevState: unknown, formData: FormData) => {
      const answer = formData.get("answer") as string;
      if (answer.trim() === "") {
        return {
          type: "error",
          name: "answer",
          error: "解答を入力してください",
        } as const;
      }

      const intent = formData.get("intent") as string;
      switch (intent) {
        case "confirm":
          setIsModalOpen(true);
          return { type: "confirm", answer } as const;
        case "submit":
          try {
            const result = await props.action(answer);
            if (result === "failure") {
              return {
                type: "error",
                error: "解答の送信に失敗しました",
              } as const;
            }
          } catch (err) {
            console.error(err);
            return {
              type: "error",
              error: "解答の送信に失敗しました",
            } as const;
          }
          return;
      }
    },
    null,
  );

  const { remainingSeconds } = useAnswerable(
    props.lastSubmittedAt,
    props.submitInterval,
  );
  const isAnswerable = remainingSeconds <= 0;
  const { minutes, seconds } = formatRemainingTimeParts(remainingSeconds);

  const formId = useId();
  const errorId = useId();

  return (
    <form
      id={formId}
      className="flex size-full flex-col"
      noValidate
      action={action}
    >
      <AnswerTextInputField
        invalid={lastResult?.type === "error" && lastResult.name === "answer"}
        defaultValue={
          lastResult?.type === "confirm" ? lastResult.answer : undefined
        }
        storageKey={props.storageKey}
      />
      <div className="mt-20 flex items-center justify-end gap-24">
        {!isAnswerable && (
          <div className="flex w-[160px] items-center justify-between">
            <span className="text-black text-16">解答可能まで</span>
            <div className="flex items-center">
              <span className="w-[24px] text-right text-20 font-bold text-primary">
                {minutes}
              </span>
              <span className="mx-2 text-20 font-bold text-primary">:</span>
              <span className="w-[24px] text-right text-20 font-bold text-primary">
                {seconds.toString().padStart(2, "0")}
              </span>
            </div>
          </div>
        )}
        {isAnswerable && lastResult?.type === "error" && (
          <label
            id={errorId}
            className="flex-shrink text-16 font-bold text-primary"
          >
            {lastResult.error}
          </label>
        )}
        <Button
          name="intent"
          value="confirm"
          type="submit"
          disabled={
            !isAnswerable || (lastResult?.type === "confirm" && isPending)
          }
          className="flex items-center justify-center self-end rounded-12 bg-surface-2 py-16 pl-24 pr-20 shadow-md transition hover:opacity-80 active:shadow-none disabled:bg-disabled"
        >
          <div className="text-16 font-bold">解答する</div>
          <MaterialSymbol icon="send" size={24} />
        </Button>
      </div>
      <Suspense>
        <ConfirmModal
          isOpen={isModalOpen}
          formId={formId}
          confirmType="submit"
          confirmName="intent"
          confirmValue="submit"
          onCancel={() => setIsModalOpen(false)}
          title="解答の確認"
          confirmText="送信する"
          cancelText="キャンセル"
          dialogClassName="w-full max-w-[1024px] transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all break-words"
        >
          <div className="my-12]">
            <p className="mb-24 text-16 text-text">
              本当に解答を送信しますか？
            </p>
            <Typography>
              <Markdown>{lastResult?.answer}</Markdown>
            </Typography>
          </div>
        </ConfirmModal>
      </Suspense>
    </form>
  );
}

function AnswerTextInputField(props: {
  disabled?: boolean;
  invalid?: boolean;
  errorID?: string;
  defaultValue?: string;
  storageKey?: string;
}) {
  const storageValue =
    props.storageKey != null
      ? localStorage.getItem(props.storageKey + "/answer")
      : null;
  return (
    <Field disabled={props.disabled} className="flex flex-1">
      <Label className="sr-only">解答(必須)</Label>
      <Textarea
        name="answer"
        className="flex-1 resize-none rounded-12 border border-text p-12 data-[disabled]:cursor-not-allowed data-[disabled]:bg-disabled/45"
        placeholder="お世話になっております、チーム◯◯◯です。"
        defaultValue={props.defaultValue ?? storageValue ?? undefined}
        onChange={(e) => {
          if (props.storageKey != null) {
            localStorage.setItem(
              props.storageKey + "/answer",
              e.currentTarget.value,
            );
          }
        }}
        required
        invalid={props.invalid}
        aria-describedby={props.invalid ? props.errorID : undefined}
      />
    </Field>
  );
}
