import { clsx } from "clsx";

export function DeploymentDetail(props: {
  readonly revision?: number; // 何回目の再展開か
  readonly remainingRedeploys?: number; // 残り許容回数
  readonly exceededRedeployLimit?: boolean; // 残り許容回数を超えているか
  readonly totalPenalty?: number; // 総減点
}) {
  return (
    <div className="flex flex-col">
      <p className="flex place-content-end items-baseline justify-end gap-4 border-b border-text pb-4 pl-8 *:inline-block">
        <span
          className={clsx(
            "text-24 font-bold",
            props.revision == null && "px-16",
          )}
        >
          {props.revision != null ? props.revision : "-"}
        </span>
        <span className="text-14 font-bold">回目</span>
      </p>
      <div className="grid grid-cols-[repeat(2,auto)] grid-rows-2 place-content-end gap-4 text-14 font-bold">
        <p className="place-self-end">残り許容回数:</p>
        <p
          className={clsx(
            "place-self-end",
            props.exceededRedeployLimit && "text-primary",
            props.remainingRedeploys == null && "px-8",
          )}
        >
          {props.remainingRedeploys != null ? props.remainingRedeploys : "-"}
        </p>
        <p className="place-self-end">総減点:</p>
        <p
          className={clsx(
            "place-self-end",
            props.totalPenalty == null && "px-8",
            props.exceededRedeployLimit && "text-primary",
          )}
        >
          {props.totalPenalty != null ? props.totalPenalty : "-"}
        </p>
      </div>
    </div>
  );
}
