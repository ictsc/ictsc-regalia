import React from "react";

import Head from "next/head";

import { site } from "@/components/_const";
import ICTSCNavBar from "@/components/Navbar";

interface Props {
  title: string;
  children: React.ReactNode;
}

const BaseLayout = ({ title, children }: Props) => {
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
};

export default BaseLayout;
