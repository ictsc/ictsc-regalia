import "../styles/globals.css";
import React from "react";

import localFont from "@next/font/local";
import { Metadata } from "next";

import Providers from "@/app/providers";
import ICTSCNavBar from "@/components/Navbar";
import { site } from "@/components/_const";

const notoSansJP = localFont({
  variable: "--font-noto-sans-jp",
  src: "dist/fonts/NotoSansJP-VF.woff2",
});

export const metadata: Metadata = {
  title: {
    default: site,
    template: `%s - ${site}`,
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja" data-theme="ictsc" className={notoSansJP.className}>
      <body>
        <Providers>
          <header>
            <ICTSCNavBar />
          </header>
          {children}
        </Providers>
      </body>
    </html>
  );
}
