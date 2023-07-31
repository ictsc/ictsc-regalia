import React from "react";

import { Metadata } from "next";

import ICTSCTitle from "@/components/Title";

const title = "参加者一覧";

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
