"use client";

import React from "react";
import { Card } from "@/components/ui/card";
import { useLanguage } from "@/i18n/LanguageContext";

interface StatsGridProps {
  totalPersonas: number;
  totalVoices: number;
}

export function StatsGrid({ totalPersonas, totalVoices }: StatsGridProps) {
  const { t } = useLanguage();

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-6">
      <Card className="p-4 bg-card/60 border-border/80 shadow-none">
        <span className="text-[11px] font-medium text-muted-foreground">
          {t("activePersonasStat")}
        </span>
        <div className="text-2xl font-bold text-foreground tracking-tight mt-0.5">
          {totalPersonas}
        </div>
      </Card>

      <Card className="p-4 bg-card/60 border-border/80 shadow-none">
        <span className="text-[11px] font-medium text-muted-foreground">
          {t("audioFormatsStat")}
        </span>
        <div className="text-2xl font-bold text-foreground tracking-tight mt-0.5">
          6
        </div>
      </Card>

      <Card className="p-4 bg-card/60 border-border/80 shadow-none">
        <span className="text-[11px] font-medium text-muted-foreground">
          {t("neuralVoicesStat")}
        </span>
        <div className="text-2xl font-bold text-foreground tracking-tight mt-0.5">
          {totalVoices}
        </div>
      </Card>

      <Card className="p-4 bg-card/60 border-border/80 shadow-none">
        <span className="text-[11px] font-medium text-muted-foreground">
          {t("engineCacheStat")}
        </span>
        <div className="text-sm font-semibold text-foreground tracking-tight mt-1.5 flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
          <span>{t("engineCacheValue")}</span>
        </div>
      </Card>
    </div>
  );
}

