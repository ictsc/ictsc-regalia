import {
  useReducer,
  useId,
  type ActionDispatch,
  useEffect,
  useState,
} from "react";
import { Field, Label, Textarea } from "@headlessui/react";
import { MaterialSymbol } from "../../components/material-symbol";

export function SubmissionForm(props: {
  readonly action: (answer: string) => Promise<"success" | "failure">;
  readonly submitInterval?: number;
  readonly lastSubmittedAt?: string;
}) {
  const [error, dispatchError] = useReducer<FormErrorState, [FormErrorAction]>(
    reduceFormError,
    null,
  );
  const [isAnswerable, setIsAnswerable] = useState(true);

  useEffect(() => {
    if (!props.lastSubmittedAt || !props.submitInterval) {
      setIsAnswerable(true);
      return;
    }

    const checkAnswerable = () => {
      const now = new Date();
      const lastSubmit = new Date(props.lastSubmittedAt!);
      const nextSubmitTime = new Date(
        lastSubmit.getTime() + props.submitInterval! * 1000,
      );
      setIsAnswerable(now >= nextSubmitTime);
    };
    checkAnswerable();
    const interval = setInterval(checkAnswerable, 1000);
    return () => clearInterval(interval);
  }, [props.lastSubmittedAt, props.submitInterval]);

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
}: {
  error: FormErrorState;
  dispatchError: ActionDispatch<[FormErrorAction]>;
  isAnswerable: boolean;
}) {
  const submitLabelID = useId();
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
      <div className="mt-20 flex items-center justify-end gap-16">
        {error != null && (
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
