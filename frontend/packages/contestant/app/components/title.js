"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Title = Title;
function Title(props) {
    return (<title>
      {props.children != null ? "".concat(props.children, " | ICTSC2025") : "ICTSC2025"}
    </title>);
}
