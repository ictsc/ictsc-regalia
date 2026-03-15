"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DAY = exports.HOUR = exports.MINUTE = void 0;
exports.ContestStateView = ContestStateView;
var material_symbol_1 = require("@app/components/material-symbol");
var stateMap = {
    before: "競技開始前",
    running: "競技中",
    break: "休憩中",
    finished: "競技終了",
};
exports.MINUTE = 60;
exports.HOUR = 60 * exports.MINUTE;
exports.DAY = 24 * exports.HOUR;
function ContestStateView(_a) {
    var state = _a.state, restDurationSeconds = _a.restDurationSeconds;
    var rest = restDurationSeconds;
    var days = Math.floor(rest / exports.DAY);
    rest %= exports.DAY;
    var hours = Math.floor(rest / exports.HOUR);
    rest %= exports.HOUR;
    var minutes = Math.floor(rest / exports.MINUTE);
    var seconds = rest % exports.MINUTE;
    return (<div className="bg-surface-1 text-text flex h-[48px] w-[288px] items-center justify-between rounded-[8px] px-[8px]">
      <div className="flex items-center">
        <material_symbol_1.MaterialSymbol icon="schedule" size={24}/>
        <span className="text-16 ml-[4px] line-clamp-1 w-[80px] text-clip">
          {stateMap[state]}
        </span>
      </div>
      {restDurationSeconds !== 0 && (<div className="flex items-baseline">
          <span className="text-12">残り</span>
          <span className="w-[128px] text-end">
            {days > 0 ? (<>
                <span className="text-24">{days}</span>
                <span className="text-16">日</span>
              </>) : (<span className="text-24">
                {hours} : {minutes} : {seconds}
              </span>)}
          </span>
        </div>)}
    </div>);
}
