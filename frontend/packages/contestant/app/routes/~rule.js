"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var rule_1 = require("@app/features/rule");
var react_router_1 = require("@tanstack/react-router");
exports.Route = (0, react_router_1.createFileRoute)("/rule")({
    loader: function (_a) {
        var transport = _a.context.transport;
        return ({
            rule: (0, rule_1.fetchRule)(transport),
        });
    },
});
