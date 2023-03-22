import React from "react";

import BaseLayout from "@/layouts/BaseLayout";

interface Props {
  title: string;
  children: React.ReactNode;
}

const CommonLayout = ({ title, children }: Props) => {
  return (
    <BaseLayout title={title}>
      <h1 className={"title-ictsc text-center py-12"}>{title}</h1>
      {children}
    </BaseLayout>
  );
};

export default CommonLayout;
