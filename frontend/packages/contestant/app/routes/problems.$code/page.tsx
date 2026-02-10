import {
  type ReactNode,
  startTransition,
  Suspense,
  use,
  useReducer,
  useState,
} from "react";
import {
  Button,
  Tab,
  TabGroup,
  TabList,
  TabPanel,
  TabPanels,
  Transition,
} from "@headlessui/react";
import { clsx } from "clsx";
import { MaterialSymbol } from "../../components/material-symbol";
import { Markdown, Typography } from "../../components/markdown";
import { NavbarLayoutContext } from "../../components/app-shell";
import { Title } from "../../components/title";

export { SubmissionForm } from "./submission-form";
export {
  SubmissionListContainer,
  SubmissionList,
  SubmissionListItem,
  EmptySubmissionList,
} from "./submission-list";
export {
  Deployments,
  EmptyDeploymentList,
  DeploymentList,
  DeploymentItem,
} from "./deployments";

export function Page(props: {
  content: ReactNode;
  submissionForm: ReactNode;
  submissionList: ReactNode;
  deploymentList: ReactNode;
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
          deploymentList={props.deploymentList}
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
          "w-(--span) xl:w-[calc(var(--span)*2)]",
          navbarTransitioning && "transition-[width]",
        )}
      >
        {props.content}
      </div>
      <div
        className={clsx(
          "fixed top-(--header-height) right-0 flex h-(--content-height) w-(--span) gap-4 px-12 pt-20 pb-64",
          showSidebar && "bg-surface-0",
          !showSidebar &&
            "translate-x-[calc(var(--span)-64px)] lg:translate-x-0",
          "lg:w-(--span) lg:pl-0",
          navbarTransitioning
            ? "transition-[width,translate]"
            : "transition-translate",
        )}
      >
        <Button
          className="data-[hover]:bg-surface-1 grid h-40 w-[40px] place-items-center rounded-full transition data-[active]:opacity-80 lg:hidden"
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
    <>
      <Title>{`${props.code}:${props.title}`}</Title>
      <Typography className="flex-1">
        <h1>
          {props.code}: {props.title}
        </h1>
        <Markdown>{props.body}</Markdown>
      </Typography>
    </>
  );
}

function Sidebar(props: {
  submissionForm: ReactNode;
  submissionList: ReactNode;
  deploymentList: ReactNode;
  redeployable: boolean;
  onChange?: () => void;
}) {
  const [tabIndex, setTabIndex] = useState(0);
  const onChange = (index: number) => {
    startTransition(() => {
      setTabIndex(index);
      if (props.onChange) props.onChange();
    });
  };
  return (
    <TabGroup
      tabIndex={tabIndex}
      onChange={onChange}
      className="flex flex-1 flex-col"
    >
      <TabList className="flex flex-row gap-4">
        <SidebarTab>新規解答</SidebarTab>
        <SidebarTab>解答一覧</SidebarTab>
        <SidebarTab disabled={!props.redeployable}>再展開</SidebarTab>
      </TabList>
      <TabPanels className="mt-16 size-full bg-transparent px-8">
        <Suspense>
          <SidebarTabPanel>{props.submissionForm}</SidebarTabPanel>
          <SidebarTabPanel>{props.submissionList}</SidebarTabPanel>
          <SidebarTabPanel>{props.deploymentList}</SidebarTabPanel>
        </Suspense>
      </TabPanels>
    </TabGroup>
  );
}

function SidebarTab(props: { disabled?: boolean; children?: ReactNode }) {
  return (
    <Tab
      disabled={props.disabled}
      className="group rounded-8 text-16 data-[hover]:bg-surface-1 data-[selected]:text-primary relative px-8 transition data-[active]:opacity-80"
    >
      <div className="group-data-[disabled]:text-disabled group-data-[selected]:text-primary py-8">
        {props.children}
      </div>
      <div className="group-data-[selected]:bg-primary absolute bottom-0 mx-2 h-2 w-[calc(100%-20px)] rounded-t-full bg-transparent" />
    </Tab>
  );
}

function SidebarTabPanel(props: { className?: string; children?: ReactNode }) {
  return (
    <TabPanel className="size-full" unmount={false}>
      {({ selected }) => (
        <Transition show={selected}>
          <div
            className={clsx(
              "size-full transition-opacity duration-150 data-[closed]:opacity-0",
              props.className,
            )}
          >
            {props.children}
          </div>
        </Transition>
      )}
    </TabPanel>
  );
}
