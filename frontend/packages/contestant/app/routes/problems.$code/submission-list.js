"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubmissionListContainer = SubmissionListContainer;
exports.SubmissionList = SubmissionList;
exports.EmptySubmissionList = EmptySubmissionList;
exports.SubmissionListItem = SubmissionListItem;
var clsx_1 = require("clsx");
var score_1 = require("../../components/score");
function SubmissionListContainer(props) {
    return (<div className="rounded-12 bg-surface-1 size-full py-12">
      <div className="size-full overflow-y-auto px-12 [scrollbar-gutter:stable_both-edges]">
        {props.children}
      </div>
    </div>);
}
function SubmissionList(props) {
    return (<ul className={(0, clsx_1.clsx)("flex size-full flex-col gap-16 py-12", props.isPending && "opacity-75")}>
      {props.children}
    </ul>);
}
function EmptySubmissionList() {
    return (<div className="text-16 text-text grid size-full place-items-center font-bold">
      解答はまだありません！
    </div>);
}
var submissionListDateTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
    dateStyle: "medium",
    timeStyle: "short",
});
function SubmissionListItem(props) {
    return (<li className="rounded-12 bg-surface-0 flex justify-between gap-8 p-16">
      <div className="flex flex-col justify-between">
        <div>
          <h2 className="text-20 font-bold text-[#000]">#{props.id}</h2>
          <h3 className="text-12">
            提出:{" "}
            {submissionListDateTimeFormatter.format(new Date(props.submittedAt))}
          </h3>
        </div>
        <a href="#" className="text-8" onClick={props.downloadAnswer}>
          ダウンロード
        </a>
      </div>
      <score_1.Score {...props.score}/>
    </li>);
}
