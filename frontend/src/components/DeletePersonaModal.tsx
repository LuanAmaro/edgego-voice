"use client";

import React, { useState, useEffect, useRef } from "react";
import { Persona } from "@/types";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogContent,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";
import { useLanguage } from "@/i18n/LanguageContext";

interface DeletePersonaModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  persona: Persona | null;
  onConfirm: (persona: Persona) => Promise<void>;
}

export function DeletePersonaModal({
  open,
  onOpenChange,
  persona,
  onConfirm,
}: DeletePersonaModalProps) {
  const { t } = useLanguage();
  const [inputValue, setInputValue] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const isMatch = persona ? inputValue.trim() === persona.id : false;

  useEffect(() => {
    if (open) {
      setInputValue("");
      setIsDeleting(false);
      const timer = setTimeout(() => inputRef.current?.focus(), 100);
      return () => clearTimeout(timer);
    }
  }, [open]);

  const handleConfirm = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!persona || !isMatch || isDeleting) return;

    setIsDeleting(true);
    try {
      await onConfirm(persona);
      onOpenChange(false);
    } finally {
      setIsDeleting(false);
    }
  };

  if (!persona) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <form onSubmit={handleConfirm}>
        <DialogHeader onClose={() => onOpenChange(false)}>
          <DialogTitle>{t("deleteModalTitle")}</DialogTitle>
          <DialogDescription>
            {t("deleteModalDescription")}{" "}
            <span className="font-medium text-foreground">{persona.name}</span>.
          </DialogDescription>
        </DialogHeader>

        <DialogContent className="space-y-3">
          <label
            htmlFor="confirm-slug-input"
            className="text-xs text-muted-foreground leading-relaxed block"
          >
            {t("deleteModalConfirmLabel")}{" "}
            <span className="font-mono font-semibold text-foreground bg-secondary/80 px-1.5 py-0.5 rounded border border-border/60 select-all">
              {persona.id}
            </span>{" "}
            {t("deleteModalInField")}
          </label>

          <Input
            id="confirm-slug-input"
            ref={inputRef}
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder={persona.id}
            className="font-mono text-xs h-9"
            autoComplete="off"
            spellCheck={false}
          />
        </DialogContent>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isDeleting}
            className="h-8 text-xs"
          >
            {t("cancel")}
          </Button>
          <Button
            type="submit"
            variant="destructive"
            disabled={!isMatch || isDeleting}
            className="h-8 text-xs font-medium"
          >
            {isDeleting ? (
              <>
                <IonIcon
                  name="sync-outline"
                  className="text-xs animate-spin mr-1.5"
                />
                <span>{t("deletingButton")}</span>
              </>
            ) : (
              <span>{t("deletePermanentButton")}</span>
            )}
          </Button>
        </DialogFooter>
      </form>
    </Dialog>
  );
}



