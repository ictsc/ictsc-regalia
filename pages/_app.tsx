import {useEffect} from "react";

import '../styles/globals.css'
import type {AppProps} from 'next/app'
import localFont from '@next/font/local'

const notoSansJP = localFont({
  variable: '--font-noto-sans-jp',
  src: './NotoSansJP-VF.woff2'
})

export default function App({Component, pageProps}: AppProps) {
  useEffect(() => {
    // 埋め込み構文の対応に必要 (1)
    import("zenn-embed-elements");

    // 埋め込み構文の対応に必要 (2)
    const getEmbeddedMessage = (n: any) => {
      try {
        const e = JSON.parse(n);
        if (typeof e.type != "string") throw new Error("Bad Request");
        if (typeof e.data != "object") throw new Error("Bad Request");
        return e
      } catch (e) {
        return {type: "none", data: {}}
      }
    };
    window.addEventListener("message", n => {
      const e = getEmbeddedMessage(n.data);
      switch (e.type) {
        case"ready": {
          const t = document.getElementById(e.data.id);
          if (t instanceof HTMLIFrameElement) {
            const a = {type: "rendering", data: {src: t.getAttribute("data-content")}};
            t.contentWindow && t.contentWindow.postMessage && t.contentWindow.postMessage(JSON.stringify(a), "*")
          }
          break
        }
        case"resize": {
          if (e.data.height <= 0) break;
          const t = document.getElementById(e.data.id);
          t instanceof HTMLIFrameElement && (t.height = `${e.data.height || 0}`);
          break
        }
      }
    });
  }, []);


  return (
      <div data-theme='ictsc' className={`${notoSansJP.className}`}>
        <Component {...pageProps} />
      </div>
  )
}
