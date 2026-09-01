import * as React from "react";
import { cn } from "@/lib/utils";
import { X } from "lucide-react";

export interface DialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: React.ReactNode;
}

export function Dialog({ open, onOpenChange, children }: DialogProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm animate-in fade-in-0">
      <div
        className="fixed inset-0"
        onClick={() => onOpenChange(false)}
      />
      <div className="relative z-50 w-full max-w-lg rounded-lg border border-border bg-card shadow-2xl overflow-hidden animate-in zoom-in-95">
        {children}
      </div>
    </div>
  );
}

export function DialogHeader({
  className,
  children,
  onClose,
  ...props
}: React.HTMLAttributes<HTMLDivElement> & { onClose?: () => void }) {
  return (
    <div
      className={cn(
        "flex items-center justify-between p-6 border-b border-border",
        className
      )}
      {...props}
    >
      <div className="flex flex-col space-y-1.5">{children}</div>
      {onClose && (
        <button
          onClick={onClose}
          className="rounded-md p-1.5 text-muted-foreground hover:text-foreground hover:bg-white/5 transition-colors"
        >
          <X className="h-4 w-4" />
        </button>
      )}
    </div>
  );
}

export function DialogTitle({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h2
      className={cn(
        "text-lg font-semibold leading-none tracking-tight text-foreground",
        className
      )}
      {...props}
    />
  );
}

export function DialogDescription({
  className,
  ...props
}: React.HTMLAttributes<HTMLParagraphElement>) {
  return (
    <p
      className={cn("text-sm text-muted-foreground", className)}
      {...props}
    />
  );
}

export function DialogContent({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn("p-6 space-y-4 max-h-[75vh] overflow-y-auto", className)} {...props} />;
}

export function DialogFooter({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "flex items-center justify-end gap-3 p-4 border-t border-border bg-white/[0.01]",
        className
      )}
      {...props}
    />
  );
}
