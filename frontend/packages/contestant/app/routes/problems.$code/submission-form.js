"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubmissionForm = SubmissionForm;
var react_1 = require("react");
var react_2 = require("@headlessui/react");
var material_symbol_1 = require("../../components/material-symbol");
var confirmModal_1 = require("./confirmModal");
var markdown_1 = require("@app/components/markdown");
function formatRemainingTimeParts(remainingSeconds) {
    var minutes = Math.floor(remainingSeconds / 60);
    var seconds = remainingSeconds % 60;
    return { minutes: minutes, seconds: seconds };
}
function useAnswerable(lastSubmittedAt, submitInterval) {
    var _a = (0, react_1.useState)(0), seconds = _a[0], setSeconds = _a[1];
    (0, react_1.useEffect)(function () {
        if (lastSubmittedAt == null || submitInterval == null) {
            return;
        }
        var checkAnswerable = function () {
            var now = new Date();
            var lastSubmit = new Date(lastSubmittedAt);
            var nextSubmitTime = new Date(lastSubmit.getTime() + submitInterval * 1000);
            var diffMs = nextSubmitTime.getTime() - now.getTime();
            var diffSec = Math.ceil(diffMs / 1000);
            var remainingSeconds = diffSec > 0 ? diffSec : 0;
            setSeconds(remainingSeconds);
        };
        checkAnswerable();
        var interval = setInterval(checkAnswerable, 1000);
        return function () { return clearInterval(interval); };
    }, [lastSubmittedAt, submitInterval]);
    var remainingSeconds = lastSubmittedAt != null && submitInterval != null ? seconds : 0;
    return { remainingSeconds: remainingSeconds };
}
function SubmissionForm(props) {
    var _this = this;
    var _a, _b, _c;
    var _d = (0, react_1.useState)(false), isModalOpen = _d[0], setIsModalOpen = _d[1];
    var _e = (0, react_1.useActionState)(function (_prevState, formData) { return __awaiter(_this, void 0, void 0, function () {
        var answer, intent, _a, result, err_1;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    answer = formData.get("answer");
                    if (answer.trim() === "") {
                        return [2 /*return*/, {
                                type: "error",
                                name: "answer",
                                error: "解答を入力してください",
                            }];
                    }
                    intent = formData.get("intent");
                    _a = intent;
                    switch (_a) {
                        case "confirm": return [3 /*break*/, 1];
                        case "submit": return [3 /*break*/, 2];
                    }
                    return [3 /*break*/, 6];
                case 1:
                    setIsModalOpen(true);
                    return [2 /*return*/, { type: "confirm", answer: answer }];
                case 2:
                    _b.trys.push([2, 4, , 5]);
                    return [4 /*yield*/, props.action(answer)];
                case 3:
                    result = _b.sent();
                    if (result === "failure") {
                        return [2 /*return*/, {
                                type: "error",
                                error: "解答の送信に失敗しました",
                            }];
                    }
                    return [3 /*break*/, 5];
                case 4:
                    err_1 = _b.sent();
                    console.error(err_1);
                    return [2 /*return*/, {
                            type: "error",
                            error: "解答の送信に失敗しました",
                        }];
                case 5: return [2 /*return*/];
                case 6: return [2 /*return*/];
            }
        });
    }); }, null), lastResult = _e[0], action = _e[1], isPending = _e[2];
    var remainingSeconds = useAnswerable(props.lastSubmittedAt, props.submitInterval).remainingSeconds;
    var isRateLimitOk = remainingSeconds <= 0;
    var isScheduleOk = (_b = (_a = props.submissionStatus) === null || _a === void 0 ? void 0 : _a.isSubmittable) !== null && _b !== void 0 ? _b : true;
    var isAnswerable = isRateLimitOk && isScheduleOk;
    var _f = formatRemainingTimeParts(remainingSeconds), minutes = _f.minutes, seconds = _f.seconds;
    var formId = (0, react_1.useId)();
    var errorId = (0, react_1.useId)();
    return (<form id={formId} className="flex size-full flex-col" noValidate action={action}>
      <AnswerTextInputField invalid={(lastResult === null || lastResult === void 0 ? void 0 : lastResult.type) === "error" && lastResult.name === "answer"} defaultValue={(lastResult === null || lastResult === void 0 ? void 0 : lastResult.type) === "confirm" ? lastResult.answer : undefined} storageKey={props.storageKey}/>
      <div className="mt-20 flex items-center justify-end gap-24">
        {!isRateLimitOk && (<div className="flex w-[160px] items-center justify-between">
            <span className="text-16 text-black">解答可能まで</span>
            <div className="flex items-center">
              <span className="text-20 text-primary w-[24px] text-right font-bold">
                {minutes}
              </span>
              <span className="text-20 text-primary mx-2 font-bold">:</span>
              <span className="text-20 text-primary w-[24px] text-right font-bold">
                {seconds.toString().padStart(2, "0")}
              </span>
            </div>
          </div>)}
        {isRateLimitOk && !isScheduleOk && (<div className="text-16 text-text">
            {((_c = props.submissionStatus) === null || _c === void 0 ? void 0 : _c.submittableUntil) ? (<span>次回提出可能時刻まで提出できません</span>) : (<span>この問題は現在提出できません</span>)}
          </div>)}
        {isAnswerable && (lastResult === null || lastResult === void 0 ? void 0 : lastResult.type) === "error" && (<label id={errorId} className="text-16 text-primary shrink font-bold">
            {lastResult.error}
          </label>)}
        <react_2.Button name="intent" value="confirm" type="submit" disabled={!isAnswerable || ((lastResult === null || lastResult === void 0 ? void 0 : lastResult.type) === "confirm" && isPending)} className="rounded-12 bg-surface-2 disabled:bg-disabled flex items-center justify-center self-end py-16 pr-20 pl-24 shadow-md transition hover:opacity-80 active:shadow-none">
          <div className="text-16 font-bold">解答する</div>
          <material_symbol_1.MaterialSymbol icon="send" size={24}/>
        </react_2.Button>
      </div>
      <react_1.Suspense>
        <confirmModal_1.ConfirmModal isOpen={isModalOpen} formId={formId} confirmType="submit" confirmName="intent" confirmValue="submit" onCancel={function () { return setIsModalOpen(false); }} title="解答の確認" confirmText="送信する" cancelText="キャンセル" dialogClassName="w-full max-w-[1024px] transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all">
          <div className="my-12]">
            <p className="text-16 text-text mb-24">
              本当に解答を送信しますか？
            </p>
            <markdown_1.Typography>
              <markdown_1.Markdown>{lastResult === null || lastResult === void 0 ? void 0 : lastResult.answer}</markdown_1.Markdown>
            </markdown_1.Typography>
          </div>
        </confirmModal_1.ConfirmModal>
      </react_1.Suspense>
    </form>);
}
function AnswerTextInputField(props) {
    var _a, _b;
    var storageValue = props.storageKey != null
        ? localStorage.getItem(props.storageKey + "/answer")
        : null;
    return (<react_2.Field disabled={props.disabled} className="flex flex-1">
      <react_2.Label className="sr-only">解答(必須)</react_2.Label>
      <react_2.Textarea name="answer" className="rounded-12 border-text data-[disabled]:bg-disabled/45 flex-1 resize-none border p-12 data-[disabled]:cursor-not-allowed" placeholder="お世話になっております、チーム◯◯◯です。" defaultValue={(_b = (_a = props.defaultValue) !== null && _a !== void 0 ? _a : storageValue) !== null && _b !== void 0 ? _b : undefined} onChange={function (e) {
            if (props.storageKey != null) {
                localStorage.setItem(props.storageKey + "/answer", e.currentTarget.value);
            }
        }} required invalid={props.invalid} aria-describedby={props.invalid ? props.errorID : undefined}/>
    </react_2.Field>);
}
