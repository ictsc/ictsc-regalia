import { Fragment } from "react";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import type { Notice } from "@ictsc/proto/contestant/v1";
import { clsx } from "clsx";

type AnnounceProps = {
  announces: Notice[]
};

export function AnnounceList(props: AnnounceProps) {
  return (
    <div className="flex mt-64  flex-col items-center justify-center">
      {props.announces.map((announce)=> (
      <Button as={Fragment}>
        {({ active }) => {
          return (
            <Link
              to="/"
              title="アナウンス"
              className={clsx(
                "mb-16 w-[700px] rounded-8 border-text bg-surface-2 py-4 px-40 text-16 font-bold",
                active ? "opacity-50" : "opacity-100",
              )}
            >
              <div>{announce.title}</div>
            </Link>
          );
        }}
      </Button>
      ))};
    </div>
  );
}
