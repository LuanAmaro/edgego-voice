import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["class"],
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--bg-page)",
        foreground: "var(--text-main)",
        card: {
          DEFAULT: "var(--bg-card)",
          foreground: "var(--text-main)",
          hover: "var(--bg-card-hover)",
        },
        muted: {
          DEFAULT: "var(--bg-muted)",
          foreground: "var(--text-muted)",
        },
        primary: {
          DEFAULT: "var(--primary)",
          hover: "var(--primary-hover)",
          foreground: "var(--primary-foreground)",
        },
        border: "var(--border-subtle)",
        input: "var(--bg-input)",
        ring: "var(--primary)",
        destructive: {
          DEFAULT: "var(--color-danger)",
          foreground: "var(--text-main)",
        },
      },
      borderRadius: {
        lg: "12px",
        md: "10px",
        sm: "6px",
      },
      fontFamily: {
        sans: ["var(--font-inter)", "Inter", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
    },
  },
  plugins: [],
};

export default config;
