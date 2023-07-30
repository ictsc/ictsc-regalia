"use client";

import "../styles/globals.css";
import React from "react";

import localFont from "@next/font/local";

import { RecoilRoot } from "recoil";

import { site } from "@/components/_const";
import {
  DismissNoticeStateInit,
  WatchDismissNoticeIds,
} from "@/hooks/state/recoil";

const notoSansJP = localFont({
  variable: "--font-noto-sans-jp",
  src: "dist/fonts/NotoSansJP-VF.woff2",
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <head>
        <title>{site}</title>
        <link rel="icon" type="image/png" sizes="32x32" href="/favicon.png" />
      </head>
      <body>
        <RecoilRoot>
          <div data-theme="ictsc" className={notoSansJP.className}>
            {children}
          </div>

          {/* 通知の永続化まわり */}
          <DismissNoticeStateInit />
          <WatchDismissNoticeIds />
        </RecoilRoot>
      </body>
    </html>
  );
}
