"use client";

import React from "react";

import { RecoilRoot } from "recoil";

import {
  DismissNoticeStateInit,
  WatchDismissNoticeIds,
} from "@/hooks/state/recoil";

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <RecoilRoot>
      {children}

      {/* 通知の永続化まわり */}
      <DismissNoticeStateInit />
      <WatchDismissNoticeIds />
    </RecoilRoot>
  );
}
