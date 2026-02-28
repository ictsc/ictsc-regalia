import { isAfter } from "date-fns";
import {
  Suspense,
  use,
  useDeferredValue,
  useEffect,
  useState,
  useTransition,
  type ReactNode,
} from "react";
import { type Schedule } from "@ictsc/proto/contestant/v1";
import { currentEndAt as endAt } from "./feature";
import { ScheduleContext } from "./use-schedule";

export function ScheduleProvider(props: {
  readonly initialData: Promise<Schedule | null>;
  readonly loadData: () => Promise<Schedule | null>;
  readonly children?: ReactNode;
}) {
  const deferredInitialData = useDeferredValue(props.initialData);

  const [promiseState, setPromise] = useState({
    data: props.initialData,
    base: props.initialData,
  });
  const [isStatePending, startTransision] = useTransition();

  const [promise, isPending] =
    promiseState.base === props.initialData
      ? [promiseState.data, isStatePending]
      : [deferredInitialData, deferredInitialData !== props.initialData];

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
