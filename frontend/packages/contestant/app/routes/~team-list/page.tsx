import { useState } from "react";

type TeamListProps = {
  teamName: string;
  affiliation: string;
  userNames: string[];
};

export function TeamList(props: TeamListProps) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="flex items-center justify-center gap-x-40 pb-64 pl-8 md:flex-nowrap">
      <div className="flex flex-row gap-16 px-20 py-24 w-[90%] min-w-[300px] max-w-[650px] shadow-lg md:w-[650px] rounded-16">
        {/* アコーディオンボタン */}
        <div>
          <button className="h-[110px] md:h-64"
            onClick={() => setIsOpen(!isOpen)}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className={`h-[25px] w-[25px] transition-transform ${isOpen ? 'rotate-90' : ''}`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
        {/* アコーディオンボタン以外の表示要素 */}
        <div className="flex flex-col">
          <div className="flex flex-col md:flex-row ml-20 md:ml-40 gap-20 md:gap-40 md:w-[500px] w-[200px] min-w-[150px]">
            {/* チーム名表示 */}
            <div className=" flex flex-[1] flex-col overflow-hidden">
              <p className="md:pb-20 text-14">チーム名</p>
              <p
                className="w-full overflow-hidden truncate whitespace-nowrap text-16 font-bold"
                title={props.teamName}
              >
                {props.teamName}
              </p>
            </div >
            {/* 所属表示 */}
            <div className="flex flex-[1] flex-col overflow-hidden">
              <p className="md:pb-20 text-14">所属</p>
              <p
                className="text- w-full overflow-hidden truncate whitespace-nowrap font-bold"
                title={props.affiliation}
              >
                {props.affiliation}
              </p>
            </div>
          </div>
          {/* 名前表示 */}
          {isOpen &&(
            <div className="ml-20 md:ml-40 mt-40">
              <p className="pb-20">名前</p>
              {props.userNames.map((name, index) => (
                <p key={index} className="mb-8 font-bold">{name}</p>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
