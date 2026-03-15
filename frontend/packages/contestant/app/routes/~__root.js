"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_1 = require("react");
var react_error_boundary_1 = require("react-error-boundary");
var react_router_1 = require("@tanstack/react-router");
var connect_1 = require("@connectrpc/connect");
var app_shell_1 = require("./app-shell");
var viewer_1 = require("../features/viewer");
var schedule_1 = require("../features/schedule");
exports.Route = (0, react_router_1.createRootRouteWithContext)()({
    component: Root,
    loader: function (_a) {
        var transport = _a.context.transport;
        return ({
            viewer: (0, viewer_1.fetchViewer)(transport),
            schedule: (0, schedule_1.fetchSchedule)(transport),
            loadSchedule: function () { return (0, schedule_1.fetchSchedule)(transport); },
        });
    },
});
var TanStackRouterDevtools = import.meta.env.DEV
    ? (0, react_1.lazy)(function () {
        return Promise.resolve().then(function () { return require("@tanstack/router-devtools"); }).then(function (mod) { return ({
            default: mod.TanStackRouterDevtools,
        }); });
    })
    : function () { return null; };
function Root() {
    var _a = exports.Route.useLoaderData(), viewer = _a.viewer, schedule = _a.schedule, loadSchedule = _a.loadSchedule;
    return (<>
      <schedule_1.ScheduleProvider initialData={schedule} loadData={loadSchedule}>
        <app_shell_1.AppShell viewer={viewer}>
          <react_error_boundary_1.ErrorBoundary FallbackComponent={ErrorFallback}>
            <react_1.Suspense fallback={null}>
              <react_router_1.Outlet />
            </react_1.Suspense>
          </react_error_boundary_1.ErrorBoundary>
        </app_shell_1.AppShell>
      </schedule_1.ScheduleProvider>
      <react_1.Suspense>
        <Redirector viewer={viewer}/>
      </react_1.Suspense>
      <react_1.Suspense>
        <TanStackRouterDevtools />
      </react_1.Suspense>
    </>);
}
function ErrorFallback(props) {
    var error = props.error;
    if (error instanceof connect_1.ConnectError && error.code === connect_1.Code.Unauthenticated) {
        return <UnauthorizedFallback {...props}/>;
    }
    if (error instanceof connect_1.ConnectError && error.code === connect_1.Code.PermissionDenied) {
        return <PermissionFallback {...props}/>;
    }
    // TODO: エラー画面を表示する
    return null;
}
function UnauthorizedFallback(_a) {
    var resetErrorBoundary = _a.resetErrorBoundary;
    // viewer の取得よりも先に認証エラーが発生した場合 ErrorBoundary が反応するが，
    // 認証エラーは Redirector により処理されるため，適当に再開しつつエラーを無視する
    (0, react_1.useEffect)(function () {
        resetErrorBoundary();
    }, [resetErrorBoundary]);
    return null;
}
function PermissionFallback(_a) {
    var resetErrorBoundary = _a.resetErrorBoundary;
    var navigate = (0, react_router_1.useNavigate)();
    (0, react_1.useEffect)(function () {
        resetErrorBoundary();
        (0, react_1.startTransition)(function () { return navigate({ to: "/" }); });
    }, [resetErrorBoundary, navigate]);
    return null;
}
function Redirector(_a) {
    var viewerPromise = _a.viewer;
    var routerState = (0, react_router_1.useRouterState)();
    var navigate = (0, react_router_1.useNavigate)();
    var viewer = (0, react_1.use)(viewerPromise);
    (0, react_1.useEffect)(function () {
        switch (viewer.type) {
            case "unauthenticated": {
                if (routerState.location.pathname !== "/signin") {
                    (0, react_1.startTransition)(function () {
                        return navigate({
                            to: "/signin",
                            search: { next: routerState.location.pathname },
                        });
                    });
                }
                break;
            }
            case "pre-signup": {
                if (routerState.location.pathname !== "/signup") {
                    (0, react_1.startTransition)(function () {
                        return navigate({
                            to: "/signup",
                            search: { next: routerState.location.pathname },
                        });
                    });
                }
                break;
            }
        }
    }, [routerState, navigate, viewer]);
    return null;
}
