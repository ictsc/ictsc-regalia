"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubmitAnswer = exports.SubmitDeployment = void 0;
var confirmModal_1 = require("./confirmModal");
var markdown_1 = require("@app/components/markdown");
exports.default = {
    title: "pages/problem/confirmModal",
};
exports.SubmitDeployment = {
    render: function () { return (<confirmModal_1.ConfirmModal isOpen={true} onConfirm={function () { }} onCancel={function () { }} title="再展開の確認" confirmText="再展開する" cancelText="キャンセル" dialogClassName="w-full max-w-md transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all">
      <span>ここに本文を書く</span>
    </confirmModal_1.ConfirmModal>); },
};
exports.SubmitAnswer = {
    render: function () { return (<confirmModal_1.ConfirmModal isOpen={true} onConfirm={function () { }} onCancel={function () { }} title="解答の確認" confirmText="送信する" cancelText="キャンセル" dialogClassName="w-full max-w-[1024px] transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all">
      <div style={{ padding: "16px" }}>
        <markdown_1.Typography>
          <markdown_1.Markdown>{markdownContent}</markdown_1.Markdown>
        </markdown_1.Typography>
      </div>
    </confirmModal_1.ConfirmModal>); },
};
var markdownContent = "\n# \u898B\u51FA\u30571\n## \u898B\u51FA\u30572\n### \u898B\u51FA\u30573\n#### \u898B\u51FA\u30574\n##### \u898B\u51FA\u30575\n###### \u898B\u51FA\u30576\n\n```shell\necho hoge\npwd\nls\n```\n\n| \u52171 | \u52172 | \u52173 |\n|-----|-----|-----|\n| A1  | B1  | C1  |\n| A2  | B2  | C2  |\n| A3  | B3  | C3  |\n";
