import { Fragment, ReactNode } from "react";
import {
  Dialog,
  DialogBackdrop,
  DialogPanel,
  DialogTitle,
  Transition,
  TransitionChild,
} from "@headlessui/react";
import { MaterialSymbol } from "@app/components/material-symbol";

interface ConfirmModalProps {
  isOpen: boolean;
  onConfirm: () => void;
  onCansel: () => void;
  title: string;
  confirmText: string;
  cancelText: string;
  children?: ReactNode;
}

export function ConfirmModal({
  isOpen,
  onConfirm,
  onCansel,
  title,
  confirmText,
  cancelText,
  children,
}: ConfirmModalProps) {
  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onCansel}>
        <DialogBackdrop
          transition
          className="fixed inset-0 bg-disabled/30 duration-300 ease-out data-[closed]:opacity-0"
        />
        <TransitionChild
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="bg-black/25 fixed inset-0" />
        </TransitionChild>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <TransitionChild
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <DialogPanel className="w-full max-w-md transform rounded-8 bg-surface-0 p-16 text-left align-middle shadow-xl transition-all">
                <div className="mb-4 flex items-center gap-4 text-primary">
                  <MaterialSymbol icon="error" size={24} fill />
                  <DialogTitle
                    as="h3"
                    className="text-20 font-bold text-primary"
                  >
                    {title}
                  </DialogTitle>
                </div>
                {children}

                <div className="mt-4 flex justify-end gap-8">
                  <button
                    className="inline-flex justify-center rounded-[6px] border border-text bg-surface-0 px-8 py-2 text-16 font-medium text-text transition hover:bg-surface-1"
                    onClick={onCansel}
                  >
                    <div className="py-4">{cancelText}</div>
                  </button>
                  <button
                    className="inline-flex justify-center rounded-[6px] border border-transparent bg-surface-2 px-8 py-2 text-16 font-medium text-text transition hover:opacity-80"
                    onClick={() => {
                      onConfirm();
                      onCansel();
                    }}
                  >
                    <div className="py-4">{confirmText}</div>
                  </button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
