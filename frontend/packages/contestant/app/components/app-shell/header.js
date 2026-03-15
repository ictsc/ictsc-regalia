"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Header = Header;
var react_router_1 = require("@tanstack/react-router");
var logo_1 = require("@app/components/logo");
function Header(_a) {
    var contestState = _a.contestState, accountMenu = _a.accountMenu;
    return (<div className="border-primary bg-surface-0 flex size-full items-center border-b-[3px]">
      <div className="ml-4 flex-none sm:ml-16">
        <react_router_1.Link to="/">
          <logo_1.Logo className="scale-75 sm:scale-100" height={50}/>
        </react_router_1.Link>
      </div>
      <div className="ml-auto flex h-full items-center">
        <div className="mr-[30px]">{contestState}</div>
        <div className="bg-primary flex h-full w-[140px] items-center justify-end pt-[3px] [clip-path:polygon(40%_0,100%_0,100%_100%,0_100%)]">
          <div className="mr-[30px]">{accountMenu}</div>
        </div>
      </div>
    </div>);
}
