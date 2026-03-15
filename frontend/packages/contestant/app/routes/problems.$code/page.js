"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeploymentItem = exports.DeploymentList = exports.EmptyDeploymentList = exports.Deployments = exports.EmptySubmissionList = exports.SubmissionListItem = exports.SubmissionList = exports.SubmissionListContainer = exports.SubmissionForm = void 0;
exports.Page = Page;
exports.Content = Content;
var react_1 = require("react");
var react_2 = require("@headlessui/react");
var clsx_1 = require("clsx");
var material_symbol_1 = require("../../components/material-symbol");
var markdown_1 = require("../../components/markdown");
var app_shell_1 = require("../../components/app-shell");
var title_1 = require("../../components/title");
var submission_form_1 = require("./submission-form");
Object.defineProperty(exports, "SubmissionForm", { enumerable: true, get: function () { return submission_form_1.SubmissionForm; } });
var submission_list_1 = require("./submission-list");
Object.defineProperty(exports, "SubmissionListContainer", { enumerable: true, get: function () { return submission_list_1.SubmissionListContainer; } });
Object.defineProperty(exports, "SubmissionList", { enumerable: true, get: function () { return submission_list_1.SubmissionList; } });
Object.defineProperty(exports, "SubmissionListItem", { enumerable: true, get: function () { return submission_list_1.SubmissionListItem; } });
Object.defineProperty(exports, "EmptySubmissionList", { enumerable: true, get: function () { return submission_list_1.EmptySubmissionList; } });
var deployments_1 = require("./deployments");
Object.defineProperty(exports, "Deployments", { enumerable: true, get: function () { return deployments_1.Deployments; } });
Object.defineProperty(exports, "EmptyDeploymentList", { enumerable: true, get: function () { return deployments_1.EmptyDeploymentList; } });
Object.defineProperty(exports, "DeploymentList", { enumerable: true, get: function () { return deployments_1.DeploymentList; } });
Object.defineProperty(exports, "DeploymentItem", { enumerable: true, get: function () { return deployments_1.DeploymentItem; } });
function Page(props) {
    return (<Layout content={props.content} sidebar={<Sidebar onChange={props.onTabChange} redeployable={props.redeployable} submissionForm={props.submissionForm} submissionList={props.submissionList} deploymentList={props.deploymentList}/>}/>);
}
function Layout(props) {
    var navbarTransitioning = (0, react_1.use)(app_shell_1.NavbarLayoutContext).navbarTransitioning;
    var _a = (0, react_1.useReducer)(function (o) { return !o; }, false), showSidebar = _a[0], toggleSidebar = _a[1];
    return (<div className={(0, clsx_1.clsx)("relative mt-20", "[--span:calc(var(--content-width))] lg:[--span:calc(var(--content-width)/2)] xl:[--span:calc(var(--content-width)/3)]")}>
      <div className={(0, clsx_1.clsx)("overflow-y-auto px-40 pb-64", "w-(--span) xl:w-[calc(var(--span)*2)]", navbarTransitioning && "transition-[width]")}>
        {props.content}
      </div>
      <div className={(0, clsx_1.clsx)("fixed top-(--header-height) right-0 flex h-(--content-height) w-(--span) gap-4 px-12 pt-20 pb-64", showSidebar && "bg-surface-0", !showSidebar &&
            "translate-x-[calc(var(--span)-64px)] lg:translate-x-0", "lg:w-(--span) lg:pl-0", navbarTransitioning
            ? "transition-[width,translate]"
            : "transition-translate")}>
        <react_2.Button className="data-[hover]:bg-surface-1 grid h-40 w-[40px] place-items-center rounded-full transition data-[active]:opacity-80 lg:hidden" onClick={toggleSidebar} title={showSidebar ? "閉じる" : "サイドバーを開く"}>
          <material_symbol_1.MaterialSymbol size={24} icon={showSidebar ? "close" : "menu"}/>
        </react_2.Button>
        <div className={(0, clsx_1.clsx)("flex flex-1 transition-opacity ease-out", !showSidebar && "opacity-0 lg:opacity-100")}>
          {props.sidebar}
        </div>
      </div>
    </div>);
}
function Content(props) {
    return (<>
      <title_1.Title>{"".concat(props.code, ":").concat(props.title)}</title_1.Title>
      <markdown_1.Typography className="flex-1">
        <h1>
          {props.code}: {props.title}
        </h1>
        <markdown_1.Markdown>{props.body}</markdown_1.Markdown>
      </markdown_1.Typography>
    </>);
}
function Sidebar(props) {
    var _a = (0, react_1.useState)(0), tabIndex = _a[0], setTabIndex = _a[1];
    var onChange = function (index) {
        (0, react_1.startTransition)(function () {
            setTabIndex(index);
            if (props.onChange)
                props.onChange();
        });
    };
    return (<react_2.TabGroup tabIndex={tabIndex} onChange={onChange} className="flex flex-1 flex-col">
      <react_2.TabList className="flex flex-row gap-4">
        <SidebarTab>新規解答</SidebarTab>
        <SidebarTab>解答一覧</SidebarTab>
        <SidebarTab disabled={!props.redeployable}>再展開</SidebarTab>
      </react_2.TabList>
      <react_2.TabPanels className="mt-16 size-full bg-transparent px-8">
        <react_1.Suspense>
          <SidebarTabPanel>{props.submissionForm}</SidebarTabPanel>
          <SidebarTabPanel>{props.submissionList}</SidebarTabPanel>
          <SidebarTabPanel>{props.deploymentList}</SidebarTabPanel>
        </react_1.Suspense>
      </react_2.TabPanels>
    </react_2.TabGroup>);
}
function SidebarTab(props) {
    return (<react_2.Tab disabled={props.disabled} className="group rounded-8 text-16 data-[hover]:bg-surface-1 data-[selected]:text-primary relative px-8 transition data-[active]:opacity-80">
      <div className="group-data-[disabled]:text-disabled group-data-[selected]:text-primary py-8">
        {props.children}
      </div>
      <div className="group-data-[selected]:bg-primary absolute bottom-0 mx-2 h-2 w-[calc(100%-20px)] rounded-t-full bg-transparent"/>
    </react_2.Tab>);
}
function SidebarTabPanel(props) {
    return (<react_2.TabPanel className="size-full" unmount={false}>
      {function (_a) {
            var selected = _a.selected;
            return (<react_2.Transition show={selected}>
          <div className={(0, clsx_1.clsx)("size-full transition-opacity duration-150 data-[closed]:opacity-0", props.className)}>
            {props.children}
          </div>
        </react_2.Transition>);
        }}
    </react_2.TabPanel>);
}
