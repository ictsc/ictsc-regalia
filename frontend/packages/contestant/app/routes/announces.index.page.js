"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AnnounceList = AnnounceList;
var react_1 = require("react");
var react_2 = require("@headlessui/react");
var react_router_1 = require("@tanstack/react-router");
var material_symbol_1 = require("../components/material-symbol");
var title_1 = require("../components/title");
var unread_announces_banner_1 = require("./problems.index/unread-announces-banner");
function AnnounceList(props) {
    return (<>
      <title_1.Title>アナウンス一覧</title_1.Title>
      <div className="mx-40 mt-64 flex flex-col items-center justify-center gap-16">
        {props.announces.length === 0 ? (<h1 className="font-bold">現在アナウンスはありません</h1>) : (props.announces.map(function (announce) { return (<div key={announce.slug} className="flex w-full items-center gap-8">
              <react_2.Button as={react_1.Fragment}>
                <react_router_1.Link className="rounded-8 bg-surface-1 text-16 data-[hover]:bg-surface-2 flex min-w-0 flex-1 items-center gap-8 py-4 pr-40 pl-20 font-bold transition data-[active]:opacity-50" to="/announces/$slug" params={{ slug: announce.slug }}>
                  <material_symbol_1.MaterialSymbol icon="arrow_forward_ios" size={20} className="text-icon shrink-0"/>
                  <span className="truncate">{announce.title}</span>
                </react_router_1.Link>
              </react_2.Button>
              <unread_announces_banner_1.ReadToggleButton slug={announce.slug}/>
            </div>); }))}
      </div>
    </>);
}
