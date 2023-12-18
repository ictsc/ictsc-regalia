import React from "react";

import { Metadata } from "next";
import { notFound } from "next/navigation";

import { preRoundMode } from "@/components/_const";
import ICTSCTitle from "@/components/title";

const title = "ランキング";

export const metadata: Metadata = {
  title,
};

export default function Layout({ children }: { children: React.ReactNode }) {
  if (preRoundMode) {
    return notFound();
  }

  return (
    <>
      <ICTSCTitle title={title} />
      <main className="container-ictsc overflow-x-auto">{children}</main>;
    </>
  );
}
