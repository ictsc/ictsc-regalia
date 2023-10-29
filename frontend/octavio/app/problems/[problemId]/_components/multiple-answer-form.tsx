import { useEffect, useState } from "react";

import {
  errorNotify,
  successNotify,
} from "@/app/problems/[problemId]/_components/notify";
import useAnswers from "@/hooks/answer";
import useApi from "@/hooks/api";
import useAuth from "@/hooks/auth";
import useProblem from "@/hooks/problem";

interface SelectionData {
  group: number;
  value: number[];
  type: "radio" | "checkbox";
}

function MultipleAnswerForm({ code }: { code: string }) {
  const { client } = useApi();
  const { user } = useAuth();
  const { problem } = useProblem(code);
  const { answers, mutate } = useAnswers(problem?.id ?? null);
  const sortedAnswers = answers.sort((a, b) => {
    if (a.created_at < b.created_at) {
      return 1;
    }
    if (a.created_at > b.created_at) {
      return -1;
    }
    return 0;
  });
  const [latestAnswerID, setLatestAnswerID] = useState<null | string>(null);

  // 初期選択状態をラジオボタンとチェックボックスに反映
  useEffect(() => {
    if (sortedAnswers.length <= 0) {
      return;
    }
    const answer = sortedAnswers[0];

    // 最後の回答が違っていたらアラートを出す
    if (latestAnswerID !== null && latestAnswerID !== answer.id) {
      // eslint-disable-next-line no-alert
      const isConfirmed = window.confirm(
        "回答の更新があります。フォームを更新しますか？",
      );

      setLatestAnswerID(answer.id);

      if (!isConfirmed) {
        return;
      }
    }
    setLatestAnswerID(answer.id);

    const initialSelections: SelectionData[] = JSON.parse(answer.body).map(
      (selection: SelectionData) => ({
        ...selection,
        value: selection.value,
      }),
    );

    initialSelections.forEach(({ group, value, type }) => {
      value?.forEach((val) => {
        const elementId = `${type}-${group}-${val}`;
        const inputElement = document.getElementById(
          elementId,
        ) as HTMLInputElement;
        if (inputElement) {
          inputElement.checked = true;
        }
      });
    });
  }, [latestAnswerID, sortedAnswers]);

  const sendButton = async () => {
    const checkedRadioElements = document.querySelectorAll(
      'input[type="radio"]:checked',
    );
    const checkedCheckboxElements = document.querySelectorAll(
      'input[type="checkbox"]:checked',
    );

    const processData = (
      element: HTMLInputElement,
      type: "radio" | "checkbox",
    ): SelectionData => {
      const ids = element.id.split("-");
      return {
        group: parseInt(ids[1], 10),
        value: [parseInt(ids[2], 10)],
        type,
      };
    };

    const radioData = Array.from(checkedRadioElements).map((e) =>
      processData(e as HTMLInputElement, "radio"),
    );

    const checkboxDataMap: Map<number, number[]> = new Map();
    Array.from(checkedCheckboxElements).forEach((element) => {
      const data = processData(element as HTMLInputElement, "checkbox");
      const existingValues = checkboxDataMap.get(data.group) || [];
      checkboxDataMap.set(data.group, [...existingValues, ...data.value]);
    });

    const checkboxData = Array.from(checkboxDataMap.entries()).map(
      ([group, value]) => ({ group, value, type: "checkbox" }),
    );

    // radioDataとcheckboxDataを統合
    const selectedDataArray = [...radioData, ...checkboxData];

    const response = await client.post(`problems/${problem?.id}/answers`, {
      user_group_id: user?.user_group_id,
      problem_id: problem?.id,
      body: JSON.stringify(selectedDataArray),
    });

    if (response.code === 201) {
      setLatestAnswerID(response.data.answer.id);
      successNotify();

      await mutate();
    } else {
      errorNotify();
    }
  };

  const clearButton = () => {
    const checkedRadioElements = document.querySelectorAll<HTMLInputElement>(
      'input[type="radio"]:checked',
    );
    const checkedCheckboxElements = document.querySelectorAll<HTMLInputElement>(
      'input[type="checkbox"]:checked',
    );

    checkedRadioElements.forEach((radioElem) => {
      // eslint-disable-next-line no-param-reassign
      radioElem.checked = false;
    });

    // checkbox の選択をクリア
    checkedCheckboxElements.forEach((checkboxElem) => {
      // eslint-disable-next-line no-param-reassign
      checkboxElem.checked = false;
    });
  };

  return (
    <div className="flex flex-row justify-between">
      <input
        type="button"
        className="btn btn-primary mt-4"
        onClick={sendButton}
        value="提出確認"
      />
      <input
        type="button"
        className="btn btn-ghost mt-4 text-primary"
        onClick={clearButton}
        value="フォームをクリア"
      />
    </div>
  );
}

export default MultipleAnswerForm;
