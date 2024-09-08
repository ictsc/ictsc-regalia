import React from "react";

import { Metadata } from "next";

const title = "問題一覧";

export const metadata: Metadata = {
  title,
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return children;
}
