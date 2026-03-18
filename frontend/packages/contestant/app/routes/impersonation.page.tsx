import { type ReactNode } from "react";
import { Link } from "@tanstack/react-router";
import { Title } from "../components/title";

type Props = {
  children: ReactNode;
};

export function ImpersonationPage({ children }: Props) {
  return (
    <>
      <Title>管理者ログイン</Title>
      <div className="bg-surface-1 flex min-h-full flex-col items-center justify-center px-40 py-24">
        <div className="border-disabled border-t-primary bg-surface-0 rounded-8 flex w-full max-w-[480px] flex-col gap-16 border border-t-4 p-24">
          <h1 className="text-icon text-20 font-bold">
            特定の競技者としてログイン
          </h1>
          {children}
          <Link
            to="/signin"
            className="text-14 text-text flex w-fit items-center gap-4 underline"
          >
            ← 戻る
          </Link>
        </div>
      </div>
    </>
  );
}
