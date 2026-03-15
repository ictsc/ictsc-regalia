"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SignInPage = SignInPage;
var react_1 = require("@headlessui/react");
var icon_clyde_white_RGB_svg_1 = require("../../assets/icon_clyde_white_RGB.svg");
var logo_1 = require("../components/logo");
var title_1 = require("../components/title");
function SignInPage(_a) {
    var signInURL = _a.signInURL;
    return (<>
      <title_1.Title>ログイン</title_1.Title>
      <div className="mx-40 flex h-full flex-col items-center justify-center gap-[90px]">
        <logo_1.Logo width={500}/>
        <DiscordLoginButton href={signInURL}/>
      </div>
    </>);
}
function DiscordLoginButton(props) {
    return (<react_1.Button as="a" className="rounded-16 text-32 bg-[#5865f2] py-[22px] ps-16 pe-[20px] shadow-md disabled:bg-[#a0a0a0] data-[hover]:bg-[#4752c4]" {...props}>
      <span className="text-surface-0 flex flex-row gap-[12px]">
        <img src={icon_clyde_white_RGB_svg_1.default} width={40} height={40} alt=""/>
        Discord でログイン
      </span>
    </react_1.Button>);
}
