"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_router_1 = require("@tanstack/react-router");
var problem_1 = require("../features/problem");
var answer_1 = require("../features/answer");
var deployment_1 = require("../features/deployment");
var announce_1 = require("../features/announce");
exports.Route = (0, react_router_1.createFileRoute)("/problems/$code")({
    loader: function (_a) {
        var transport = _a.context.transport, code = _a.params.code;
        var fetchAnswersResult = (0, answer_1.fetchAnswers)(transport, code);
        var answers = fetchAnswersResult.then(function (r) { return r.answers; });
        var metadata = fetchAnswersResult.then(function (r) { return r.metadata; });
        var deployments = (0, deployment_1.fetchDeployments)(transport, code);
        return {
            problem: (0, problem_1.fetchProblem)(transport, code),
            notices: (0, announce_1.fetchNotices)(transport),
            answers: answers,
            metadata: metadata,
            submitAnswer: function (body) { return (0, answer_1.submitAnswer)(transport, code, body); },
            deployments: deployments,
            deploy: function () { return (0, deployment_1.deploy)(transport, code); },
            fetchAnswer: function (num) { return (0, answer_1.fetchAnswer)(transport, code, num); },
        };
    },
});
