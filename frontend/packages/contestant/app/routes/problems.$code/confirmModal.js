"use strict";
var __rest = (this && this.__rest) || function (s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConfirmModal = ConfirmModal;
var react_1 = require("react");
var react_2 = require("@headlessui/react");
var material_symbol_1 = require("@app/components/material-symbol");
function ConfirmModal(_a) {
    var isOpen = _a.isOpen, onConfirm = _a.onConfirm, onCancel = _a.onCancel, title = _a.title, confirmText = _a.confirmText, cancelText = _a.cancelText, dialogClassName = _a.dialogClassName, children = _a.children, props = __rest(_a, ["isOpen", "onConfirm", "onCancel", "title", "confirmText", "cancelText", "dialogClassName", "children"]);
    return (<react_2.Transition appear show={isOpen} as={react_1.Fragment}>
      <react_2.Dialog as="div" className="relative z-50" onClose={onCancel}>
        <react_2.DialogBackdrop transition className="bg-disabled/30 fixed inset-0 duration-300 ease-out data-[closed]:opacity-0"/>
        <react_2.TransitionChild as={react_1.Fragment} enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100" leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0">
          <div className="fixed inset-0 bg-black/25"/>
        </react_2.TransitionChild>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <react_2.TransitionChild as={react_1.Fragment} enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100" leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95">
              <react_2.DialogPanel className={dialogClassName}>
                <div className="text-primary mb-4 flex items-center gap-4">
                  <material_symbol_1.MaterialSymbol icon="error" size={24} fill/>
                  <react_2.DialogTitle as="h3" className="text-20 text-primary font-bold">
                    {title}
                  </react_2.DialogTitle>
                </div>
                {children}

                <div className="mt-4 flex justify-end gap-8">
                  <react_2.Button className="border-text bg-surface-0 text-16 text-text hover:bg-surface-1 inline-flex justify-center rounded-[6px] border px-8 py-2 font-medium transition" onClick={onCancel}>
                    <div className="py-4">{cancelText}</div>
                  </react_2.Button>
                  <react_2.Button className="bg-surface-2 text-16 text-text inline-flex justify-center rounded-[6px] border border-transparent px-8 py-2 font-medium transition hover:opacity-80" form={props.formId} type={props.confirmType} name={props.confirmName} value={props.confirmValue} onClick={function () {
            onConfirm === null || onConfirm === void 0 ? void 0 : onConfirm();
            onCancel();
        }}>
                    <div className="py-4">{confirmText}</div>
                  </react_2.Button>
                </div>
              </react_2.DialogPanel>
            </react_2.TransitionChild>
          </div>
        </div>
      </react_2.Dialog>
    </react_2.Transition>);
}
