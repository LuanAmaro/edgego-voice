"use client";

import React, { useState } from "react";
import { Persona } from "@/types";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogContent,
  DialogFooter,
} from "@/components/ui/dialog";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";
import { useLanguage } from "@/i18n/LanguageContext";

interface SnippetModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  persona: Persona | null;
  serverUrl: string;
  apiKey: string;
}

export function SnippetModal({
  open,
  onOpenChange,
  persona,
  serverUrl,
  apiKey,
}: SnippetModalProps) {
  const { t } = useLanguage();
  const [copied, setCopied] = useState(false);

  if (!persona) return null;

  const curlSnippet = `curl -X POST "${serverUrl}/v1/persona/${persona.id}/speech" \\
  -H "Authorization: Bearer ${apiKey}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "input": "Olá! Mensagem gerada pela persona ${persona.name}."
  }' \\
  --output audio_${persona.id}.${persona.format}`;

  const handleCopy = () => {
    navigator.clipboard.writeText(curlSnippet);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogHeader onClose={() => onOpenChange(false)}>
        <DialogTitle>{t("snippetModalTitle")}</DialogTitle>
        <DialogDescription>
          {persona.name} ({persona.id})
        </DialogDescription>
      </DialogHeader>

      <DialogContent>
        <Card className="overflow-hidden border-border bg-[#0C0A09]">
          <div className="flex items-center justify-between px-4 py-2 bg-white/[0.03] border-b border-border">
            <span className="text-xs font-semibold text-foreground">
              {t("curlExampleTitle")}
            </span>
            <Button
              variant="secondary"
              size="sm"
              onClick={handleCopy}
              className="h-7 text-xs"
            >
              {copied ? (
                <IonIcon name="checkmark-outline" className="text-primary text-sm" />
              ) : (
                <IonIcon name="copy-outline" className="text-sm" />
              )}
              <span>{copied ? t("copied") : t("copy")}</span>
            </Button>
          </div>
          <pre className="p-4 font-mono text-xs text-zinc-300 overflow-x-auto leading-relaxed">
            {curlSnippet}
          </pre>
        </Card>
      </DialogContent>

      <DialogFooter>
        <Button onClick={() => onOpenChange(false)} className="h-8 text-xs">
          {t("close")}
        </Button>
      </DialogFooter>
    </Dialog>
  );
}

