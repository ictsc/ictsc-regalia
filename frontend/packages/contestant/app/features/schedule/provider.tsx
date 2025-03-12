import { isAfter } from "date-fns";
import {
  Suspense,
  use,
  useEffect,
  useState,
  useTransition,
  type ReactNode,
} from "react";
import { type Schedule } from "@ictsc/proto/contestant/v1";
import { endAt } from "./feature";
import { ScheduleContext } from "./use-schedule";

export function ScheduleProvider(props: {
  readonly initialData: Promise<Schedule | null>;
  readonly loadData: () => Promise<Schedule | null>;
  readonly children?: ReactNode;
}) {
  const [promiseState, setPromise] = useState({
    data: props.initialData,
    base: props.initialData,
  });
  const promise =
    promiseState.base === props.initialData
      ? promiseState.data
      : props.initialData;
  const [isPending, startTransision] = useTransition();

  return (
    <>
      <ScheduleContext value={{ promise, isPending }}>
        {props.children}
      </ScheduleContext>
      <Suspense>
        <ScheduleReloader
          schedule={promise}
          load={() => {
            if (isPending) return;
            startTransision(() => {
              setPromise({
                data: props.loadData(),
                base: props.initialData,
              });
            });
          }}
        />
      </Suspense>
    </>
  );
}

function ScheduleReloader(props: {
  schedule: Promise<Schedule | null>;
  load: () => void;
}) {
  const { load } = props;
  const schedule = use(props.schedule);
  useEffect(() => {
    const timer = setInterval(() => {
      if (schedule == null) return;
      const end = endAt(schedule);
      if (end == null || isAfter(new Date(), end)) {
        load();
      }
    }, 1000);
    return () => clearInterval(timer);
  }, [schedule, load]);

  return null;
}
