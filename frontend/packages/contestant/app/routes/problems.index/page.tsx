import type { Problem } from "@ictsc/proto/contestant/v1";
import { clsx } from "clsx";
import { Fragment } from "react";
import {
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
} from "@headlessui/react";
import type { Notice } from "../../features/announce";
import { ProblemItem } from "./problem-item";
import { UnreadAnnouncesBanner } from "./unread-announces-banner";
import { protoScoreToProps } from "../../features/score";
import { Title } from "../../components/title";
import { MaterialSymbol } from "../../components/material-symbol";
import {
  groupProblems,
  type GroupScheduleInfo,
} from "../../features/problem/group";

type PageProps = {
  problems: Problem[];
  notices: Notice[];
};

function ScheduleLabel(props: { schedules: GroupScheduleInfo[] }) {
  if (props.schedules.length === 0) {
    return <span className="opacity-50">スケジュール未設定</span>;
  }
  return props.schedules.map((s, i) => (
    <Fragment key={s.name}>
      {i > 0 && " / "}
      <span
        className={clsx(
          s.temporalStatus === "past" && "opacity-50",
          s.temporalStatus === "current" && "text-primary",
        )}
      >
        {s.name}
      </span>
    </Fragment>
  ));
}

export function ProblemsPage(props: PageProps) {
  const groups = groupProblems(props.problems);

  return (
    <>
      <Title>問題一覧</Title>
      <div className="mx-16 my-64 flex flex-col gap-16">
        <div className="ml-16">
          <UnreadAnnouncesBanner notices={props.notices} />
        </div>
        {groups.map((group) => (
          <Disclosure key={group.key} as="section" defaultOpen>
            <DisclosureButton className="group/disc flex w-full cursor-pointer items-center gap-16">
              <MaterialSymbol
                icon="arrow_forward_ios"
                size={24}
                className="text-disabled transition-transform group-data-[open]/disc:rotate-90"
              />
              <div className="border-disabled flex-1 border-t" />
              <h2 className="text-24 shrink-0 text-center font-bold">
                <ScheduleLabel schedules={group.schedules} />
              </h2>
              <div className="border-disabled flex-1 border-t" />
            </DisclosureButton>
            <DisclosurePanel className="mt-16">
              <ul className="grid grid-flow-row grid-cols-1 gap-x-40 gap-y-24 lg:grid-cols-2">
                {group.problems.map((problem) => (
                  <li key={problem.code}>
                    <ProblemItem
                      code={problem.code}
                      title={problem.title}
                      score={protoScoreToProps(problem.maxScore, problem.score)}
                      submissionStatus={problem.submissionStatus}
                    />
                  </li>
                ))}
              </ul>
            </DisclosurePanel>
          </Disclosure>
        ))}
      </div>
    </>
  );
}
