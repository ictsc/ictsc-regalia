"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UnreadAnnouncesBanner = UnreadAnnouncesBanner;
exports.ReadToggleButton = ReadToggleButton;
var react_router_1 = require("@tanstack/react-router");
var announce_read_status_1 = require("../../features/announce-read-status");
var material_symbol_1 = require("../../components/material-symbol");
function UnreadAnnouncesBanner(_a) {
    var notices = _a.notices;
    var _b = (0, announce_read_status_1.useReadAnnouncements)(), getUnreadNotices = _b.getUnreadNotices, markAllAsRead = _b.markAllAsRead;
    var unreadNotices = getUnreadNotices(notices);
    if (unreadNotices.length === 0) {
        return null;
    }
    var label = "\u672A\u8AAD\u306E\u65B0\u7740\u30A2\u30CA\u30A6\u30F3\u30B9\u304C ".concat(unreadNotices.length, " \u4EF6\u3042\u308A\u307E\u3059");
    return (<div className="flex min-w-0 items-center gap-8">
      <react_router_1.Link to="/announces" title={label} className="text-14 bg-surface-1 hover:bg-surface-2 flex min-w-0 items-center gap-4 rounded-full py-4 pr-12 pl-8 font-bold transition">
        <material_symbol_1.MaterialSymbol icon="notifications" size={20} className="text-icon shrink-0"/>
        <span className="truncate">{label}</span>
      </react_router_1.Link>
      <button type="button" title="既読にする" className="text-14 bg-surface-1 hover:bg-surface-2 flex shrink-0 items-center gap-4 rounded-full py-4 pr-12 pl-8 transition" onClick={function () { return markAllAsRead(notices); }}>
        <material_symbol_1.MaterialSymbol icon="done_all" size={20} className="text-icon"/>
        既読にする
      </button>
    </div>);
}
function ReadToggleButton(_a) {
    var slug = _a.slug;
    var _b = (0, announce_read_status_1.useReadAnnouncements)(), isRead = _b.isRead, markAsRead = _b.markAsRead, markAsUnread = _b.markAsUnread;
    var read = isRead(slug);
    return (<button type="button" className="text-14 bg-surface-1 hover:bg-surface-2 flex shrink-0 items-center gap-4 rounded-full py-4 pr-12 pl-8 transition" onClick={function () { return (read ? markAsUnread(slug) : markAsRead(slug)); }}>
      <material_symbol_1.MaterialSymbol icon={read ? "mark_email_read" : "mark_email_unread"} size={20} className="text-icon"/>
      {read ? "未読にする" : "既読にする"}
    </button>);
}
