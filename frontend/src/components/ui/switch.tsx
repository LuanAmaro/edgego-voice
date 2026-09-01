"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface SwitchProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  checked?: boolean;
  defaultChecked?: boolean;
  defaultSelected?: boolean;
  onCheckedChange?: (checked: boolean) => void;
}

const Switch = React.forwardRef<HTMLButtonElement, SwitchProps>(
  (
    {
      className,
      checked,
      defaultChecked,
      defaultSelected,
      onCheckedChange,
      onClick,
      ...props
    },
    ref
  ) => {
    const [internalChecked, setInternalChecked] = React.useState(
      checked !== undefined
        ? checked
        : defaultChecked !== undefined
        ? defaultChecked
        : defaultSelected !== undefined
        ? defaultSelected
        : false
    );

    const isChecked = checked !== undefined ? checked : internalChecked;

    const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
      onClick?.(e);
      const nextChecked = !isChecked;
      if (checked === undefined) {
        setInternalChecked(nextChecked);
      }
      onCheckedChange?.(nextChecked);
    };

    return (
      <button
        type="button"
        role="switch"
        aria-checked={isChecked}
        data-state={isChecked ? "checked" : "unchecked"}
        ref={ref}
        onClick={handleClick}
        className={cn(
          "peer inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50",
          isChecked ? "bg-primary" : "bg-zinc-800",
          className
        )}
        {...props}
      >
        <span
          data-state={isChecked ? "checked" : "unchecked"}
          className={cn(
            "pointer-events-none block h-4 w-4 rounded-full bg-white shadow-md ring-0 transition-transform duration-150 ease-in-out",
            isChecked ? "translate-x-4" : "translate-x-0"
          )}
        />
      </button>
    );
  }
);
Switch.displayName = "Switch";

export { Switch };
