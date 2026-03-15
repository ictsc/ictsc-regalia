"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Empty = exports.Default = void 0;
var actions_1 = require("storybook/actions");
var page_1 = require("./page");
exports.default = {
    title: "pages/problem",
    render: function (props) { return (<div style={{
            "--header-height": "0",
            "--content-height": "100vh",
            "--content-width": "100%",
        }}>
      <page_1.Page {...props}/>
    </div>); },
    args: {
        redeployable: false,
        content: (<page_1.Content code="AAA" title="Title" body={"\n## \u6982\u8981\n\n\u3042\u308C\u3053\u308C\n\n## \u524D\u63D0\u6761\u4EF6\n\n- \u3044\u308D\u3044\u308D\n- \u3042\u308B\u3088\u306D\n\n## \u521D\u671F\u72B6\u614B\n\n- \u3053\u3046\u306A\u3063\u3066\u308B\n\n## \u7D42\u4E86\u72B6\u614B\n\n- \u3053\u3046\u306A\u308B\n\n## \u63A5\u7D9A\u60C5\u5831\n\n| \u30DB\u30B9\u30C8\u540D | IP\u30A2\u30C9\u30EC\u30B9 | \u30E6\u30FC\u30B6\u540D | \u30D1\u30B9\u30EF\u30FC\u30C9 |\n| --- | --- | --- | --- |\n| Web | 192.168.0.1 | user | password |\n"}/>),
        submissionForm: (<page_1.SubmissionForm action={function () {
                submitAction("submit");
                return Promise.resolve("success");
            }}/>),
        submissionList: (<page_1.SubmissionList>
        {Array.from({ length: 10 }).map(function (_, i) { return (<page_1.SubmissionListItem key={i} id={i + 1} submittedAt="2025-02-03T00:00:00Z" score={{ maxScore: 100 }} downloadAnswer={function () { }}/>); })}
      </page_1.SubmissionList>),
        deploymentList: null,
    },
};
var submitAction = (0, actions_1.action)("submit");
exports.Default = {};
exports.Empty = {
    args: {
        submissionList: <page_1.EmptySubmissionList />,
    },
};
