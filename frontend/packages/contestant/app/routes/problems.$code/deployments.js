"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Deployments = Deployments;
exports.EmptyDeploymentList = EmptyDeploymentList;
exports.DeploymentList = DeploymentList;
exports.DeploymentItem = DeploymentItem;
var react_1 = require("@headlessui/react");
var clsx_1 = require("clsx");
var react_2 = require("react");
var v1_1 = require("@ictsc/proto/contestant/v1");
var confirmModal_1 = require("./confirmModal");
function Deployments(props) {
    var buttonID = (0, react_2.useId)();
    var _a = (0, react_2.useState)(false), isModalOpen = _a[0], setIsModalOpen = _a[1];
    var handleRedeployClick = function () {
        setIsModalOpen(true);
    };
    var handleConfirm = function () {
        setIsModalOpen(false);
        props.redeploy();
    };
    var handleCancel = function () {
        setIsModalOpen(false);
    };
    return (<>
      <div className="flex size-full flex-col gap-16">
        <div className="rounded-12 bg-surface-1 size-full py-12">
          <div className="size-full overflow-y-auto px-12 [scrollbar-gutter:stable_both-edges]">
            {props.list}
          </div>
        </div>
        <div className="flex items-center justify-end gap-16">
          {props.error != null && (<label htmlFor={buttonID} className="text-14 text-primary">
              {props.error}
            </label>)}
          <react_1.Button id={buttonID} className={(0, clsx_1.clsx)("rounded-12 bg-surface-2 grid place-items-center px-24 py-16 shadow-md transition", "data-[disabled]:bg-disabled data-[active]:shadow-none data-[hover]:opacity-80")} disabled={!props.canRedeploy || props.isRedeploying} onClick={handleRedeployClick}>
            <span className="text-16 font-bold">再展開する</span>
          </react_1.Button>
        </div>
      </div>
      <confirmModal_1.ConfirmModal isOpen={isModalOpen} onConfirm={handleConfirm} onCancel={handleCancel} title="再展開の確認" confirmText="再展開する" cancelText="キャンセル" dialogClassName="w-full max-w-md transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all">
        <div className="my-12">
          <p className="text-16 text-text">本当にこの問題を再展開しますか？</p>
          <p className="text-16 text-text">
            残り許容回数:
            <span className={(0, clsx_1.clsx)("text-16 pl-4 font-bold", props.allowedDeploymentCount <= 0
            ? "text-primary"
            : "text-text")}>
              {props.allowedDeploymentCount}
            </span>
          </p>
        </div>
      </confirmModal_1.ConfirmModal>
    </>);
}
function EmptyDeploymentList(props) {
    return (<div className="text-16 text-text grid size-full place-items-center">
      <p className="flex flex-col items-center gap-8">
        <h1 className="font-bold">まだ再展開されていません</h1>
        <h2>
          許容回数:
          <span className="ms-4 font-bold">{props.allowedDeploymentCount}</span>
        </h2>
      </p>
    </div>);
}
function DeploymentList(props) {
    return (<ul className={(0, clsx_1.clsx)("flex size-full flex-col gap-16", props.isPending && "opacity-75")}>
      {props.children}
    </ul>);
}
var deploymentDateTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
    dateStyle: "short",
    timeStyle: "short",
});
function DeploymentItem(props) {
    var _a;
    var status = props.status;
    return (<li className={(0, clsx_1.clsx)("rounded-12 flex justify-between gap-8 p-16", status !== v1_1.DeploymentStatus.DEPLOYED ? "bg-surface-0" : "bg-disabled", props.isPending && "opacity-75")}>
      <div className="flex flex-col">
        <h2 className="text-20 font-bold">
          {deploymentDateTimeFormatter.format(new Date(props.requestedAt))}
        </h2>
        <h3 className={(0, clsx_1.clsx)("text-12", status !== v1_1.DeploymentStatus.DEPLOYED && "text-primary")}>
          {(_a = {},
            _a[v1_1.DeploymentStatus.UNSPECIFIED] = null,
            _a[v1_1.DeploymentStatus.DEPLOYED] = "展開完了",
            _a[v1_1.DeploymentStatus.DEPLOYING] = "展開中",
            _a[v1_1.DeploymentStatus.FAILED] = "展開失敗",
            _a)[status]}
        </h3>
      </div>
      <div className="flex flex-col">
        <p className="border-text flex items-baseline justify-end gap-4 border-b pb-4 font-bold">
          <span className="text-24">{props.revision}</span>
          <span className="text-14">回目</span>
        </p>
        <div className="text-14 grid grid-cols-[repeat(2,minmax(24px,auto))] grid-rows-2 place-items-end gap-4 font-bold">
          <p>残り許容回数</p>
          <p className={(0, clsx_1.clsx)(props.thresholdExceeded && "text-primary")}>
            {props.allowedDeploymentCount}
          </p>
          <p>総減点</p>
          <p className={(0, clsx_1.clsx)(props.thresholdExceeded && "text-primary")}>
            {props.penalty}
          </p>
        </div>
      </div>
    </li>);
}
