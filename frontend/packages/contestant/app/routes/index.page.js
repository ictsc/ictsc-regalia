"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.IndexPage = IndexPage;
exports.Timer = Timer;
var react_1 = require("react");
var date_fns_1 = require("date-fns");
var logo_1 = require("../components/logo");
var material_symbol_1 = require("../components/material-symbol");
var title_1 = require("../components/title");
function IndexPage(props) {
    switch (props.state) {
        case "in_contest":
            return <InContest {...props}/>;
        case "waiting":
            return <OutOfContest {...props}/>;
        case "ended":
            return <EndOfContest />;
        default:
            return null;
    }
}
function Timer(props) {
    var _a, _b, _c, _d;
    var _e = (0, react_1.useState)(function () { return Date.now(); }), nowState = _e[0], setNow = _e[1];
    (0, react_1.useEffect)(function () {
        var interval = setInterval(function () { return setNow(Date.now()); }, 1000);
        return function () { return clearInterval(interval); };
    }, []);
    var nowMs = (_a = props.nowMs) !== null && _a !== void 0 ? _a : nowState;
    var days = (0, date_fns_1.differenceInDays)(props.endMs, nowMs);
    var dur = (0, date_fns_1.intervalToDuration)({ start: nowMs, end: props.endMs });
    return (<>
      <title_1.Title />
      <p className="flex w-[5em] items-baseline justify-end" title={(0, date_fns_1.formatDuration)(dur)}>
        {days > 0 ? ("".concat(days, "\u65E5")) : (<>
            <span className="w-[1.5em] text-center">
              {"".concat((_b = dur.hours) !== null && _b !== void 0 ? _b : 0).padStart(2, "0")}
            </span>
            <span className="">:</span>
            <span className="w-[1.5em] text-center">
              {"".concat((_c = dur.minutes) !== null && _c !== void 0 ? _c : 0).padStart(2, "0")}
            </span>
            <span className="">:</span>
            <span className="w-[1.5em] text-center">
              {"".concat((_d = dur.seconds) !== null && _d !== void 0 ? _d : 0).padStart(2, "0")}
            </span>
          </>)}
      </p>
    </>);
}
function InContest(props) {
    return (<div className="mx-40 flex h-full flex-col items-center justify-center">
      <logo_1.Logo width={500}/>
      <span className="text-16 mt-16 underline">
        左のサイドメニューからタブを選択してください
      </span>
      <div className="rounded-16 border-primary mt-[48px] flex flex-col gap-8 border-2 p-16 *:px-8">
        <div className="flex">
          <material_symbol_1.MaterialSymbol icon="schedule" size={40} className="text-icon"/>
          <div className="ml-8 flex flex-col">
            <div className="text-24 leading-[40px]">
              競技中
              {props.currentScheduleName && " (".concat(props.currentScheduleName, ")")}
            </div>
            {props.timer != null && (<div className="flex items-baseline">
                <div className="text-14">残り</div>
                <div className="text-32 w-[168px] text-end">{props.timer}</div>
              </div>)}
          </div>
        </div>
        {props.nextScheduleName != null && (<div className="border-primary flex w-full items-center border-t pt-8">
            <div className="flex size-40 items-center justify-center">
              <material_symbol_1.MaterialSymbol icon="arrow_forward_ios" size={24} className="text-icon"/>
            </div>
            <div className="text-14 mt-2 ml-8">
              次のスケジュール: {props.nextScheduleName}
            </div>
          </div>)}
      </div>
    </div>);
}
function OutOfContest(props) {
    var title = props.nextScheduleName
        ? "".concat(props.nextScheduleName, " \u307E\u3067\u3042\u3068")
        : "次の競技まであと";
    return (<div className="mx-40 flex h-full flex-col items-center justify-center">
      <h1 className="text-48 font-bold underline">{title}</h1>
      <div className="mt-40 flex items-center">
        <material_symbol_1.MaterialSymbol icon="schedule" size={48} className="text-icon"/>
        <span className="text-48 ml-16 font-bold">{props.timer}</span>
      </div>
    </div>);
}
function EndOfContest() {
    return (<div className="mx-40 flex h-full flex-col items-center justify-center">
      <h1 className="text-48 font-bold underline">競技は終了しました</h1>
    </div>);
}
