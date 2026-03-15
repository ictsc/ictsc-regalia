"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Route = void 0;
var react_1 = require("react");
var react_router_1 = require("@tanstack/react-router");
var markdown_1 = require("../components/markdown");
var title_1 = require("../components/title");
var announce_read_status_1 = require("../features/announce-read-status");
var unread_announces_banner_1 = require("./problems.index/unread-announces-banner");
exports.Route = (0, react_router_1.createLazyFileRoute)("/announces/$slug")({
    component: RouteComponent,
});
function RouteComponent() {
    var announce = exports.Route.useLoaderData().announce;
    var slug = exports.Route.useParams().slug;
    var markAsRead = (0, announce_read_status_1.useReadAnnouncements)().markAsRead;
    (0, react_1.useEffect)(function () {
        if (announce != null) {
            markAsRead(slug);
        }
    }, [slug, markAsRead, announce]);
    return (<>
      <title_1.Title>{announce === null || announce === void 0 ? void 0 : announce.title}</title_1.Title>
      <div className="mt-64 flex w-full px-40">
        {announce == null ? (<h1 className="mx-auto font-bold">アナウンスがありません</h1>) : (<markdown_1.Typography className="w-full">
            <h1 className="flex items-center justify-between gap-8">
              {announce.title}
              <unread_announces_banner_1.ReadToggleButton slug={slug}/>
            </h1>
            <markdown_1.Markdown>{announce.body}</markdown_1.Markdown>
          </markdown_1.Typography>)}
      </div>
    </>);
}
