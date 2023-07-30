import React, { useEffect } from "react";

import Head from "next/head";

import ICTSCNavBar from "@/components/Navbar";
import { site } from "@/components/_const";

interface Props {
  title: string;
  children: React.ReactNode;
}

function BaseLayout({ title, children }: Props) {
  useEffect(() => {
    // title
    document.title = `${title} - ${site}`;
  }, [title]);

  return (
    <>
      <Head>
        <title>
          {title} - {site}
        </title>
      </Head>
      <ICTSCNavBar />
      {children}
    </>
  );
}

export default BaseLayout;
