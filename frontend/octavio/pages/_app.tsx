import "../styles/globals.css";
import localFont from "@next/font/local";
import type { AppProps } from "next/app";

import { RecoilRoot } from "recoil";

import {
  DismissNoticeStateInit,
  WatchDismissNoticeIds,
} from "@/hooks/state/recoil";

const notoSansJP = localFont({
  variable: "--font-noto-sans-jp",
  src: "./NotoSansJP-VF.woff2",
});

export default function App({ Component, pageProps }: AppProps) {
  return (
    <RecoilRoot>
      <div data-theme="ictsc" className={`${notoSansJP.className}`}>
        <Component {...pageProps} />
      </div>
      {/* 通知の永続化まわり */}
      <DismissNoticeStateInit />
      <WatchDismissNoticeIds />
    </RecoilRoot>
  );
}
