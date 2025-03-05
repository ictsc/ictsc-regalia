import { type ComponentProps, type ReactNode, use, useReducer } from "react";
import {
  Button,
  Tab,
  TabGroup,
  TabList,
  TabPanel,
  TabPanels,
} from "@headlessui/react";
import { clsx } from "clsx";
import { MaterialSymbol } from "../../components/material-symbol";
import { Markdown, Typography } from "../../components/markdown";
import { NavbarLayoutContext } from "../../components/app-shell";
import { Score } from "../../components/score";

export { SubmissionForm } from "./submission-form";

export function Page(props: {
  content: ReactNode;
  submissionForm: ReactNode;
  submissionList: ReactNode;
  redeployable: boolean;

  onTabChange?: () => void;
}) {
  return (
    <Layout
      content={props.content}
      sidebar={
        <Sidebar
          onChange={props.onTabChange}
          redeployable={props.redeployable}
          submissionForm={props.submissionForm}
          submissionList={props.submissionList}
        />
      }
    />
  );
}

function Layout(props: { content: React.ReactNode; sidebar: React.ReactNode }) {
  const { navbarTransitioning } = use(NavbarLayoutContext);
  const [showSidebar, toggleSidebar] = useReducer((o) => !o, false);

  return (
    <div
      className={clsx(
        "relative mt-20",
        "[--span:calc(var(--content-width))] lg:[--span:calc(var(--content-width)/2)] xl:[--span:calc(var(--content-width)/3)]",
      )}
    >
      <div
        className={clsx(
          "overflow-y-auto px-40 pb-64",
          "w-[--span] xl:w-[calc(var(--span)*2)]",
          navbarTransitioning && "transition-[width]",
        )}
      >
        {props.content}
      </div>
      <div
        className={clsx(
          "fixed right-0 top-[--header-height] flex h-[--content-height] w-[--span] gap-4 px-12 pb-64 pt-20",
          showSidebar && "bg-surface-0",
          !showSidebar &&
            "translate-x-[calc(var(--span)-64px)] lg:translate-x-0",
          "lg:w-[--span] lg:pl-0",
          navbarTransitioning
            ? "transition-[width,transform]"
            : "transition-transform",
        )}
      >
        <Button
          className="grid h-40 w-[40px] place-items-center rounded-full transition data-[hover]:bg-surface-1 data-[active]:opacity-80 lg:hidden"
          onClick={toggleSidebar}
          title={showSidebar ? "閉じる" : "サイドバーを開く"}
        >
          <MaterialSymbol size={24} icon={showSidebar ? "close" : "menu"} />
        </Button>
        <div
          className={clsx(
            "flex flex-1 transition-opacity ease-out",
            !showSidebar && "opacity-0 lg:opacity-100",
          )}
        >
          {props.sidebar}
        </div>
      </div>
    </div>
  );
}

export function Content(props: { code: string; title: string; body: string }) {
  return (
    <Typography className="flex-1">
      <h1>
        {props.code}: {props.title}
      </h1>
      <Markdown>{props.body}</Markdown>
    </Typography>
  );
}

function Sidebar(props: {
  submissionForm: ReactNode;
  submissionList: ReactNode;
  redeployable: boolean;
  onChange?: () => void;
}) {
  return (
    <TabGroup onChange={props.onChange} className="flex flex-1 flex-col">
      <TabList className="flex flex-row gap-4">
        <SidebarTab>新規解答</SidebarTab>
        <SidebarTab>解答一覧</SidebarTab>
        <SidebarTab disabled={!props.redeployable}>再展開</SidebarTab>
      </TabList>
      <TabPanels className="mt-16 size-full bg-transparent px-8">
        <TabPanel className="size-full">{props.submissionForm}</TabPanel>
        <TabPanel className="size-full rounded-12 bg-disabled py-12">
          <div className="size-full overflow-y-auto px-12 [scrollbar-gutter:stable_both-edges]">
            {props.submissionList}
          </div>
        </TabPanel>
      </TabPanels>
    </TabGroup>
  );
}

function SidebarTab(props: { disabled?: boolean; children?: ReactNode }) {
  return (
    <Tab
      disabled={props.disabled}
      className="group relative rounded-8 px-8 text-16 transition data-[hover]:bg-surface-1 data-[selected]:text-primary data-[active]:opacity-80"
    >
      <div className="py-8 group-data-[disabled]:text-disabled group-data-[selected]:text-primary">
        {props.children}
      </div>
      <div className="absolute bottom-0 mx-2 h-2 w-[calc(100%-20px)] rounded-t-full bg-transparent group-data-[selected]:bg-primary" />
    </Tab>
  );
}

export function SubmissionList(props: { readonly children?: ReactNode }) {
  return (
    <ul className="flex size-full flex-col gap-16 py-12">{props.children}</ul>
  );
}

export function EmptySubmissionList() {
  return (
    <div className="grid size-full place-items-center text-16 font-bold text-text">
      解答はまだありません！
    </div>
  );
}

const submissionListDateTimeFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "medium",
  timeStyle: "short",
});

export function SubmissionListItem(props: {
  readonly id: number;
  readonly submittedAt: string;
  readonly score: ComponentProps<typeof Score>;
}) {
  return (
    <li className="flex justify-between gap-8 rounded-12 bg-surface-0 p-16">
      <div className="flex flex-col">
        <h2 className="text-20 font-bold text-[#000]">#{props.id}</h2>
        <h3 className="text-12">
          提出:{" "}
          {submissionListDateTimeFormatter.format(new Date(props.submittedAt))}
        </h3>
        <p className="mt-4 text-12"></p>
      </div>
      <Score {...props.score} />
    </li>
  );
}
