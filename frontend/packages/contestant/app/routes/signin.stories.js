"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Default = void 0;
var signin_page_1 = require("./signin.page");
exports.default = {
    title: "pages/signin",
    component: signin_page_1.SignInPage,
};
exports.Default = {
    name: "デフォルト",
    args: {
        signInURL: "/",
    },
};
