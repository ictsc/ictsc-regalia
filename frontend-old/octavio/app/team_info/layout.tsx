import React from "react";

import { Metadata } from "next";

import ICTSCTitle from "@/components/title";

const title = "チーム情報";

export const metadata: Metadata = {
  title,
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <ICTSCTitle title={title} />
      <main className="container-ictsc">{children}</main>
    </>
  );
}
