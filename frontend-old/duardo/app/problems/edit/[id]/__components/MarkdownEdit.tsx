import React, { useEffect, useState } from "react";

import { UseFormRegister } from "react-hook-form";

import { FormSchema } from "@/app/problems/edit/[id]/__type/MarkdownEditFormSchema";

type Props = {
  register: UseFormRegister<FormSchema>;
  onChange?: (value: string) => void;
};

export default function MarkdownEdit({ register, onChange }: Props) {
  const [text, setText] = useState("");

  useEffect(() => {
    const textarea: HTMLTextAreaElement | null = document.querySelector(
      ".auto-expand-textarea",
    );
    if (textarea) {
      textarea.style.height = "inherit";
      textarea.style.height = `${textarea.scrollHeight}px`;
    }
  }, [text]);

  function handleChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setText(event.target.value);
    if (onChange) {
      onChange(event.target.value);
    }
  }

  return (
    <textarea
      className="textarea auto-expand-textarea focus:outline-none focus:border-0 border-0 w-full resize-none"
      value={text}
      {...register("body", { required: true })}
      onChange={handleChange}
    />
  );
}
