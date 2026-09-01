"use client";

import React from "react";
import { Persona, VoiceOption } from "@/types";
import { StatsGrid } from "@/components/StatsGrid";
import { PersonaCard } from "@/components/PersonaCard";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";
import { useLanguage } from "@/i18n/LanguageContext";

interface PersonasListProps {
  personas: Persona[];
  voices: VoiceOption[];
  onOpenNewModal: () => void;
  onTest: (persona: Persona) => void;
  onSnippet: (persona: Persona) => void;
  onEdit: (persona: Persona) => void;
  onDelete: (persona: Persona) => void;
}

export function PersonasList({
  personas,
  voices,
  onOpenNewModal,
  onTest,
  onSnippet,
  onEdit,
  onDelete,
}: PersonasListProps) {
  const { t } = useLanguage();

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <StatsGrid
        totalPersonas={personas.length}
        totalVoices={voices.length}
      />

      {/* Section Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-foreground tracking-tight">
            {t("personasTitle")}
          </h2>
          <p className="text-sm text-muted-foreground">
            {t("personasSubtitle")}
          </p>
        </div>

        <Button onClick={onOpenNewModal} className="h-10">
          <IonIcon name="add-outline" className="text-lg" />
          <span>{t("newPersonaButton")}</span>
        </Button>
      </div>

      {/* Personas Grid */}
      {personas.length === 0 ? (
        <div className="rounded-xl border border-dashed border-border bg-card p-12 text-center flex flex-col items-center justify-center animate-card-enter">
          <div className="text-3xl text-muted-foreground mb-3">
            <IonIcon name="people-outline" />
          </div>
          <h3 className="text-base font-semibold text-foreground mb-1">
            {t("emptyPersonasTitle")}
          </h3>
          <p className="text-xs text-muted-foreground mb-4">
            {t("emptyPersonasDesc")}
          </p>
          <Button onClick={onOpenNewModal}>
            <IonIcon name="add-outline" className="text-lg" />
            <span>{t("createPersonaButton")}</span>
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {personas.map((persona, index) => (
            <PersonaCard
              key={persona.id}
              persona={persona}
              voices={voices}
              index={index}
              onTest={onTest}
              onSnippet={onSnippet}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </div>
      )}
    </div>
  );
}

