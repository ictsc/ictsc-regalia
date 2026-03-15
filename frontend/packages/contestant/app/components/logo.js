"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Logo = Logo;
var ictsc2025_logo_svg_1 = require("@assets/ictsc2025_logo.svg");
var RATIO = 1960 / 524.3;
function Logo(_a) {
    var className = _a.className, width = _a.width, height = _a.height;
    if (width != null && height == null) {
        height = width / RATIO;
    }
    if (height != null && width == null) {
        width = height * RATIO;
    }
    return (<img className={className} width={width} height={height} src={ictsc2025_logo_svg_1.default} alt="ICTSC: ICT Trouble Shooting Contest"/>);
}
