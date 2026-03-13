import { Fragment } from "react";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import type { Notice } from "@ictsc/proto/contestant/v1";
import { MaterialSymbol } from "../components/material-symbol";
import { Title } from "../components/title";
import { ReadToggleButton } from "./problems.index/unread-announces-banner";

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
          props.announces.map((announce) => (
            <div key={announce.slug} className="flex w-full items-center gap-8">
              <Button as={Fragment}>
                <Link
                  className="rounded-8 bg-surface-1 text-16 data-[hover]:bg-surface-2 flex min-w-0 flex-1 items-center gap-8 py-4 pr-40 pl-20 font-bold transition data-[active]:opacity-50"
                  to="/announces/$slug"
                  params={{ slug: announce.slug }}
                >
                  <MaterialSymbol
                    icon="arrow_forward_ios"
                    size={20}
                    className="text-icon shrink-0"
                  />
                  <span className="truncate">{announce.title}</span>
                </Link>
              </Button>
              <ReadToggleButton slug={announce.slug} />
            </div>
          ))
        )}
      </div>
    </>
  );
}
