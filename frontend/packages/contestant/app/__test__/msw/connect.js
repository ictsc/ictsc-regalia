"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
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
var __classPrivateFieldSet = (this && this.__classPrivateFieldSet) || function (receiver, state, value, kind, f) {
    if (kind === "m") throw new TypeError("Private method is not writable");
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a setter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot write private member to an object whose class did not declare it");
    return (kind === "a" ? f.call(receiver, value) : f ? f.value = value : state.set(receiver, value)), value;
};
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var _ConnectHandler_instances, _ConnectHandler_router, _ConnectHandler_handler;
Object.defineProperty(exports, "__esModule", { value: true });
exports.connect = exports.ConnectHandler = void 0;
var msw_1 = require("msw");
var protocol_1 = require("@connectrpc/connect/protocol");
var connect_1 = require("@connectrpc/connect");
var defaultConnectResolver = function (_a) { return __awaiter(void 0, [_a], void 0, function (_b) {
    var ureq, ures;
    var request = _b.request, handler = _b.handler;
    return __generator(this, function (_c) {
        switch (_c.label) {
            case 0:
                ureq = (0, protocol_1.universalServerRequestFromFetch)(request.clone(), {});
                return [4 /*yield*/, handler(ureq)];
            case 1:
                ures = _c.sent();
                return [2 /*return*/, (0, protocol_1.universalServerResponseToFetch)(ures)];
        }
    });
}); };
var ConnectHandler = /** @class */ (function (_super) {
    __extends(ConnectHandler, _super);
    function ConnectHandler(info, routes, options) {
        var _this = _super.call(this, {
            info: __assign(__assign({}, info), { header: "".concat(info.kind, " ").concat(info.name) }),
            resolver: defaultConnectResolver,
            options: options,
        }) || this;
        _ConnectHandler_instances.add(_this);
        _ConnectHandler_router.set(_this, void 0);
        var router = (0, connect_1.createConnectRouter)();
        routes(router);
        __classPrivateFieldSet(_this, _ConnectHandler_router, router, "f");
        return _this;
    }
    ConnectHandler.prototype.parse = function (_a) {
        var request = _a.request;
        return Promise.resolve({
            handler: __classPrivateFieldGet(this, _ConnectHandler_instances, "m", _ConnectHandler_handler).call(this, request),
        });
    };
    ConnectHandler.prototype.predicate = function (_a) {
        var handler = _a.parsedResult.handler;
        if (handler == null) {
            return false;
        }
        return true;
    };
    ConnectHandler.prototype.extendResolverArgs = function (_a) {
        var handler = _a.parsedResult.handler;
        return {
            handler: handler,
        };
    };
    ConnectHandler.prototype.log = function (_a) {
        var request = _a.request, response = _a.response, handler = _a.parsedResult.handler;
        if (handler == null) {
            throw new Error("handler is null");
        }
        var method = handler.method;
        console.groupCollapsed("".concat(method.methodKind, " ").concat(method.parent.typeName, "/").concat(method.name, " (").concat(response.status, " ").concat(response.statusText, ")"));
        console.log("Request:", {
            url: new URL(request.url),
            method: request.method,
            headers: Object.fromEntries(request.headers.entries()),
        });
        console.log("Response:", {
            status: response.status,
            statusText: response.statusText,
            headers: Object.fromEntries(response.headers.entries()),
        });
        console.groupEnd();
    };
    return ConnectHandler;
}(msw_1.RequestHandler));
exports.ConnectHandler = ConnectHandler;
_ConnectHandler_router = new WeakMap(), _ConnectHandler_instances = new WeakSet(), _ConnectHandler_handler = function _ConnectHandler_handler(req) {
    var url = new URL(req.url);
    var handler = __classPrivateFieldGet(this, _ConnectHandler_router, "f").handlers.find(function (h) { return h.requestPath == url.pathname; });
    if (handler == null) {
        return;
    }
    if (!handler.allowedMethods.includes(req.method)) {
        return;
    }
    return handler;
};
exports.connect = {
    rpc: function (method, impl, options) {
        return new ConnectHandler({ kind: method.kind, name: "".concat(method.parent.typeName, "/").concat(method.name) }, function (router) { return router.rpc(method, impl); }, options);
    },
    service: function (service, impl, options) {
        return new ConnectHandler({ kind: "service", name: service.typeName }, function (router) { return router.service(service, impl); }, options);
    },
};
