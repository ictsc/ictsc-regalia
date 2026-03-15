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
exports.renderMarkdown = renderMarkdown;
var unified_1 = require("unified");
var remark_parse_1 = require("remark-parse");
var remark_gfm_1 = require("remark-gfm");
var remark_math_1 = require("remark-math");
var remark_rehype_1 = require("remark-rehype");
var remark_breaks_1 = require("remark-breaks");
var core_1 = require("@shikijs/rehype/core");
var rehype_react_1 = require("rehype-react");
var jsx_runtime_1 = require("react/jsx-runtime");
var core_2 = require("shiki/core");
var javascript_1 = require("shiki/engine/javascript");
var highlighterPromise = (0, core_2.createHighlighterCore)({
    themes: [Promise.resolve().then(function () { return require("@shikijs/themes/material-theme-lighter"); })],
    langs: [
        Promise.resolve().then(function () { return require("@shikijs/langs/diff"); }),
        Promise.resolve().then(function () { return require("@shikijs/langs/shellscript"); }),
        Promise.resolve().then(function () { return require("@shikijs/langs/shellsession"); }),
        Promise.resolve().then(function () { return require("@shikijs/langs/hcl"); }),
        Promise.resolve().then(function () { return require("@shikijs/langs/sql"); }),
    ],
    engine: (0, javascript_1.createJavaScriptRegexEngine)(),
});
function renderMarkdown(content) {
    return __awaiter(this, void 0, void 0, function () {
        var highlighter, file;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, highlighterPromise];
                case 1:
                    highlighter = _a.sent();
                    return [4 /*yield*/, (0, unified_1.unified)()
                            .use(remark_breaks_1.default)
                            .use(remark_parse_1.default, { fragment: true })
                            .use(remark_gfm_1.default)
                            .use(remark_math_1.default)
                            .use(remark_rehype_1.default)
                            .use(core_1.default, highlighter, {
                            theme: "material-theme-lighter",
                        })
                            .use(rehype_react_1.default, jsx_runtime_1.default)
                            .process(content)];
                case 2:
                    file = _a.sent();
                    return [2 /*return*/, file.result];
            }
        });
    });
}
