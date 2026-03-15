"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Score = Score;
var clsx_1 = require("clsx");
function Score(props) {
    return (<div className="flex flex-col">
      <p className="border-text flex flex-row items-baseline gap-4 border-b pb-4 pl-8 *:inline-block">
        <span className={(0, clsx_1.clsx)("text-24 font-bold", props.fullScore && "text-primary", props.score == null && "px-16")}>
          {props.score != null ? props.score : "-"}
        </span>
        <span className="text-20 -translate-y-2">/</span>
        <span className="text-14 font-bold">{props.maxScore}</span>
      </p>
      <div className="text-14 grid grid-cols-[repeat(2,auto)] grid-rows-2 place-content-end gap-4 font-bold">
        <p>素点:</p>
        <p className={(0, clsx_1.clsx)("place-self-end", props.rawFullScore && "text-primary", props.rawScore == null && "px-8")}>
          {props.rawScore != null ? props.rawScore : "-"}
        </p>
        <p>減点:</p>
        <p className={(0, clsx_1.clsx)("place-self-end", props.penalty == null && "px-8")}>
          {props.penalty != null ? props.penalty : "-"}
        </p>
      </div>
    </div>);
}
