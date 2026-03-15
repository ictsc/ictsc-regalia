"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Default = void 0;
var logo_1 = require("./logo");
exports.default = {
    title: "components/Logo",
    component: logo_1.Logo,
};
exports.Default = {
    render: function () { return <logo_1.Logo height={100}/>; },
};
