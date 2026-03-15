"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.TeamsPage = TeamsPage;
var react_1 = require("react");
var clsx_1 = require("clsx");
var material_symbol_1 = require("../components/material-symbol");
var title_1 = require("../components/title");
function TeamsPage(props) {
    var _a = (0, react_1.useState)({}), openStates = _a[0], setOpenStates = _a[1];
    var toggleAccordion = function (index) {
        var _a;
        setOpenStates(__assign(__assign({}, openStates), (_a = {}, _a[index] = !openStates[index], _a)));
    };
    return (<>
      <title_1.Title>チーム一覧</title_1.Title>
      <div className="pt-64">
        {props.teamProfile.map(function (team, index) {
            return (<div key={index} className="flex items-center justify-center gap-x-40 pb-64 pl-8 md:flex-nowrap">
              <div className="rounded-16 flex w-[90%] max-w-[650px] min-w-[300px] flex-row gap-16 px-20 py-24 shadow-lg md:w-[650px]">
                {/* アコーディオンボタン */}
                <div>
                  <button className="h-[110px] md:h-64" onClick={function () { return toggleAccordion(index); }}>
                    <material_symbol_1.MaterialSymbol icon="arrow_forward_ios" size={24} className={(0, clsx_1.clsx)("transition-transform", openStates[index] && "rotate-90")}/>
                  </button>
                </div>
                {/* アコーディオンボタン以外の表示要素 */}
                <div className="flex flex-col">
                  <div className="ml-20 flex w-[200px] min-w-[150px] flex-col gap-20 md:ml-40 md:w-[500px] md:flex-row md:gap-40">
                    {/* チーム名表示 */}
                    <div className="flex flex-[1] flex-col overflow-hidden">
                      <p className="text-14 md:pb-20">チーム名</p>
                      <p className="text-16 w-full truncate overflow-hidden font-bold whitespace-nowrap" title={team.name}>
                        {team.name}
                      </p>
                    </div>
                    {/* 所属表示 */}
                    <div className="flex flex-[1] flex-col overflow-hidden">
                      <p className="text-14 md:pb-20">所属</p>
                      <p className="text- w-full truncate overflow-hidden font-bold whitespace-nowrap" title={team.organization}>
                        {team.organization}
                      </p>
                    </div>
                  </div>
                  {/* 名前表示 */}
                  {openStates[index] && (<div className="mt-40 ml-20 md:ml-40">
                      <p className="pb-20">名前</p>
                      {team.members.map(function (member, index) { return (<p key={index} className="mb-8 font-bold">
                          {member.displayName}
                        </p>); })}
                    </div>)}
                </div>
              </div>
            </div>);
        })}
      </div>
    </>);
}
