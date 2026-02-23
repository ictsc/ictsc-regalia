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
  const [seconds, setSeconds] = useState(0);

  useEffect(() => {
    if (lastSubmittedAt == null || submitInterval == null) {
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

      setSeconds(remainingSeconds);
    };

    checkAnswerable();
    const interval = setInterval(checkAnswerable, 1000);
    return () => clearInterval(interval);
  }, [lastSubmittedAt, submitInterval]);

  const remainingSeconds =
    lastSubmittedAt != null && submitInterval != null ? seconds : 0;

  return { remainingSeconds };
}

export function SubmissionForm(props: {
  readonly action: (answer: string) => Promise<"success" | "failure">;
  readonly submitInterval?: number;
  readonly lastSubmittedAt?: string;
  readonly storageKey?: string;
  readonly submissionStatus?: {
    isSubmittable: boolean;
    submittableUntil?: Date;
  };
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
  const isRateLimitOk = remainingSeconds <= 0;
  const isScheduleOk = props.submissionStatus?.isSubmittable ?? true;
  const isAnswerable = isRateLimitOk && isScheduleOk;
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
        {!isRateLimitOk && (
          <div className="flex w-[160px] items-center justify-between">
            <span className="text-16 text-black">解答可能まで</span>
            <div className="flex items-center">
              <span className="text-20 text-primary w-[24px] text-right font-bold">
                {minutes}
              </span>
              <span className="text-20 text-primary mx-2 font-bold">:</span>
              <span className="text-20 text-primary w-[24px] text-right font-bold">
                {seconds.toString().padStart(2, "0")}
              </span>
            </div>
          </div>
        )}
        {isRateLimitOk && !isScheduleOk && (
          <div className="text-16 text-text">
            {props.submissionStatus?.submittableUntil ? (
              <span>
                次回提出可能時刻まで提出できません
              </span>
            ) : (
              <span>この問題は現在提出できません</span>
            )}
          </div>
        )}
        {isAnswerable && lastResult?.type === "error" && (
          <label id={errorId} className="text-16 text-primary shrink font-bold">
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
          className="rounded-12 bg-surface-2 disabled:bg-disabled flex items-center justify-center self-end py-16 pr-20 pl-24 shadow-md transition hover:opacity-80 active:shadow-none"
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
          dialogClassName="w-full max-w-[1024px] transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all"
        >
          <div className="my-12]">
            <p className="text-16 text-text mb-24">
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
        className="rounded-12 border-text data-[disabled]:bg-disabled/45 flex-1 resize-none border p-12 data-[disabled]:cursor-not-allowed"
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
