import {
  useReducer,
  useId,
  type ActionDispatch,
  useEffect,
  useState,
} from "react";
import { Field, Label, Textarea } from "@headlessui/react";
import { MaterialSymbol } from "../../components/material-symbol";

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
  submitInterval?: number
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
      const nextSubmitTime = new Date(lastSubmit.getTime() + submitInterval * 1000);

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
}) {
  const [error, dispatchError] = useReducer<FormErrorState, [FormErrorAction]>(
    reduceFormError,
    null,
  );
  const { remainingSeconds } = useAnswerable(
    props.lastSubmittedAt,
    props.submitInterval
  );
  const isAnswerable = remainingSeconds <= 0;

  return (
    <form
      className="flex size-full flex-col"
      onReset={() => dispatchError("reset")}
      action={async (data) => {
        if (!isAnswerable) return;
        const answer = data.get("answer");
        if (typeof answer !== "string") {
          return;
        }
        const result = await props.action(answer);
        switch (result) {
          case "success":
            dispatchError("reset");
            break;
          case "failure":
            dispatchError("failure");
            break;
        }
      }}
    >
      <SubmissionFormInner
        error={error}
        dispatchError={dispatchError}
        isAnswerable={isAnswerable}
        remainingSeconds={remainingSeconds}
      />
    </form>
  );
}

type FormErrorAction = "reset" | "missing-answer" | "failure";
type FormErrorState = string | null;

function reduceFormError(
  _: FormErrorState,
  action: FormErrorAction,
): FormErrorState {
  switch (action) {
    case "reset":
      return null;
    case "missing-answer":
      return "回答を入力してください";
    case "failure":
      return "回答の送信に失敗しました";
  }
}

function SubmissionFormInner({
  error,
  dispatchError,
  isAnswerable,
  remainingSeconds,
}: {
  error: FormErrorState;
  dispatchError: ActionDispatch<[FormErrorAction]>;
  isAnswerable: boolean;
  remainingSeconds: number;
}) {
  const submitLabelID = useId();
  const { minutes, seconds } = formatRemainingTimeParts(remainingSeconds);

  return (
    <>
      <Field className="flex flex-1">
        <Label className="sr-only">回答(必須)</Label>
        <Textarea
          name="answer"
          className="flex-1 resize-none rounded-12 border border-text p-12"
          required
          placeholder="お世話になっております、チーム◯◯◯です。"
          onInvalid={(e) => {
            e.preventDefault();
            if (!(e.target instanceof HTMLTextAreaElement)) {
              return;
            }
            e.target.focus();
            const validity = e.target.validity;
            if (validity.valueMissing) {
              dispatchError("missing-answer");
            } else {
              dispatchError("failure");
            }
          }}
        />
      </Field>

      <div className="mt-20 flex items-center justify-end gap-24">
          {!isAnswerable && (
            <label
              id={submitLabelID}
              className="flex items-center justify-between w-[160px]"
            >
              <span className="text-16 text-black">解答可能まで</span>
              <div className="flex items-center">
                <span className="text-primary font-bold text-20 text-right w-[24px]">
                  {minutes}
                </span>
                <span className="text-primary font-bold text-20 mx-2">:</span>
                <span className="text-primary font-bold text-20 text-right w-[24px]">
                  {seconds.toString().padStart(2, "0")}
                </span>
              </div>
            </label>
          )}

          {/* 回答可かつエラーがある場合 */}
          {isAnswerable && error != null && (
            <label
              id={submitLabelID}
              className="flex-shrink text-16 font-bold text-primary"
            >
              {error}
            </label>
          )}

        <button
          aria-labelledby={submitLabelID}
          type="submit"
          disabled={!isAnswerable}
          className="flex items-center justify-center self-end rounded-12 bg-surface-2 py-16 pl-24 pr-20 shadow-md transition hover:opacity-80 active:shadow-none disabled:bg-disabled"
        >
          <div className="text-16 font-bold">回答する</div>
          <MaterialSymbol icon="send" size={24} />
        </button>
      </div>
    </>
  );
}
