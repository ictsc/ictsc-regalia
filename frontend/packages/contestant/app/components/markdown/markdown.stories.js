"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LinksAndImages = exports.CodeAndTables = exports.TextFormatting = exports.ListsAndIndentation = exports.HeadingsAndParagraphs = void 0;
var markdown_1 = require("./markdown");
exports.default = {
    title: "components/Markdown",
    component: markdown_1.Markdown,
};
exports.HeadingsAndParagraphs = {
    render: function () { return (<markdown_1.Typography>
      <markdown_1.Markdown>
        {"\n# \u898B\u51FA\u30571\n## \u898B\u51FA\u30572\n### \u898B\u51FA\u30573\n#### \u898B\u51FA\u30574\n##### \u898B\u51FA\u30575\n###### \u898B\u51FA\u30576\n\n\u901A\u5E38\u306E\u6BB5\u843D\u30C6\u30AD\u30B9\u30C8\u3067\u3059\u3002\n\u8907\u6570\u884C\u306E\u6BB5\u843D\u3082\n\u3053\u306E\u3088\u3046\u306B\u66F8\u3051\u307E\u3059\u3002\n\n\u6BB5\u843D\u3068\u6BB5\u843D\u306E\u9593\u306F\n\n\u3053\u306E\u3088\u3046\u306B\u7A7A\u884C\u3092\u5165\u308C\u307E\u3059\u3002\n"}
      </markdown_1.Markdown>
    </markdown_1.Typography>); },
};
exports.ListsAndIndentation = {
    render: function () { return (<markdown_1.Typography>
      <markdown_1.Markdown>
        {"\n- \u7B87\u6761\u66F8\u304D1\u30EC\u30D9\u30EB\u76EE\n  - 2\u30EC\u30D9\u30EB\u76EE\n    - 3\u30EC\u30D9\u30EB\u76EE\n      - 4\u30EC\u30D9\u30EB\u76EE\n- \u5225\u306E\u9805\u76EE\n\n1. \u756A\u53F7\u4ED8\u304D\u30EA\u30B9\u30C8\n2. 2\u756A\u76EE\u306E\u9805\u76EE\n   1. \u30CD\u30B9\u30C8\u3057\u305F\u756A\u53F7\u4ED8\u304D\n   2. 2\u756A\u76EE\u306E\u30CD\u30B9\u30C8\n3. 3\u756A\u76EE\u306E\u9805\u76EE\n\n- [ ] \u30C1\u30A7\u30C3\u30AF\u30DC\u30C3\u30AF\u30B9\uFF08\u672A\u30C1\u30A7\u30C3\u30AF\uFF09\n- [x] \u30C1\u30A7\u30C3\u30AF\u30DC\u30C3\u30AF\u30B9\uFF08\u30C1\u30A7\u30C3\u30AF\u6E08\u307F\uFF09\n"}
      </markdown_1.Markdown>
    </markdown_1.Typography>); },
};
exports.TextFormatting = {
    render: function () { return (<markdown_1.Typography>
      <markdown_1.Markdown>
        {"\n*\u30A4\u30BF\u30EA\u30C3\u30AF* _\u30A4\u30BF\u30EA\u30C3\u30AF_\n**\u592A\u5B57** __\u592A\u5B57__\n***\u592A\u5B57\u30A4\u30BF\u30EA\u30C3\u30AF*** ___\u592A\u5B57\u30A4\u30BF\u30EA\u30C3\u30AF___\n~\u6253\u3061\u6D88\u3057\u7DDA~\n`inline code`\n\n---\n\n> \u5F15\u7528\u6587\n> \u8907\u6570\u884C\u306E\n> \u5F15\u7528\u6587\n\n> \u30CD\u30B9\u30C8\u3057\u305F\n>> \u5F15\u7528\u6587\n"}
      </markdown_1.Markdown>
    </markdown_1.Typography>); },
};
exports.CodeAndTables = {
    render: function () { return (<markdown_1.Typography>
      <markdown_1.Markdown>
        {"\n```hcl\nresource \"aws_s3_bucket\" \"bucket\" {\n  bucket = \"my-bucket\"\n}\n```\n\n| \u52171 | \u52172 | \u52173 |\n|-----|-----|-----|\n| A1  | B1  | C1  |\n| A2  | B2  | C2  |\n| A3  | B3  | C3  |\n"}
      </markdown_1.Markdown>
    </markdown_1.Typography>); },
};
exports.LinksAndImages = {
    render: function () { return (<markdown_1.Typography>
      <markdown_1.Markdown>
        {"\n[\u30EA\u30F3\u30AF\u30C6\u30AD\u30B9\u30C8](https://example.com)\n![\u753B\u50CF\u306E\u4EE3\u66FF\u30C6\u30AD\u30B9\u30C8](https://example.com/image.jpg)\n\n\u81EA\u52D5\u30EA\u30F3\u30AF: https://example.com\n"}
      </markdown_1.Markdown>
    </markdown_1.Typography>); },
};
