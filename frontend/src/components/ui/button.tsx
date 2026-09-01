import * as React from "react";
import { cn } from "@/lib/utils";

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "outline" | "destructive" | "ghost";
  size?: "default" | "sm" | "lg" | "icon";
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "primary", size = "default", ...props }, ref) => {
    const variantStyles = {
      primary:
        "bg-primary text-primary-foreground font-semibold hover:bg-primary-hover shadow-[0_0_15px_rgba(0,223,129,0.25)] hover:shadow-[0_0_22px_rgba(0,223,129,0.45)] hover:-translate-y-0.5",
      secondary:
        "bg-card hover:bg-card-hover text-foreground border border-border hover:border-zinc-700 hover:-translate-y-0.5 shadow-sm",
      outline:
        "bg-transparent hover:bg-white/5 text-foreground border border-border hover:border-zinc-600 hover:-translate-y-0.5",
      destructive:
        "bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/20 hover:text-white hover:-translate-y-0.5",
      ghost:
        "bg-transparent hover:bg-white/5 text-muted-foreground hover:text-foreground",
    };

    const sizeStyles = {
      default: "h-9 px-4 py-2 text-sm rounded-md",
      sm: "h-8 px-3 text-xs rounded-md",
      lg: "h-11 px-6 text-base rounded-md font-semibold",
      icon: "h-9 w-9 p-0 rounded-md flex items-center justify-center",
    };

    return (
      <button
        ref={ref}
        className={cn(
          "inline-flex items-center justify-center gap-2 font-medium select-none motion-btn focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary disabled:pointer-events-none disabled:opacity-50",
          variantStyles[variant],
          sizeStyles[size],
          className
        )}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

export { Button };
