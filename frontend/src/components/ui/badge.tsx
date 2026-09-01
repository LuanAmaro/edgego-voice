import * as React from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "success" | "secondary" | "outline" | "format";
}

export function Badge({
  className,
  variant = "default",
  ...props
}: BadgeProps) {
  const variantStyles = {
    default: "bg-primary/10 text-primary border-primary/20",
    success:
      "bg-emerald-500/10 text-emerald-400 border-emerald-500/20 shadow-[0_0_8px_rgba(0,223,129,0.15)]",
    secondary: "bg-muted text-foreground border-border",
    outline: "border-border text-foreground bg-transparent",
    format: "bg-white/5 text-foreground border-border font-mono uppercase text-[11px]",
  };

  return (
    <div
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors select-none",
        variantStyles[variant],
        className
      )}
      {...props}
    />
  );
}
