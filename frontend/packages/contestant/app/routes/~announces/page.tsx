import { Fragment } from "react";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import type { Notice } from "@ictsc/proto/contestant/v1";
import { clsx } from "clsx";

type AnnounceProps = {
  announces: Notice[];
};

export function AnnounceList(props: AnnounceProps) {
  return (
    <div className="mt-64 flex flex-col items-center justify-center">
      {props.announces.map((announce, index) => (
        <Button key={index} as={Fragment}>
          {({ active }) => {
            return (
              <Link
                key={index}
                to="/"
                className={clsx(
                  "mb-16 w-[700px] rounded-8 border-text bg-surface-2 px-40 py-4 text-16 font-bold",
                  active ? "opacity-50" : "opacity-100",
                )}
              >
                <div>{announce.title}</div>
              </Link>
            );
          }}
        </Button>
      ))}
      ;
    </div>
  );
}
