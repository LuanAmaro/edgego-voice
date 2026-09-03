"use client";

import React from "react";
import { Persona, VoiceOption } from "@/types";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";
import { CountryFlag } from "@/components/ui/country-flag";
import { useLanguage } from "@/i18n/LanguageContext";

interface PersonaCardProps {
  persona: Persona;
  voices: VoiceOption[];
  index?: number;
  onTest: (persona: Persona) => void;
  onSnippet: (persona: Persona) => void;
  onEdit: (persona: Persona) => void;
  onDelete: (persona: Persona) => void;
}

export function PersonaCard({
  persona,
  voices,
  onTest,
  onSnippet,
  onEdit,
  onDelete,
}: PersonaCardProps) {
  const { t } = useLanguage();
  const voiceObj = voices.find((v) => v.id === persona.voice);
  const countryCode = voiceObj?.country_code || "BR";
  const voiceName = voiceObj ? voiceObj.name : persona.voice;

  return (
    <div className="rounded-lg border border-border/80 bg-card/60 p-4 flex flex-col gap-3 transition-colors hover:border-border">
      {/* Header */}
      <div className="flex items-start justify-between gap-2">
        <div>
          <h4 className="font-semibold text-sm text-foreground">
            {persona.name}
          </h4>
          <span className="text-[11px] font-mono text-muted-foreground">
            {persona.id}
          </span>
        </div>

        <span className="inline-flex items-center gap-1 text-[10px] font-mono text-muted-foreground bg-secondary px-1.5 py-0.5 rounded">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
          <span>{t("activeStatus")}</span>
        </span>
      </div>

      {/* Description */}
      <p className="text-xs text-muted-foreground line-clamp-2 min-h-[32px]">
        {persona.description || t("noDescription")}
      </p>

      {/* Attributes */}
      <div className="rounded-md border border-border/60 bg-background/50 p-2.5 space-y-1.5 text-xs">
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-[11px]">{t("voiceLabel")}</span>
          <span className="font-medium text-foreground flex items-center gap-1.5 text-[11px]">
            <CountryFlag
              countryCode={countryCode}
              className="w-3.5 h-2.5 rounded-[2px]"
            />
            <span className="truncate max-w-[140px]">{voiceName}</span>
          </span>
        </div>

        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-[11px]">{t("formatLabel")}</span>
          <span className="font-mono text-[10px] uppercase text-muted-foreground bg-secondary px-1.5 py-0.2 rounded">
            {persona.format}
          </span>
        </div>

        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-[11px]">{t("speedLabel")}</span>
          <span className="font-mono text-[11px] text-foreground">
            {persona.speed.toFixed(2)}x
          </span>
        </div>

        {persona.pitch && persona.pitch !== "+0Hz" && (
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground text-[11px]">{t("pitchLabel")}</span>
            <span className="font-mono text-[10px] text-primary bg-primary/10 px-1.5 py-0.5 rounded font-semibold">
              {persona.pitch}
            </span>
          </div>
        )}

        {(persona.break_comma || persona.break_period) && (
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground text-[11px]">{t("ssmlBreakSectionTitle")}</span>
            <span className="font-mono text-[10px] text-muted-foreground bg-secondary px-1.5 py-0.5 rounded">
              {persona.break_comma || "0ms"} / {persona.break_period || "0ms"}
            </span>
          </div>
        )}
      </div>

      {/* Actions */}
      <div className="flex items-center gap-1.5 pt-2 border-t border-border/40 mt-auto">
        <Button
          variant="secondary"
          size="sm"
          onClick={() => onTest(persona)}
          className="h-7 text-xs px-2.5 font-normal"
        >
          <IonIcon name="play-outline" className="text-xs mr-1" />
          <span>{t("testAction")}</span>
        </Button>

        <Button
          variant="outline"
          size="sm"
          onClick={() => onSnippet(persona)}
          className="h-7 text-xs px-2.5 font-normal"
        >
          <IonIcon name="code-slash-outline" className="text-xs mr-1" />
          <span>{t("apiAction")}</span>
        </Button>

        <div className="flex items-center gap-1 ml-auto">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onEdit(persona)}
            title={t("editPersonaTitle")}
            className="h-7 w-7 text-muted-foreground hover:text-foreground"
          >
            <IonIcon name="create-outline" className="text-xs" />
          </Button>

          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(persona)}
            title={t("deletePersonaTitle")}
            className="h-7 w-7 text-muted-foreground hover:text-destructive"
          >
            <IonIcon name="trash-outline" className="text-xs" />
          </Button>
        </div>
      </div>
    </div>
  );
}

