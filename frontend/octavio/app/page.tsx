"use client";

import { rule } from "@/components/_const";
import ICTSCCard from "@/components/card";
import MarkdownPreview from "@/components/markdown-preview";
import ICTSCTitle from "@/components/title";

interface SelectedData {
  group: number;
  values: number[];
}

function Home() {
  const button = () => {
    const checkedRadioElements = document.querySelectorAll(
      'input[type="radio"]:checked',
    );
    const checkedCheckboxElements = document.querySelectorAll(
      'input[type="checkbox"]:checked',
    );

    const processData = (element: HTMLInputElement): SelectedData => {
      const ids = element.id.split("-");
      return {
        group: parseInt(ids[1], 10),
        values: [parseInt(ids[2], 10)],
      };
    };

    const radioData = Array.from(checkedRadioElements).map((e) =>
      processData(e as HTMLInputElement),
    );

    const checkboxDataMap: Map<number, number[]> = new Map();
    Array.from(checkedCheckboxElements).forEach((element) => {
      const data = processData(element as HTMLInputElement);
      const existingValues = checkboxDataMap.get(data.group) || [];
      checkboxDataMap.set(data.group, [...existingValues, ...data.values]);
    });

    const checkboxData = Array.from(checkboxDataMap.entries()).map(
      ([group, values]) => ({ group, values }),
    );

    // radioDataとcheckboxDataを統合
    const selectedDataArray: SelectedData[] = [...radioData, ...checkboxData];

    console.log(selectedDataArray);
  };

  return (
    <>
      <ICTSCTitle title="ルール" />
      <main className="container-ictsc">
        <ICTSCCard>
          <MarkdownPreview
            content={`
# ルール

問題文１

\`\`\`ictscr
- [ ] 問題１
- [ ] 問題２
- [ ] 問題３
- [ ] 問題４
\`\`\`

問題文２

\`\`\`ictscr
- [ ] 問題１
- [ ] 問題２
- [ ] 問題３
- [ ] 問題４
\`\`\`

問題文３

\`\`\`ictscc
- [ ] 問題１
- [ ] 問題２
- [ ] 問題３
- [ ] 問題４
\`\`\`
          
          `}
          />
          <input
            type="button"
            className={"btn"}
            onClick={button}
            value="ボタン"
          />
        </ICTSCCard>
      </main>
    </>
  );
}

export default Home;
