import { use, useReducer } from "react";
import type { ProblemDetail } from "@app/features/problem";
import {
  Button,
  Field,
  Label,
  Tab,
  TabGroup,
  TabList,
  TabPanel,
  TabPanels,
  Textarea,
} from "@headlessui/react";
import { clsx } from "clsx";
import { MaterialSymbol } from "@app/components/material-symbol";
import { Markdown, Typography } from "@app/components/markdown";
import { NavbarLayoutContext } from "@app/components/app-shell";

export function ProblemPage(props: { problem: Promise<ProblemDetail> }) {
  const problem = use(props.problem);

  return <Layout content={<Content {...problem} />} sidebar={<Sidebar />} />;
}

export function Layout(props: {
  content: React.ReactNode;
  sidebar: React.ReactNode;
}) {
  const { navbarTransitioning } = use(NavbarLayoutContext);
  const [showSidebar, toggleSidebar] = useReducer((o) => !o, false);

  return (
    <div
      className={clsx(
        "relative mt-20",
        "[--span:calc(var(--content-width))] sm:[--span:calc(var(--content-width)/2)] lg:[--span:calc(var(--content-width)/3)]",
      )}
    >
      <div
        className={clsx(
          "overflow-y-auto px-40 pb-64",
          "w-[--span] lg:w-[calc(var(--span)*2)]",
          navbarTransitioning && "transition-[width]",
        )}
      >
        {props.content}
      </div>
      <div
        className={clsx(
          "fixed right-0 top-[--header-height] flex h-[--content-height] w-[--span] gap-4 px-12 pb-64 pt-20 transition duration-200",
          showSidebar && "bg-surface-0",
          !showSidebar &&
            "translate-x-[calc(var(--span)-64px)] sm:translate-x-0",
          "sm:w-[--span] sm:pl-0",
          navbarTransitioning && "transition-[width]",
        )}
      >
        <Button
          className="grid h-40 w-[40px] place-items-center rounded-full transition data-[hover]:bg-surface-1 data-[active]:opacity-80 sm:hidden"
          onClick={toggleSidebar}
          title={showSidebar ? "閉じる" : "サイドバーを開く"}
        >
          <MaterialSymbol size={24} icon={showSidebar ? "close" : "menu"} />
        </Button>
        <div
          className={clsx(
            "flex flex-1 transition-opacity ease-out",
            !showSidebar && "opacity-0 sm:opacity-100",
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

export function Sidebar() {
  return (
    <TabGroup className="flex flex-1 flex-col">
      <TabList className="flex flex-row gap-4">
        <SidebarTab>新規回答</SidebarTab>
        <SidebarTab>回答一覧</SidebarTab>
        {/* <Tab>再展開</Tab> */}
      </TabList>
      <TabPanels className="mx-8 mt-16 flex flex-1">
        <TabPanel className="flex flex-1">
          <SubmissionForm />
        </TabPanel>
        <TabPanel>
          <div>回答一覧</div>
        </TabPanel>
      </TabPanels>
    </TabGroup>
  );
}

function SidebarTab(props: { children?: React.ReactNode }) {
  return (
    <Tab className="group relative rounded-8 px-8 text-16 transition data-[hover]:bg-surface-1 data-[selected]:text-primary data-[active]:opacity-80">
      <div className="py-8 group-data-[selected]:text-primary">
        {props.children}
      </div>
      <div className="absolute bottom-0 mx-2 h-2 w-[calc(100%-20px)] rounded-t-full bg-transparent group-data-[selected]:bg-primary" />
    </Tab>
  );
}

function SubmissionForm() {
  return (
    <form className="flex flex-1 flex-col">
      <Field className="flex flex-1">
        <Label className="sr-only">回答</Label>
        <Textarea
          className="flex-1 resize-none rounded-12 border border-text p-12"
          placeholder="お世話になっております、チーム◯◯◯です。"
        />
      </Field>
      <button
        type="submit"
        className="mt-20 flex items-center justify-center self-end rounded-12 bg-surface-2 py-16 pl-24 pr-20 shadow-md transition hover:opacity-80 active:shadow-none"
      >
        <div className="text-16 font-bold">回答する</div>
        <MaterialSymbol icon="send" size={24} />
      </button>
    </form>
  );
}
