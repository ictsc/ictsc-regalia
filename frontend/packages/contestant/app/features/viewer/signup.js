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
exports.signUp = signUp;
function signUp(request, baseURL) {
    return __awaiter(this, void 0, void 0, function () {
        var result, url, resp, body, codes, _i, codes_1, code;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    result = {};
                    if (request.invitationCode === "") {
                        result.invitationCodeError = "required";
                        result.error = "invalid";
                    }
                    if (request.name === "") {
                        result.nameError = "required";
                        result.error = "invalid";
                    }
                    if (request.displayName === "") {
                        result.displayNameError = "required";
                        result.error = "invalid";
                    }
                    if (result.error != null) {
                        return [2 /*return*/, result];
                    }
                    url = new URL(baseURL !== null && baseURL !== void 0 ? baseURL : window.location.href);
                    url.search = "";
                    url.pathname = "/api/auth/signup";
                    return [4 /*yield*/, fetch(url, {
                            method: "POST",
                            headers: {
                                "Content-Type": "application/json",
                            },
                            body: JSON.stringify({
                                invitation_code: request.invitationCode,
                                name: request.name,
                                display_name: request.displayName,
                            }),
                        })];
                case 1:
                    resp = _a.sent();
                    if (resp.status === 429) {
                        return [2 /*return*/, { error: "rate_limit" }];
                    }
                    if (resp.status === 400) {
                        result.error = "invalid";
                    }
                    return [4 /*yield*/, resp.json()];
                case 2:
                    body = _a.sent();
                    if (typeof body !== "object" || body === null) {
                        throw new Error("unexpected response");
                    }
                    if ("codes" in body) {
                        codes = body.codes;
                    }
                    else {
                        codes = [];
                    }
                    if (typeof codes !== "object" || codes === null || !Array.isArray(codes)) {
                        throw new Error("unexpected response");
                    }
                    for (_i = 0, codes_1 = codes; _i < codes_1.length; _i++) {
                        code = codes_1[_i];
                        switch (code) {
                            case "invalid_invitation_code":
                                result.invitationCodeError = "invalid";
                                break;
                            case "team_is_full":
                                result.invitationCodeError = "team_full";
                                break;
                            case "invalid_name":
                                result.nameError = "invalid";
                                break;
                            case "duplicate_name":
                                result.nameError = "duplicate";
                                break;
                            case "invalid_display_name":
                                result.displayNameError = "invalid";
                                break;
                            default:
                                break;
                        }
                    }
                    return [2 /*return*/, result];
            }
        });
    });
}
