"use client";

import React, { useState, useRef, useEffect } from "react";
import { VoiceOption } from "@/types";
import { CountryFlag } from "@/components/ui/country-flag";
import { IonIcon } from "@/components/ui/ion-icon";

interface VoiceSelectProps {
  voices: VoiceOption[];
  value: string;
  onChange: (value: string) => void;
  className?: string;
}

export function VoiceSelect({
  voices,
  value,
  onChange,
  className = "",
}: VoiceSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);
  const searchRef = useRef<HTMLInputElement>(null);

  const selectedVoice = voices.find((v) => v.id === value) || voices[0];

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
        setSearch("");
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => searchRef.current?.focus(), 50);
    }
  }, [isOpen]);

  const filteredVoices = voices.filter(
    (v) =>
      v.name.toLowerCase().includes(search.toLowerCase()) ||
      v.language.toLowerCase().includes(search.toLowerCase()) ||
      (v.country_code &&
        v.country_code.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className={`relative ${className}`} ref={containerRef}>
      {/* Trigger */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className={`flex h-auto w-full items-center justify-between rounded-xl border px-4 py-3 text-sm text-foreground transition-all duration-200 shadow-sm ${
          isOpen
            ? "border-primary bg-card ring-1 ring-primary/30"
            : "border-border bg-card hover:border-zinc-600"
        }`}
      >
        <div className="flex items-center gap-3 truncate">
          {selectedVoice && (
            <CountryFlag
              countryCode={selectedVoice.country_code || "BR"}
              className="w-6 h-4 rounded-[3px] shadow-sm shrink-0"
            />
          )}
          <div className="text-left truncate">
            <span className="font-bold text-foreground">
              {selectedVoice?.name}
            </span>
            <span className="text-xs text-muted-foreground ml-2">
              ({selectedVoice?.language} - {selectedVoice?.gender})
            </span>
          </div>
        </div>
        <IonIcon
          name="chevron-down-outline"
          className={`text-sm text-muted-foreground transition-transform duration-200 ml-2 shrink-0 ${
            isOpen ? "rotate-180 text-primary" : ""
          }`}
        />
      </button>

      {/* Dropdown */}
      {isOpen && (
        <div className="absolute left-0 right-0 z-50 mt-2 overflow-hidden rounded-xl border border-border bg-card shadow-2xl flex flex-col animate-in fade-in-0 slide-in-from-top-2 duration-150">
          {/* Search */}
          <div className="p-2.5 border-b border-border/60">
            <div className="relative flex items-center">
              <IonIcon
                name="search-outline"
                className="absolute left-3 text-xs text-muted-foreground pointer-events-none"
              />
              <input
                ref={searchRef}
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Buscar voz ou país (ex: Brasil, EUA, Espanhol)..."
                className="w-full bg-background/50 rounded-lg border border-border py-2 pl-8 pr-3 text-xs text-foreground placeholder:text-muted-foreground/60 focus:outline-none focus:border-primary transition-colors"
              />
            </div>
          </div>

          {/* List */}
          <div className="overflow-y-auto p-1.5 space-y-0.5 max-h-64 scrollbar-thin">
            {filteredVoices.map((v) => {
              const isSelected = v.id === value;
              return (
                <button
                  key={v.id}
                  type="button"
                  onClick={() => {
                    onChange(v.id);
                    setIsOpen(false);
                    setSearch("");
                  }}
                  className={`flex w-full items-center justify-between rounded-lg px-3 py-2.5 text-left text-sm transition-colors ${
                    isSelected
                      ? "bg-primary/10 border border-primary/30 text-primary"
                      : "text-foreground hover:bg-white/5 border border-transparent"
                  }`}
                >
                  <div className="flex items-center gap-3 truncate">
                    <CountryFlag
                      countryCode={v.country_code || "BR"}
                      className="w-6 h-4 rounded-[3px] shadow-sm shrink-0"
                    />
                    <div className="truncate">
                      <div className="flex items-center gap-1.5">
                        <span className={`font-semibold text-sm ${isSelected ? "text-primary" : "text-foreground"}`}>
                          {v.name}
                        </span>
                        <span className="text-[11px] text-zinc-500 font-normal">
                          • {v.gender}
                        </span>
                      </div>
                      <div className={`text-[11px] ${isSelected ? "text-primary/70" : "text-muted-foreground"}`}>
                        {v.language}
                      </div>
                    </div>
                  </div>

                  {isSelected && (
                    <IonIcon
                      name="checkmark-outline"
                      className="text-base text-primary ml-2 shrink-0"
                    />
                  )}
                </button>
              );
            })}

            {filteredVoices.length === 0 && (
              <div className="p-6 text-center text-xs text-muted-foreground">
                <IonIcon name="search-outline" className="text-2xl mb-2 block mx-auto opacity-30" />
                Nenhuma voz encontrada com esse filtro.
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
