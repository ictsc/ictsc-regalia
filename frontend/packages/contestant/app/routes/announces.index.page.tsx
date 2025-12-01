import { Fragment } from "react";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import type { Notice } from "@ictsc/proto/contestant/v1";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";

type AnnounceProps = {
  announces: Notice[];
};

export function AnnounceList(props: AnnounceProps) {
  return (
    <>
      <Title>アナウンス一覧</Title>
      <div className="mx-40 mt-64 flex flex-col items-center justify-center gap-16">
        {props.announces.length === 0 ? (
          <h1 className="font-bold">現在アナウンスはありません</h1>
        ) : (
          props.announces.map((announce) => {
            return (
              <Button as={Fragment} key={announce.slug}>
                <Link
                  className="flex w-full items-center gap-8 rounded-8 bg-surface-1 py-4 pl-20 pr-40 text-16 font-bold transition data-[hover]:bg-surface-2 data-[active]:opacity-50"
                  to="/announces/$slug"
                  params={{ slug: announce.slug }}
                >
                  <MaterialSymbol
                    icon="arrow_forward_ios"
                    size={20}
                    className="text-icon"
                  />
                  {announce.title}
                </Link>
              </Button>
            );
          })
        )}
      </div>
    </>
  );
}
