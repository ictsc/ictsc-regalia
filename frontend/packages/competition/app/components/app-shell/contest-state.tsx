import { MaterialSymbol } from "../material-symbol";

export type ContestState = "before" | "running" | "break" | "finished";

export function ContestStateView() {
  return (
    <div className="flex h-[50px] items-center rounded-[10px] bg-surface-1 px-[10px] text-text">
      <MaterialSymbol icon="schedule" size={24} />
      <span className="ml-[5px] text-16">競技開始前</span>
      <span className="ml-[20px]">
        <span className="text-12">残り</span>
        <span className="ml-[5px] text-24">00 : 00 : 00</span>
      </span>
    </div>
  );
}
