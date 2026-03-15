"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ActivityPage = ActivityPage;
var clsx_1 = require("clsx");
var react_router_1 = require("@tanstack/react-router");
var title_1 = require("../components/title");
var score_1 = require("../components/score");
var formatter = new Intl.DateTimeFormat("ja-JP", {
    dateStyle: "short",
    timeStyle: "medium",
});
function ActivityPage(props) {
    return (<>
      <title_1.Title>アクティビティ</title_1.Title>
      <div className="my-40 flex size-full flex-col items-center px-20">
        {props.entries.length === 0 ? (<p className="text-16 text-text mt-64">提出履歴がありません</p>) : (<div className="flex w-full max-w-screen-md flex-col">
            {props.entries.map(function (entry, index) { return (<div key={"".concat(entry.problemCode, "-").concat(entry.answerId)} className="flex flex-row items-center">
                {/* 縦線+丸 */}
                <div className="flex w-12 shrink-0 flex-col items-center self-stretch">
                  {index === 0 ? (<div className="w-2 grow" style={{
                        backgroundImage: "linear-gradient(to bottom, transparent, var(--color-text) 100%)",
                        opacity: 0.1,
                        maskImage: "repeating-linear-gradient(to bottom, black 0 3px, transparent 3px 6px)",
                    }}/>) : (<div className="bg-text/10 w-2 grow"/>)}
                  <div className={(0, clsx_1.clsx)("size-12 shrink-0 rounded-full", entry.scored ? "bg-primary" : "bg-text/20")}/>
                  <div className={(0, clsx_1.clsx)("w-2 grow", index < props.entries.length - 1
                    ? "bg-text/10"
                    : "bg-transparent")}/>
                </div>
                {/* 時刻 */}
                <div className={(0, clsx_1.clsx)("text-12 ml-12 w-[150px] shrink-0 whitespace-nowrap", !entry.scored && "text-text/40")}>
                  {formatter.format(new Date(entry.submittedAt))}
                </div>
                {/* 矢印 */}
                <div className={(0, clsx_1.clsx)("text-16 shrink-0 px-8", entry.scored ? "text-primary" : "text-text/20")}>
                  ←
                </div>
                {/* カード */}
                <div className={(0, clsx_1.clsx)("rounded-16 my-12 flex min-w-0 flex-1 flex-col gap-12 p-20 shadow-lg sm:flex-row sm:items-center sm:justify-between", !entry.scored && "opacity-50")}>
                  <div className="flex min-w-0 flex-1 flex-col gap-4">
                    <react_router_1.Link to="/problems/$code" params={{ code: entry.problemCode }} className="text-14 truncate font-bold hover:underline">
                      {entry.problemCode}: {entry.problemTitle}
                    </react_router_1.Link>
                    <div className="flex items-center gap-8">
                      <p className="text-12">提出 #{entry.answerId}</p>
                      <a href="#" className="text-8" onClick={entry.onDownload}>
                        ダウンロード
                      </a>
                    </div>
                  </div>
                  {entry.scored ? (<score_1.Score {...entry.score}/>) : (<p className="text-14 font-bold">採点中</p>)}
                </div>
              </div>); })}
          </div>)}
      </div>
    </>);
}
