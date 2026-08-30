import { useEffect, useRef, type ComponentProps } from "react";
import { clsx } from "clsx";
import { Link } from "@tanstack/react-router";
import { Title } from "../components/title";
import { Score } from "../components/score";

type ActivityEntry = {
  problemCode: string;
  problemTitle: string;
  answerId: number;
  submittedAt: string;
  score: ComponentProps<typeof Score>;
  scored: boolean;
  onDownload: () => void;
};

type ActivityPageProps = {
  entries: ActivityEntry[];
};

const formatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "short",
  timeStyle: "medium",
});

export function ActivityPage(props: ActivityPageProps) {
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const listEl = listRef.current;
    if (!listEl) return;

    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;

    // Find nearest scrollable ancestor
    let scrollContainer: HTMLElement | null = listEl.parentElement;
    while (scrollContainer) {
      const style = getComputedStyle(scrollContainer);
      if (
        style.overflowY === "auto" ||
        style.overflowY === "scroll" ||
        style.overflow === "auto" ||
        style.overflow === "scroll"
      ) {
        break;
      }
      scrollContainer = scrollContainer.parentElement;
    }
    const scrollTarget: HTMLElement | Window = scrollContainer ?? window;
    const containerEl = scrollContainer ?? document.documentElement;

    function update() {
      if (!listEl) return;
      const containerRect = containerEl.getBoundingClientRect();
      const containerCenter = containerRect.top + containerEl.clientHeight / 2;
      const halfHeight = containerEl.clientHeight / 2;

      const rows = listEl.children;
      for (let i = 0; i < rows.length; i++) {
        const row = rows[i] as HTMLElement;
        const rowRect = row.getBoundingClientRect();
        const itemCenter = rowRect.top + row.offsetHeight / 2;
        const distance = Math.abs(itemCenter - containerCenter);
        const t = Math.min(distance / halfHeight, 1);
        const scale = 1 - t * 0.5;
        const opacity = 1 - t * 0.7;
        // 縦線(child[0])はスキップし、時刻・矢印・カードだけにエフェクトを適用
        for (let j = 1; j < row.children.length; j++) {
          const child = row.children[j] as HTMLElement;
          child.style.transform = `scale(${scale})`;
          // opacity-50 クラス（未採点カード）の場合は元の 0.5 を掛け合わせる
          const baseOpacity = child.classList.contains("opacity-50") ? 0.5 : 1;
          child.style.opacity = `${opacity * baseOpacity}`;
        }
      }
    }

    function updatePadding() {
      if (!listEl || listEl.children.length === 0) return;
      const half = containerEl.clientHeight / 2;
      const firstRow = listEl.children[0] as HTMLElement;
      const lastRow = listEl.children[
        listEl.children.length - 1
      ] as HTMLElement;
      listEl.style.paddingTop = `${Math.max(0, half - firstRow.offsetHeight / 2)}px`;
      listEl.style.paddingBottom = `${Math.max(0, half - lastRow.offsetHeight / 2)}px`;
    }

    updatePadding();
    update();

    scrollTarget.addEventListener("scroll", update, { passive: true });
    const observer = new ResizeObserver(() => {
      updatePadding();
      update();
    });
    observer.observe(containerEl);

    return () => {
      scrollTarget.removeEventListener("scroll", update);
      observer.disconnect();
      if (listEl) {
        listEl.style.paddingTop = "";
        listEl.style.paddingBottom = "";
      }
    };
  }, [props.entries]);

  return (
    <>
      <Title>アクティビティ</Title>
      <div className="my-40 flex size-full flex-col items-center px-20">
        {props.entries.length === 0 ? (
          <p className="text-16 text-text mt-64">提出履歴がありません</p>
        ) : (
          <div ref={listRef} className="flex w-full max-w-screen-md flex-col">
            {props.entries.map((entry, index) => (
              <div
                key={`${entry.problemCode}-${entry.answerId}`}
                className="flex flex-row items-center"
              >
                {/* 縦線+丸 */}
                <div className="flex w-12 shrink-0 flex-col items-center self-stretch">
                  {index === 0 ? (
                    <div
                      className="w-2 grow"
                      style={{
                        backgroundImage:
                          "linear-gradient(to bottom, transparent, var(--color-text) 100%)",
                        opacity: 0.1,
                        maskImage:
                          "repeating-linear-gradient(to bottom, black 0 3px, transparent 3px 6px)",
                      }}
                    />
                  ) : (
                    <div className="bg-text/10 w-2 grow" />
                  )}
                  <div
                    className={clsx(
                      "size-12 shrink-0 rounded-full",
                      entry.scored ? "bg-primary" : "bg-text/20",
                    )}
                  />
                  <div
                    className={clsx(
                      "w-2 grow",
                      index < props.entries.length - 1
                        ? "bg-text/10"
                        : "bg-transparent",
                    )}
                  />
                </div>
                {/* 時刻 */}
                <div
                  className={clsx(
                    "text-12 ml-12 w-[150px] shrink-0 whitespace-nowrap",
                    !entry.scored && "text-text/40",
                  )}
                >
                  {formatter.format(new Date(entry.submittedAt))}
                </div>
                {/* 矢印 */}
                <div
                  className={clsx(
                    "text-16 shrink-0 px-8",
                    entry.scored ? "text-primary" : "text-text/20",
                  )}
                >
                  ←
                </div>
                {/* カード */}
                <div
                  className={clsx(
                    "rounded-16 my-12 flex min-w-0 flex-1 flex-col gap-12 p-20 shadow-lg sm:flex-row sm:items-center sm:justify-between",
                    !entry.scored && "opacity-50",
                  )}
                >
                  <div className="flex min-w-0 flex-1 flex-col gap-4">
                    <Link
                      to="/problems/$code"
                      params={{ code: entry.problemCode }}
                      className="text-14 truncate font-bold hover:underline"
                    >
                      {entry.problemCode}: {entry.problemTitle}
                    </Link>
                    <div className="flex items-center gap-8">
                      <p className="text-12">提出 #{entry.answerId}</p>
                      <a href="#" className="text-8" onClick={entry.onDownload}>
                        ダウンロード
                      </a>
                    </div>
                  </div>
                  {entry.scored ? (
                    <Score {...entry.score} />
                  ) : (
                    <p className="text-14 font-bold">採点中</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
