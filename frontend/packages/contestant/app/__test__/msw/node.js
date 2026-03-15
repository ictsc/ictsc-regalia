"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.setupMSW = setupMSW;
var vitest_1 = require("vitest");
var node_1 = require("msw/node");
function setupMSW(handlers, options) {
    if (handlers === void 0) { handlers = []; }
    if (options === void 0) { options = { onUnhandledRequest: "error" }; }
    var server = node_1.setupServer.apply(void 0, handlers);
    (0, vitest_1.beforeAll)(function () { return server.listen(options); });
    (0, vitest_1.afterAll)(function () { return server.close(); });
    (0, vitest_1.afterEach)(function () { return server.resetHandlers(); });
    return server;
}
