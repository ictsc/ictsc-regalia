import { Button } from "@headlessui/react";
import { clsx } from "clsx";
import { useState, useId, type ReactNode } from "react";
import { type Deployment as DeploymentType } from "../../features/deployment";
import { DeploymentStatus } from "@ictsc/proto/contestant/v1";
import { ConfirmModal } from "./confirmModal";

export function Deployments(props: {
  list: ReactNode;
  canRedeploy: boolean;
  isRedeploying: boolean;
  allowedDeploymentCount: number;
  redeploy: () => void;
  error?: ReactNode;
}) {
  const buttonID = useId();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleRedeployClick = () => {
    setIsModalOpen(true);
  };

  const handleConfirm = () => {
    setIsModalOpen(false);
    props.redeploy();
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <div className="flex size-full flex-col gap-16">
        <div className="size-full rounded-12 bg-surface-1 py-12">
          <div className="size-full overflow-y-auto px-12 [scrollbar-gutter:stable_both-edges]">
            {props.list}
          </div>
        </div>
        <div className="flex items-center justify-end gap-16">
          {props.error != null && (
            <label htmlFor={buttonID} className="text-14 text-primary">
              {props.error}
            </label>
          )}
          <Button
            id={buttonID}
            className={clsx(
              "grid place-items-center rounded-12 bg-surface-2 px-24 py-16 shadow-md transition",
              "data-[disabled]:bg-disabled data-[hover]:opacity-80 data-[active]:shadow-none",
            )}
            disabled={!props.canRedeploy || props.isRedeploying}
            onClick={handleRedeployClick}
          >
            <span className="text-16 font-bold">再展開する</span>
          </Button>
        </div>
      </div>
      <ConfirmModal
        isOpen={isModalOpen}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
        title="再展開の確認"
        confirmText="再展開する"
        cancelText="キャンセル"
        dialogClassName="w-full max-w-md transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all"
      >
        <div className="my-12">
          <p className="text-16 text-text">本当にこの問題を再展開しますか？</p>
          <p className="text-16 text-text">
            残り許容回数:
            <span
              className={clsx(
                "pl-4 text-16 font-bold",
                props.allowedDeploymentCount <= 0
                  ? "text-primary"
                  : "text-text",
              )}
            >
              {props.allowedDeploymentCount}
            </span>
          </p>
        </div>
      </ConfirmModal>
    </>
  );
}

export function EmptyDeploymentList(props: { allowedDeploymentCount: number }) {
  return (
    <div className="grid size-full place-items-center text-16 text-text">
      <p className="flex flex-col items-center gap-8">
        <h1 className="font-bold">まだ再展開されていません</h1>
        <h2>
          許容回数:
          <span className="ms-4 font-bold">{props.allowedDeploymentCount}</span>
        </h2>
      </p>
    </div>
  );
}

export function DeploymentList(props: {
  isPending: boolean;
  children?: ReactNode;
}) {
  return (
    <ul
      className={clsx(
        "flex size-full flex-col gap-16",
        props.isPending && "opacity-75",
      )}
    >
      {props.children}
    </ul>
  );
}

const deploymentDateTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "short",
});

export function DeploymentItem(
  props: DeploymentType & { isPending?: boolean },
) {
  const { status } = props;
  return (
    <li
      className={clsx(
        "flex justify-between gap-8 rounded-12 p-16",
        status !== DeploymentStatus.DEPLOYED ? "bg-surface-0" : "bg-disabled",
        props.isPending && "opacity-75",
      )}
    >
      <div className="flex flex-col">
        <h2 className="text-20 font-bold">
          {deploymentDateTimeFormatter.format(new Date(props.requestedAt))}
        </h2>
        <h3
          className={clsx(
            "text-12",
            status !== DeploymentStatus.DEPLOYED && "text-primary",
          )}
        >
          {
            {
              [DeploymentStatus.UNSPECIFIED]: null,
              [DeploymentStatus.DEPLOYED]: "展開完了",
              [DeploymentStatus.DEPLOYING]: "展開中",
              [DeploymentStatus.FAILED]: "展開失敗",
            }[status]
          }
        </h3>
      </div>
      <div className="flex flex-col">
        <p className="flex items-baseline justify-end gap-4 border-b border-text pb-4 font-bold">
          <span className="text-24">{props.revision}</span>
          <span className="text-14">回目</span>
        </p>
        <div className="grid grid-cols-[repeat(2,minmax(24px,auto))] grid-rows-2 place-items-end gap-4 text-14 font-bold">
          <p>残り許容回数</p>
          <p className={clsx(props.thresholdExceeded && "text-primary")}>
            {props.allowedDeploymentCount}
          </p>
          <p>総減点</p>
          <p className={clsx(props.thresholdExceeded && "text-primary")}>
            {props.penalty}
          </p>
        </div>
      </div>
    </li>
  );
}
