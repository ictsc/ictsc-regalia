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
            <div
              key={announce.slug}
              className="rounded-8 bg-surface-1 flex w-full items-center transition"
            >
              <Button as={Fragment}>
                <Link
                  className="text-16 hover:bg-surface-2 rounded-l-8 flex min-w-0 flex-1 items-center gap-8 py-4 pl-20 font-bold transition data-[active]:opacity-50"
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
              <ReadToggleButton
                slug={announce.slug}
                className="text-14 hover:bg-surface-2 rounded-r-8 flex shrink-0 items-center gap-4 px-16 py-4 transition"
              />
            </div>
          ))
        )}
      </div>
    </>
  );
}
