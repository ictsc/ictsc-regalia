import {
  Suspense,
  use,
  useDeferredValue,
  useEffect,
  useEffectEvent,
  useState,
  useTransition,
  type ReactNode,
} from "react";
import { type Schedule } from "@ictsc/proto/contestant/v1";
import { nextReloadAt } from "./feature";
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
  const onLoad = useEffectEvent(props.load);
  const schedule = use(props.schedule);

  useEffect(() => {
    const reloadAt = nextReloadAt(schedule);
    if (reloadAt == null) return;

    const delayMs = reloadAt.getTime() - Date.now();
    if (delayMs <= 0) {
      onLoad();
      return;
    }

    const timer = window.setTimeout(() => {
      onLoad();
    }, delayMs);
    return () => window.clearTimeout(timer);
  }, [schedule]);

  return null;
}
