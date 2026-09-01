import * as React from "react";
import { cn } from "@/lib/utils";

export interface AlertProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "destructive" | "success" | "info";
}

const Alert = React.forwardRef<HTMLDivElement, AlertProps>(
  ({ className, variant = "default", ...props }, ref) => {
    const variantStyles = {
      default: "bg-secondary/40 text-foreground border-border",
      destructive:
        "border-red-500/30 text-red-400 bg-red-500/10 [&>ion-icon]:text-red-400",
      success:
        "border-emerald-500/30 text-emerald-400 bg-emerald-500/10 [&>ion-icon]:text-emerald-400",
      info: "border-primary/30 text-foreground bg-primary/5 [&>ion-icon]:text-primary",
    };

    return (
      <div
        ref={ref}
        role="alert"
        className={cn(
          "relative w-full rounded-lg border p-3 text-xs flex items-start gap-2.5 transition-all",
          variantStyles[variant],
          className
        )}
        {...props}
      />
    );
  }
);
Alert.displayName = "Alert";

const AlertTitle = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h5
    ref={ref}
    className={cn("font-medium text-xs leading-none tracking-tight mb-1 text-foreground", className)}
    {...props}
  />
));
AlertTitle.displayName = "AlertTitle";

const AlertDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("text-[11px] text-muted-foreground leading-relaxed", className)}
    {...props}
  />
));
AlertDescription.displayName = "AlertDescription";

export { Alert, AlertTitle, AlertDescription };
