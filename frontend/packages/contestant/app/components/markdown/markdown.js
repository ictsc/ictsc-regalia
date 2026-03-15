"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Typography = Typography;
exports.Markdown = Markdown;
var react_1 = require("react");
var clsx_1 = require("clsx");
var markdown_module_css_1 = require("./markdown.module.css");
var markdown_utils_1 = require("./markdown-utils");
function Typography(props) {
    return (<div className={(0, clsx_1.clsx)(markdown_module_css_1.default.content, props.className)}>
      {props.children}
    </div>);
}
function Markdown(_a) {
    var children = _a.children;
    var nodePromise = (0, react_1.useMemo)(function () { return (0, markdown_utils_1.renderMarkdown)(children !== null && children !== void 0 ? children : ""); }, [children]);
    return (0, react_1.use)(nodePromise);
}
