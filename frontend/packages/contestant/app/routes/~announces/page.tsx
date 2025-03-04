import { Fragment } from "react";
import { Button } from "@headlessui/react";
import { Link } from "@tanstack/react-router";
import { clsx } from "clsx";

export type AnnounceProps = {
  announce: string;
};

export function AnnounceList(props: AnnounceProps) {
  return (
    <div className="flex h-full flex-col items-center justify-center">
      <Button as={Fragment}>
        {({ active }) => {
          return (
            <Link
              to="/problems"
              title="アナウンス"
              className={clsx(
                "mb-16 w-full rounded-8 border-text bg-surface-2 py-4 pl-40 text-16 font-bold",
                active ? "opacity-50" : "opacity-100",
              )}
            >
              <div>{props.announce}</div>
            </Link>
          );
        }}
      </Button>
    </div>
  );
}
