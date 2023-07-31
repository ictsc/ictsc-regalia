"use client";

import React from "react";

interface Props {
  title: string;
  children: React.ReactNode;
}

function CommonLayout({ title, children }: Props) {
  return (
    <>
      <h1 id="page-title" className="title-ictsc text-center py-12">
        {title}
      </h1>
      {children}
    </>
  );
}

export default CommonLayout;
