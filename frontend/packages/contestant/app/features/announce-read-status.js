"use strict";
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.useReadAnnouncements = useReadAnnouncements;
var react_1 = require("react");
var STORAGE_KEY = "readAnnouncements";
function getSnapshot() {
    if (typeof window === "undefined")
        return [];
    try {
        var stored = localStorage.getItem(STORAGE_KEY);
        if (stored == null)
            return [];
        return JSON.parse(stored);
    }
    catch (_a) {
        return [];
    }
}
var cachedSlugs = getSnapshot();
function getSnapshotCached() {
    return cachedSlugs;
}
function setSlugs(slugs) {
    cachedSlugs = slugs;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(slugs));
}
var listeners = [];
function handleStorageEvent(e) {
    if (e.key === STORAGE_KEY) {
        cachedSlugs = getSnapshot();
        notifyListeners();
    }
}
function subscribeInternal(callback) {
    if (listeners.length === 0) {
        window.addEventListener("storage", handleStorageEvent);
    }
    listeners.push(callback);
    return function () {
        listeners = listeners.filter(function (l) { return l !== callback; });
        if (listeners.length === 0) {
            window.removeEventListener("storage", handleStorageEvent);
        }
    };
}
function notifyListeners() {
    for (var _i = 0, listeners_1 = listeners; _i < listeners_1.length; _i++) {
        var listener = listeners_1[_i];
        listener();
    }
}
function useReadAnnouncements() {
    var readSlugs = (0, react_1.useSyncExternalStore)(subscribeInternal, getSnapshotCached, function () { return []; });
    var markAsRead = (0, react_1.useCallback)(function (slug) {
        var current = getSnapshot();
        if (!current.includes(slug)) {
            var next = __spreadArray(__spreadArray([], current, true), [slug], false);
            setSlugs(next);
            notifyListeners();
        }
    }, []);
    var markAllAsRead = (0, react_1.useCallback)(function (notices) {
        var current = getSnapshot();
        var currentSet = new Set(current);
        var newSlugs = notices
            .map(function (n) { return n.slug; })
            .filter(function (s) { return !currentSet.has(s); });
        if (newSlugs.length === 0)
            return;
        setSlugs(__spreadArray(__spreadArray([], current, true), newSlugs, true));
        notifyListeners();
    }, []);
    var markAsUnread = (0, react_1.useCallback)(function (slug) {
        var current = getSnapshot();
        if (current.includes(slug)) {
            var next = current.filter(function (s) { return s !== slug; });
            setSlugs(next);
            notifyListeners();
        }
    }, []);
    var isRead = (0, react_1.useCallback)(function (slug) {
        return readSlugs.includes(slug);
    }, [readSlugs]);
    var getUnreadNotices = (0, react_1.useCallback)(function (notices) {
        var readSet = new Set(readSlugs);
        return notices.filter(function (n) { return !readSet.has(n.slug); });
    }, [readSlugs]);
    return { markAsRead: markAsRead, markAsUnread: markAsUnread, markAllAsRead: markAllAsRead, isRead: isRead, getUnreadNotices: getUnreadNotices };
}
