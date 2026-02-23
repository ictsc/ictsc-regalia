import { useState } from "react";
import { clsx } from "clsx";
import type { TeamProfile } from "@ictsc/proto/contestant/v1";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";

type TeamListProps = {
  teamProfile: TeamProfile[];
};

export function TeamsPage(props: TeamListProps) {
  const [openStates, setOpenStates] = useState<{ [key: number]: boolean }>({});

  const toggleAccordion = (index: number) => {
    setOpenStates({
      ...openStates,
      [index]: !openStates[index],
    });
  };

  return (
    <>
      <Title>チーム一覧</Title>
      <div className="pt-64">
        {props.teamProfile.map((team, index) => {
          return (
            <div
              key={index}
              className="flex items-center justify-center gap-x-40 pb-64 pl-8 md:flex-nowrap"
            >
              <div className="rounded-16 flex w-[90%] max-w-[650px] min-w-[300px] flex-row gap-16 px-20 py-24 shadow-lg md:w-[650px]">
                {/* アコーディオンボタン */}
                <div>
                  <button
                    className="h-[110px] md:h-64"
                    onClick={() => toggleAccordion(index)}
                  >
                    <MaterialSymbol
                      icon="arrow_forward_ios"
                      size={24}
                      className={clsx(
                        "transition-transform",
                        openStates[index] && "rotate-90",
                      )}
                    />
                  </button>
                </div>
                {/* アコーディオンボタン以外の表示要素 */}
                <div className="flex flex-col">
                  <div className="ml-20 flex w-[200px] min-w-[150px] flex-col gap-20 md:ml-40 md:w-[500px] md:flex-row md:gap-40">
                    {/* チーム名表示 */}
                    <div className="flex flex-[1] flex-col overflow-hidden">
                      <p className="text-14 md:pb-20">チーム名</p>
                      <p
                        className="text-16 w-full truncate overflow-hidden font-bold whitespace-nowrap"
                        title={team.name}
                      >
                        {team.name}
                      </p>
                    </div>
                    {/* 所属表示 */}
                    <div className="flex flex-[1] flex-col overflow-hidden">
                      <p className="text-14 md:pb-20">所属</p>
                      <p
                        className="text- w-full truncate overflow-hidden font-bold whitespace-nowrap"
                        title={team.organization}
                      >
                        {team.organization}
                      </p>
                    </div>
                  </div>
                  {/* 名前表示 */}
                  {openStates[index] && (
                    <div className="mt-40 ml-20 md:ml-40">
                      <p className="pb-20">名前</p>
                      {team.members.map((member, index) => (
                        <p key={index} className="mb-8 font-bold">
                          {member.displayName}
                        </p>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </>
  );
}
