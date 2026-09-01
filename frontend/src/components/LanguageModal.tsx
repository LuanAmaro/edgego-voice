"use client";

import React from "react";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogContent,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { CountryFlag } from "@/components/ui/country-flag";
import { IonIcon } from "@/components/ui/ion-icon";
import { useLanguage } from "@/i18n/LanguageContext";
import { LanguageCode } from "@/i18n/translations";
import { cn } from "@/lib/utils";

interface LanguageModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function LanguageModal({ open, onOpenChange }: LanguageModalProps) {
  const { language, setLanguage, languages, t } = useLanguage();

  const handleSelect = (code: LanguageCode) => {
    setLanguage(code);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <div className="w-full">
        <DialogHeader onClose={() => onOpenChange(false)}>
          <div className="flex items-center gap-2">
            <IonIcon name="globe-outline" className="text-primary text-lg" />
            <DialogTitle>{t("selectLanguage")}</DialogTitle>
          </div>
          <DialogDescription>
            {t("changeLanguageDescription")}
          </DialogDescription>
        </DialogHeader>

        <DialogContent className="space-y-2.5 pt-2">
          {languages.map((item) => {
            const isSelected = language === item.code;
            return (
              <button
                key={item.code}
                type="button"
                onClick={() => handleSelect(item.code)}
                className={cn(
                  "w-full flex items-center justify-between p-3 rounded-lg border transition-all text-left",
                  "hover:border-zinc-700 hover:bg-secondary/60",
                  isSelected
                    ? "border-primary/60 bg-primary/10 shadow-[0_0_15px_rgba(0,223,129,0.08)]"
                    : "border-border/80 bg-card/60"
                )}
              >
                <div className="flex items-center gap-3.5">
                  <div className="flex items-center justify-center w-8 h-6 rounded overflow-hidden shadow-sm border border-white/10 shrink-0">
                    <CountryFlag
                      countryCode={item.countryCode}
                      className="w-full h-full object-cover"
                    />
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-foreground flex items-center gap-2">
                      <span>{item.nativeName}</span>
                      {isSelected && (
                        <span className="text-[10px] uppercase font-mono px-1.5 py-0.2 rounded bg-primary/20 text-primary font-medium">
                          {item.code}
                        </span>
                      )}
                    </div>
                    <span className="text-xs text-muted-foreground">
                      {item.name} • {item.flagTitle}
                    </span>
                  </div>
                </div>

                <div
                  className={cn(
                    "w-5 h-5 rounded-full flex items-center justify-center border transition-all shrink-0",
                    isSelected
                      ? "border-primary bg-primary text-primary-foreground"
                      : "border-muted-foreground/30 bg-transparent"
                  )}
                >
                  {isSelected && (
                    <IonIcon name="checkmark-outline" className="text-xs stroke-[3]" />
                  )}
                </div>
              </button>
            );
          })}
        </DialogContent>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="h-8 text-xs"
          >
            {t("close")}
          </Button>
        </DialogFooter>
      </div>
    </Dialog>
  );
}
