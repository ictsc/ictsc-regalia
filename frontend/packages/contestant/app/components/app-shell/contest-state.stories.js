"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.All = void 0;
var contest_state_1 = require("./contest-state");
exports.default = {
    title: "ContentState",
    component: contest_state_1.ContestStateView,
};
exports.All = {
    render: function () { return (<div className="flex flex-col items-start gap-[20px]">
      <contest_state_1.ContestStateView state="break" restDurationSeconds={15 * contest_state_1.HOUR + 30 * contest_state_1.MINUTE + 50}/>
      <contest_state_1.ContestStateView state="running" restDurationSeconds={3 * contest_state_1.HOUR + 56 * contest_state_1.MINUTE + 24}/>
      <contest_state_1.ContestStateView state="before" restDurationSeconds={20 * contest_state_1.HOUR + 30 * contest_state_1.MINUTE + 50}/>
      <contest_state_1.ContestStateView state="before" restDurationSeconds={6 * contest_state_1.DAY + 50}/>
      <contest_state_1.ContestStateView state="finished" restDurationSeconds={0}/>
    </div>); },
};
