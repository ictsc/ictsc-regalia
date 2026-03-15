"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ProblemItem = ProblemItem;
var react_1 = require("react");
var clsx_1 = require("clsx");
var react_2 = require("@headlessui/react");
var react_router_1 = require("@tanstack/react-router");
var score_1 = require("../../components/score");
function ProblemItem(props) {
    var _a, _b;
    var isSubmittable = (_b = (_a = props.submissionStatus) === null || _a === void 0 ? void 0 : _a.isSubmittable) !== null && _b !== void 0 ? _b : true;
    return (<react_2.Button as={react_1.Fragment}>
      {function (_a) {
            var hover = _a.hover, active = _a.active, disabled = _a.disabled;
            return (<react_router_1.Link to="/problems/$code" params={{ code: props.code }} disabled={disabled} className={(0, clsx_1.clsx)("rounded-16 flex w-full max-w-[512px] justify-between gap-24 px-20 py-12 transition", active ? "shadow-transparent" : "shadow-lg", props.score.rawFullScore
                    ? hover
                        ? "bg-surface-2"
                        : "bg-disabled"
                    : hover
                        ? "bg-surface-1"
                        : "bg-surface-0", !isSubmittable && "opacity-50 grayscale")}>
          <div className="flex flex-col items-start justify-between gap-4">
            <div className="flex flex-col">
              <h3 className="text-24 text-primary font-bold">{props.code}</h3>
              <p className="text-16 line-clamp-1">{props.title}</p>
            </div>
          </div>
          <score_1.Score {...props.score}/>
        </react_router_1.Link>);
        }}
    </react_2.Button>);
}
