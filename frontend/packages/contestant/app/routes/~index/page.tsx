import { Logo } from "@app/components/logo";
import { MaterialSymbol } from "@app/components/material-symbol";

export type IndexPageProps = {
  readonly inContest: boolean;
};
export function IndexPage({ inContest }: IndexPageProps) {
  return inContest ? <InContest /> : <OutOfContest />;
}

function InContest() {
  return (
    <div className="flex h-full flex-col items-center justify-center">
      <Logo width={500} />
      <span className="mt-16 text-16 underline">
        左のサイドメニューからタブを選択してください
      </span>
      <div className="mt-[48px] flex flex-col gap-8 rounded-16 border-2 border-primary p-16 *:px-8">
        <div className="flex">
          <MaterialSymbol icon="schedule" size={40} className="text-icon" />
          <div className="ml-8 flex flex-col">
            <div className="text-24 leading-[40px]">競技中</div>
            <div className="flex items-baseline">
              <div className="text-14">残り</div>
              <div className="w-[168px] text-end text-32">01 : 23 : 45</div>
            </div>
          </div>
        </div>
        <div className="flex w-full items-center border-t border-primary">
          <div className="flex size-40 items-center justify-center">
            <MaterialSymbol
              icon="arrow_forward_ios"
              size={24}
              className="text-icon"
            />
          </div>
          <div className="ml-8 text-14">次のフェーズ: 競技終了</div>
        </div>
      </div>
    </div>
  );
}

function OutOfContest() {
  return (
    <div className="flex h-full flex-col items-center justify-center">
      <h1 className="text-48 font-bold underline">競技開始まであと</h1>
      <div className="mt-40 flex items-center">
        <MaterialSymbol icon="schedule" size={48} className="text-icon" />
        <span className="ml-16 text-48 font-bold">03 : 30 : 50</span>
      </div>
    </div>
  );
}
