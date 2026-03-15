"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_1 = require("react");
var react_router_1 = require("@tanstack/react-router");
var ___root_1 = require("./~__root");
var signin_page_1 = require("./signin.page");
exports.Route = (0, react_router_1.createFileRoute)("/signin")({
    component: RouteComponent,
    validateSearch: function (search) {
        return {
            next: typeof search.next === "string" ? search.next : undefined,
        };
    },
});
function RouteComponent() {
    var viewer = ___root_1.Route.useLoaderData().viewer;
    return (<react_1.Suspense>
      <Page viewerPromise={viewer}/>
    </react_1.Suspense>);
}
function Page(_a) {
    var viewerPromise = _a.viewerPromise;
    var _b = exports.Route.useSearch().next, nextPath = _b === void 0 ? "/" : _b;
    // 既にログインしているならそのまま next に遷移する
    // ログインしているかどうかはログインの操作が可能かどうかには関係ない
    // そのためログインしているかの確認の前にページ自体は表示しておき，状態が確定したら遷移するか判断する
    var navigate = (0, react_router_1.useNavigate)();
    var deferredViewer = (0, react_1.useDeferredValue)(viewerPromise, null);
    var viewer = deferredViewer != null ? (0, react_1.use)(deferredViewer) : null;
    (0, react_1.useEffect)(function () {
        if (viewer == null)
            return;
        if (viewer.type !== "unauthenticated") {
            (0, react_1.startTransition)(function () { return navigate({ to: nextPath }); });
        }
    }, [viewer, navigate, nextPath]);
    var authURL = new URL("/api/auth/signin", window.location.origin);
    authURL.searchParams.set("next", nextPath);
    return <signin_page_1.SignInPage signInURL={authURL.toString()}/>;
}
