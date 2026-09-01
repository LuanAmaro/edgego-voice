"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { IonIcon } from "@/components/ui/ion-icon";
import { CountryFlag } from "@/components/ui/country-flag";
import { LanguageModal } from "@/components/LanguageModal";
import { useLanguage } from "@/i18n/LanguageContext";

interface LoginViewProps {
  serverUrl: string;
  setServerUrl: (url: string) => void;
  onLogin: (apiKey: string, serverUrl: string) => Promise<void>;
  isLoading: boolean;
}

export function LoginView({
  serverUrl,
  setServerUrl,
  onLogin,
  isLoading,
}: LoginViewProps) {
  const { currentLanguage, t } = useLanguage();
  const [apiKey, setApiKey] = useState("");
  const [isLangModalOpen, setIsLangModalOpen] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!apiKey.trim()) return;
    onLogin(apiKey.trim(), serverUrl.trim().replace(/\/+$/, ""));
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative">
      {/* Top right language switch on login page */}
      <div className="absolute top-4 right-4 z-10">
        <button
          onClick={() => setIsLangModalOpen(true)}
          className="flex items-center gap-2 px-2.5 py-1 rounded-md border border-border/80 bg-card/60 hover:bg-secondary hover:border-zinc-700 transition-all text-xs text-foreground group shadow-sm"
          title={t("selectLanguage")}
        >
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
      </div>

      <div className="w-full max-w-md bg-card border border-border rounded-xl p-8 shadow-2xl text-center relative overflow-hidden animate-card-enter">
        {/* Top Accent Line */}
        <div className="absolute top-0 left-0 right-0 h-[2px] bg-gradient-to-r from-transparent via-primary to-transparent" />

        {/* Official Logo Banner */}
        <div className="flex items-center justify-center mb-6">
          <img
            src="/logo.svg"
            alt="EdgeGo Voice Logo"
            className="h-14 w-auto object-contain hover:scale-105 transition-transform duration-200"
          />
        </div>

        <h1 className="text-xl font-bold tracking-tight text-foreground mb-1 font-sans">
          {t("loginTitle")}
        </h1>
        <p className="text-xs text-muted-foreground mb-6">
          {t("loginSubtitle")}
        </p>

        <form onSubmit={handleSubmit} className="space-y-4 text-left">
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              {t("serverUrlLabel")}
            </label>
            <Input
              type="text"
              value={serverUrl}
              onChange={(e) => setServerUrl(e.target.value)}
              placeholder="http://localhost:5050"
              icon={<IonIcon name="server-outline" className="text-base" />}
              className="font-mono text-xs"
              required
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              {t("apiKeyLabel")}
            </label>
            <Input
              type="password"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder="••••••••••••"
              icon={<IonIcon name="key-outline" className="text-base" />}
              required
            />
          </div>

          <Button
            type="submit"
            disabled={isLoading}
            className="w-full h-11 text-base mt-2"
          >
            <IonIcon name="log-in-outline" className="text-lg mr-1" />
            <span>{isLoading ? t("loggingIn") : t("loginButton")}</span>
          </Button>
        </form>
      </div>

      <LanguageModal
        open={isLangModalOpen}
        onOpenChange={setIsLangModalOpen}
      />
    </div>
  );
}

