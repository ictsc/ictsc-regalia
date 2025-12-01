import { createFileRoute } from "@tanstack/react-router";
import { IndexPage, Timer } from "./index.page";
import { useSchedule } from "@app/features/schedule";
import { Phase } from "@ictsc/proto/contestant/v1";
import { timestampMs } from "@bufbuild/protobuf/wkt";

export const Route = createFileRoute("/")({
  component: Page,
});

function Page() {
  const [schedule] = useSchedule();
  return (
    <IndexPage
      phase={schedule?.phase ?? Phase.UNSPECIFIED}
      nextPhase={schedule?.nextPhase}
      timer={
        <Timer
          endMs={schedule?.endAt != null ? timestampMs(schedule.endAt) : 0}
        />
      }
    />
  );
}
