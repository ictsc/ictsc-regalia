import { Fragment, type ReactNode } from "react";
import {
  Button,
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

  onConfirm?: () => void;
  formId?: string;
  confirmType?: "button" | "submit";
  confirmName?: string;
  confirmValue?: string;

  onCancel: () => void;
  title: string;
  confirmText: string;
  cancelText: string;
  dialogClassName?: string;
  children?: ReactNode;
}

export function ConfirmModal({
  isOpen,
  onConfirm,
  onCancel,
  title,
  confirmText,
  cancelText,
  dialogClassName,
  children,
  ...props
}: ConfirmModalProps) {
  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onCancel}>
        <DialogBackdrop
          transition
          className="bg-disabled/30 fixed inset-0 duration-300 ease-out data-[closed]:opacity-0"
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
          <div className="fixed inset-0 bg-black/25" />
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
              <DialogPanel className={dialogClassName}>
                <div className="text-primary mb-4 flex items-center gap-4">
                  <MaterialSymbol icon="error" size={24} fill />
                  <DialogTitle
                    as="h3"
                    className="text-20 text-primary font-bold"
                  >
                    {title}
                  </DialogTitle>
                </div>
                {children}

                <div className="mt-4 flex justify-end gap-8">
                  <Button
                    className="border-text bg-surface-0 text-16 text-text hover:bg-surface-1 inline-flex justify-center rounded-[6px] border px-8 py-2 font-medium transition"
                    onClick={onCancel}
                  >
                    <div className="py-4">{cancelText}</div>
                  </Button>
                  <Button
                    className="bg-surface-2 text-16 text-text inline-flex justify-center rounded-[6px] border border-transparent px-8 py-2 font-medium transition hover:opacity-80"
                    form={props.formId}
                    type={props.confirmType}
                    name={props.confirmName}
                    value={props.confirmValue}
                    onClick={() => {
                      onConfirm?.();
                      onCancel();
                    }}
                  >
                    <div className="py-4">{confirmText}</div>
                  </Button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}
