"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SignUpPage = SignUpPage;
var react_dom_1 = require("react-dom");
var clsx_1 = require("clsx");
var react_1 = require("@headlessui/react");
var material_symbol_1 = require("../components/material-symbol");
var title_1 = require("../components/title");
function SignUpPage(props) {
    return (<>
      <title_1.Title>アカウント登録</title_1.Title>
      <div className="grid h-full items-center justify-center">
        <form className="rounded-16 flex w-96 flex-col p-64 shadow-lg" onSubmit={function (e) {
            e.preventDefault();
            var form = e.target;
            var formData = new FormData(form);
            props.submit({
                invitationCode: formData.get("invitation_code"),
                name: formData.get("screen_name"),
                displayName: formData.get("display_name"),
            });
        }}>
          <h2 className="text-24 font-bold">アカウント登録</h2>
          <div className="mt-16 grid grid-cols-2 gap-20">
            <TextField name="invitation_code" className="col-span-2" label="招待コード" placeholder="(メールをご確認ください)" ignoreCompletion errorMessage={props.invitationCodeError != null
            ? {
                required: "未入力です",
                invalid: "無効なコードです",
                team_full: "チームが満員です",
            }[props.invitationCodeError]
            : undefined}/>
            <TextField name="screen_name" autoComplete="username" label="ID" placeholder="ictsc_taro" defaultValue={props.defaultName} errorMessage={props.nameError != null
            ? {
                required: "未入力です",
                invalid: "内容に誤りがあります",
                duplicate: "既に使われています",
            }[props.nameError]
            : undefined}/>
            <TextField name="display_name" autoComplete="name" label="表示名" placeholder="ICTSC太郎" defaultValue={props.defaultDisplayName} errorMessage={props.displayNameError != null
            ? {
                required: "未入力です",
                invalid: "内容に誤りがあります",
            }[props.displayNameError]
            : undefined}/>
          </div>
          <div className="mt-64 flex items-center justify-end gap-24">
            {props.error && (<p className="text-16 text-primary font-bold">
                {{
                rate_limit: "リクエストが多すぎます",
                invalid: "正しく入力されていない項目があります",
                unknown: "エラーが発生しました",
            }[props.error]}
              </p>)}
            <Submit />
          </div>
        </form>
      </div>
    </>);
}
function TextField(props) {
    var pending = (0, react_dom_1.useFormStatus)().pending;
    return (<react_1.Field className={props.className} disabled={pending}>
      <div className="text-16 flex items-center gap-16 font-bold">
        <react_1.Label>{props.label}</react_1.Label>
        {props.errorMessage && (<react_1.Description className="text-primary">
            {props.errorMessage}
          </react_1.Description>)}
      </div>
      <react_1.Input name={props.name} type="text" className="rounded-12 bg-surface-2 mt-8 w-full px-12 py-8 transition" placeholder={props.placeholder} invalid={Boolean(props.errorMessage)} defaultValue={props.defaultValue} autoComplete={props.autoComplete} data-1p-ignore={props.ignoreCompletion}/>
    </react_1.Field>);
}
function Submit() {
    var pending = (0, react_dom_1.useFormStatus)().pending;
    return (
    // FIXME:
    // useFormStatus には配下のコンポーネントで状態更新が発生すると状態がリセットされるバグがある(FYI: https://github.com/facebook/react/issues/30368)
    // これを回避するため，headlessui の Button コンポーネントを使わずに button 要素を使う
    <button type="submit" className={(0, clsx_1.clsx)("group rounded-12 bg-surface-2 py-[14px] pr-[28px] pl-[36px] shadow-md transition", "hover:bg-surface-2/90 active:shadow-transparent")} disabled={pending}>
      <div className={(0, clsx_1.clsx)("flex items-center gap-8 group-disabled:opacity-50")}>
        <span className="text-16 font-bold">登録</span>
        <material_symbol_1.MaterialSymbol size={24} icon="send"/>
      </div>
    </button>);
}
