import * as React from "react";
import { cn } from "@/lib/utils";

export function FieldGroup({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn("space-y-3", className)} {...props} />;
}

export function FieldLabel({
  className,
  ...props
}: React.LabelHTMLAttributes<HTMLLabelElement>) {
  return (
    <label
      className={cn("block cursor-pointer select-none", className)}
      {...props}
    />
  );
}

export function Field({
  className,
  orientation = "vertical",
  ...props
}: React.HTMLAttributes<HTMLDivElement> & {
  orientation?: "horizontal" | "vertical";
}) {
  return (
    <div
      className={cn(
        "p-3 rounded-lg border border-border bg-secondary/30 transition-colors hover:border-zinc-700",
        orientation === "horizontal" &&
          "flex items-center justify-between gap-3",
        className
      )}
      {...props}
    />
  );
}

export function FieldContent({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn("space-y-0.5 min-w-0 flex-1", className)} {...props} />;
}

export function FieldTitle({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h4
      className={cn("text-xs font-medium text-foreground tracking-tight", className)}
      {...props}
    />
  );
}

export function FieldDescription({
  className,
  ...props
}: React.HTMLAttributes<HTMLParagraphElement>) {
  return (
    <p
      className={cn("text-[11px] text-muted-foreground leading-relaxed", className)}
      {...props}
    />
  );
}
