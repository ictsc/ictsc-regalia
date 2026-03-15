"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NotSubmittable = exports.Default = void 0;
var protobuf_1 = require("@bufbuild/protobuf");
var v1_1 = require("@ictsc/proto/contestant/v1");
var problem_item_1 = require("./problem-item");
exports.default = {
    title: "pages/problems/ProblemItem",
};
exports.Default = {
    render: function () { return (<div className="grid grid-cols-2 gap-64">
      <problem_item_1.ProblemItem code="ABC" title="あいしーてぃーえすしーだよあああああああ" score={{
            maxScore: 200,
            score: 200,
            rawScore: 200,
            penalty: 0,
            fullScore: true,
            rawFullScore: true,
        }}/>
      <problem_item_1.ProblemItem code="ABC" title="あいしーてぃーえすしーだよあああああああ" score={{
            maxScore: 200,
            score: 160,
            rawScore: 200,
            penalty: -40,
            fullScore: false,
            rawFullScore: true,
        }}/>
      <problem_item_1.ProblemItem code="ABC" title="あいしーてぃーえすしーだよあああああああ" score={{
            maxScore: 200,
            score: 100,
            rawScore: 120,
            penalty: -20,
            fullScore: false,
            rawFullScore: false,
        }}/>
      <problem_item_1.ProblemItem code="ABC" title="あいしーてぃーえすしーだよあああああああ" score={{
            maxScore: 200,
        }}/>
    </div>); },
};
exports.NotSubmittable = {
    render: function () { return (<div className="grid grid-cols-2 gap-64">
      <problem_item_1.ProblemItem code="ABC" title="提出不可（スコアあり）" score={{
            maxScore: 200,
            score: 100,
            rawScore: 120,
            penalty: -20,
            fullScore: false,
            rawFullScore: false,
        }} submissionStatus={(0, protobuf_1.create)(v1_1.SubmissionStatusSchema, {
            isSubmittable: false,
        })}/>
      <problem_item_1.ProblemItem code="ABC" title="提出不可（未回答）" score={{
            maxScore: 200,
        }} submissionStatus={(0, protobuf_1.create)(v1_1.SubmissionStatusSchema, {
            isSubmittable: false,
        })}/>
    </div>); },
};
