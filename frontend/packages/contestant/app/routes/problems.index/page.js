"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ProblemsPage = ProblemsPage;
var clsx_1 = require("clsx");
var react_1 = require("react");
var react_2 = require("@headlessui/react");
var problem_item_1 = require("./problem-item");
var unread_announces_banner_1 = require("./unread-announces-banner");
var score_1 = require("../../features/score");
var title_1 = require("../../components/title");
var material_symbol_1 = require("../../components/material-symbol");
var group_1 = require("../../features/problem/group");
function ScheduleLabel(props) {
    if (props.schedules.length === 0) {
        return <span className="opacity-50">スケジュール未設定</span>;
    }
    return props.schedules.map(function (s, i) { return (<react_1.Fragment key={s.name}>
      {i > 0 && " / "}
      <span className={(0, clsx_1.clsx)(s.temporalStatus === "past" && "opacity-50", s.temporalStatus === "current" && "text-primary")}>
        {s.name}
      </span>
    </react_1.Fragment>); });
}
function ProblemsPage(props) {
    var groups = (0, group_1.groupProblems)(props.problems);
    return (<>
      <title_1.Title>問題一覧</title_1.Title>
      <div className="mx-16 my-64 flex justify-center">
        <div className="flex flex-col gap-16">
          <div className="ml-16">
            <unread_announces_banner_1.UnreadAnnouncesBanner notices={props.notices}/>
          </div>
          {groups.map(function (group) { return (<react_2.Disclosure key={group.key} as="section" defaultOpen>
              <react_2.DisclosureButton className="group/disc flex w-full cursor-pointer items-center gap-16">
                <material_symbol_1.MaterialSymbol icon="arrow_forward_ios" size={24} className="text-disabled transition-transform group-data-[open]/disc:rotate-90"/>
                <div className="border-disabled flex-1 border-t"/>
                <h2 className="text-24 shrink-0 text-center font-bold">
                  <ScheduleLabel schedules={group.schedules}/>
                </h2>
                <div className="border-disabled flex-1 border-t"/>
              </react_2.DisclosureButton>
              <react_2.DisclosurePanel static className="h-0 overflow-hidden data-[open]:mt-16 data-[open]:h-auto data-[open]:overflow-visible">
                <ul className="grid grid-flow-row grid-cols-1 gap-x-40 gap-y-24 lg:grid-cols-2">
                  {group.problems.map(function (problem) { return (<li key={problem.code}>
                      <problem_item_1.ProblemItem code={problem.code} title={problem.title} score={(0, score_1.protoScoreToProps)(problem.maxScore, problem.score)} submissionStatus={problem.submissionStatus}/>
                    </li>); })}
                </ul>
              </react_2.DisclosurePanel>
            </react_2.Disclosure>); })}
        </div>
      </div>
    </>);
}
