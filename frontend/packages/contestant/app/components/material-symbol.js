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
exports.MaterialSymbol = MaterialSymbol;
var clsx_1 = require("clsx");
function MaterialSymbol(_a) {
    var icon = _a.icon, _b = _a.fill, fill = _b === void 0 ? false : _b, _c = _a.size, size = _c === void 0 ? 24 : _c, className = _a.className, propStyle = _a.style;
    var style = __assign(__assign({}, propStyle), { fontVariationSettings: "\"FILL\" ".concat(fill ? "1" : "0"), fontSize: size, width: size, height: size });
    return (<span 
    // material symbols の提供するクラスを使わないとフォントを指定できない
    className={(0, clsx_1.clsx)("material-symbols-outlined select-none", className)} style={style}>
      {icon}
    </span>);
}
