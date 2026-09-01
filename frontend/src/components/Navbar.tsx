"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";
import { CountryFlag } from "@/components/ui/country-flag";
import { LanguageModal } from "@/components/LanguageModal";
import { useLanguage } from "@/i18n/LanguageContext";
import { TabType } from "@/types";

interface NavbarProps {
  activeTab: TabType;
  setActiveTab: (tab: TabType) => void;
  personasCount: number;
  onLogout: () => void;
}

export function Navbar({
  activeTab,
  setActiveTab,
  personasCount,
  onLogout,
}: NavbarProps) {
  const { currentLanguage, t } = useLanguage();
  const [isLangModalOpen, setIsLangModalOpen] = useState(false);

  return (
    <>
      <header className="h-14 bg-background/95 backdrop-blur border-b border-border sticky top-0 z-40 px-6 flex items-center justify-between">
        {/* Brand with Official Logo & v1.0 Badge */}
        <div className="flex items-center gap-3">
          <img
            src="/logo.svg"
            alt="EdgeGo Voice"
            className="h-6 w-auto object-contain"
          />
          <span className="text-[11px] text-muted-foreground font-mono bg-secondary px-1.5 py-0.5 rounded border border-border">
            v1.0
          </span>
        </div>

        {/* Minimalist Tabs Navigation */}
        <nav className="flex items-center gap-1 bg-secondary/50 border border-border/60 rounded-lg p-0.5">
          <button
            onClick={() => setActiveTab("personas")}
            className={`flex items-center gap-2 px-3 py-1.5 text-xs font-medium rounded-md transition-all ${
              activeTab === "personas"
                ? "bg-background text-foreground shadow-sm border border-border/80"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <IonIcon name="people-outline" className="text-sm" />
            <span>{t("personas")}</span>
            <span className="text-[10px] text-muted-foreground font-mono">
              {personasCount}
            </span>
          </button>

          <button
            onClick={() => setActiveTab("playground")}
            className={`flex items-center gap-2 px-3 py-1.5 text-xs font-medium rounded-md transition-all ${
              activeTab === "playground"
                ? "bg-background text-foreground shadow-sm border border-border/80"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <IonIcon name="sparkles-outline" className="text-sm" />
            <span>{t("playground")}</span>
          </button>

          <button
            onClick={() => setActiveTab("apidocs")}
            className={`flex items-center gap-2 px-3 py-1.5 text-xs font-medium rounded-md transition-all ${
              activeTab === "apidocs"
                ? "bg-background text-foreground shadow-sm border border-border/80"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <IonIcon name="code-slash-outline" className="text-sm" />
            <span>{t("apiDocs")}</span>
          </button>
        </nav>

        {/* Language Switcher & Logout */}
        <div className="flex items-center gap-2">
          {/* Language Selector Button (replaces static online indicator) */}
          <button
            onClick={() => setIsLangModalOpen(true)}
            className="flex items-center gap-2 px-2.5 py-1 rounded-md border border-border/80 bg-secondary/40 hover:bg-secondary hover:border-zinc-700 transition-all text-xs text-foreground group"
            title={t("selectLanguage")}
          >
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0" />
            <CountryFlag
              countryCode={currentLanguage.countryCode}
              className="w-4 h-2.5 rounded-[2px] shadow-sm shrink-0"
            />
            <span className="font-medium text-[11px] uppercase tracking-wide font-mono">
              {currentLanguage.code.split("-")[0]}
            </span>
            <IonIcon
              name="chevron-down-outline"
              className="text-[10px] text-muted-foreground group-hover:text-foreground transition-colors"
            />
          </button>

          <Button
            variant="ghost"
            size="icon"
            onClick={onLogout}
            title={t("logout")}
            className="h-8 w-8 text-muted-foreground hover:text-foreground"
          >
            <IonIcon name="log-out-outline" className="text-sm" />
          </Button>
        </div>
      </header>

      {/* Language Modal */}
      <LanguageModal
        open={isLangModalOpen}
        onOpenChange={setIsLangModalOpen}
      />
    </>
  );
}

