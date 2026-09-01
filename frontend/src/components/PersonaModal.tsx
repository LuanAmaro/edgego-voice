"use client";

import React, { useState, useEffect } from "react";
import { Persona, VoiceOption } from "@/types";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogContent,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldLabel,
  FieldTitle,
} from "@/components/ui/field";
import { IonIcon } from "@/components/ui/ion-icon";
import { VoiceSelect } from "@/components/ui/voice-select";
import { useLanguage } from "@/i18n/LanguageContext";

interface PersonaModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  editingPersona: Persona | null;
  voices: VoiceOption[];
  onSave: (personaData: Partial<Persona>) => Promise<void>;
}

export function PersonaModal({
  open,
  onOpenChange,
  editingPersona,
  voices,
  onSave,
}: PersonaModalProps) {
  const { t } = useLanguage();
  const [id, setId] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [voice, setVoice] = useState("pt-BR-ThalitaMultilingualNeural");
  const [format, setFormat] = useState("mp3");
  const [speed, setSpeed] = useState(1.0);
  const [sanitize, setSanitize] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isIdCustomized, setIsIdCustomized] = useState(false);

  // Função para sanitizar e transformar texto em slug seguro
  const slugify = (text: string) => {
    return text
      .toLowerCase()
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "") // remove acentuações
      .replace(/[^a-z0-9\s-_]/g, "") // remove caracteres especiais
      .replace(/[\s_]+/g, "-") // converte espaços para hífens
      .replace(/-+/g, "-"); // colapsa múltiplos hífens
  };

  useEffect(() => {
    if (editingPersona) {
      setId(editingPersona.id);
      setName(editingPersona.name);
      setDescription(editingPersona.description || "");
      setVoice(editingPersona.voice);
      setFormat(editingPersona.format || "mp3");
      setSpeed(editingPersona.speed || 1.0);
      setSanitize(!editingPersona.remove_filter);
      setIsIdCustomized(true);
    } else {
      setId("");
      setName("");
      setDescription("");
      setVoice(voices[0]?.id || "pt-BR-ThalitaMultilingualNeural");
      setFormat("mp3");
      setSpeed(1.0);
      setSanitize(true);
      setIsIdCustomized(false);
    }
  }, [editingPersona, open, voices]);

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setName(val);
    if (!editingPersona && !isIdCustomized) {
      setId(slugify(val));
    }
  };

  const handleIdChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setIsIdCustomized(true);
    setId(slugify(e.target.value));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const finalId = slugify(id || name);
    if (!name.trim() || (!editingPersona && !finalId)) return;

    setIsSubmitting(true);
    try {
      await onSave({
        id: finalId,
        name: name.trim(),
        description: description.trim(),
        voice,
        format,
        speed,
        pitch: "+0Hz",
        remove_filter: !sanitize,
      });
      onOpenChange(false);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <form onSubmit={handleSubmit}>
        <DialogHeader onClose={() => onOpenChange(false)}>
          <DialogTitle>
            {editingPersona ? t("modalEditTitle") : t("modalNewTitle")}
          </DialogTitle>
          <DialogDescription>
            {t("modalDescription")}
          </DialogDescription>
        </DialogHeader>

        <DialogContent className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              {t("displayNameLabel")}
            </label>
            <Input
              value={name}
              onChange={handleNameChange}
              placeholder={t("displayNamePlaceholder")}
              required
              className="h-9 text-xs"
            />
          </div>

          <div className="space-y-1.5">
            <div className="flex justify-between items-center">
              <label className="text-xs font-medium text-muted-foreground">
                {t("uniqueIdLabel")}
              </label>
              <span className="text-[10px] text-muted-foreground font-mono">
                {t("uniqueIdHint")}
              </span>
            </div>
            <Input
              value={id}
              onChange={handleIdChange}
              placeholder={t("uniqueIdPlaceholder")}
              disabled={!!editingPersona}
              className="font-mono text-xs h-9"
              required
            />
            <span className="text-[10px] text-muted-foreground font-mono">
              {t("endpointPreview")}: /v1/persona/{id || "slug"}/speech
            </span>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              {t("descriptionLabel")}
            </label>
            <Textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t("descriptionPlaceholder")}
              rows={2}
              className="text-xs resize-none"
            />
          </div>

          {/* Voice Select */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              {t("neuralVoiceLabel")}
            </label>
            <VoiceSelect
              voices={voices}
              value={voice}
              onChange={(newVoice) => setVoice(newVoice)}
            />
          </div>

          {/* Format & Speed */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">
                {t("outputFormatLabel")}
              </label>
              <select
                value={format}
                onChange={(e) => setFormat(e.target.value)}
                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1.5 text-xs text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring transition-colors"
              >
                <option value="mp3">MP3 (Universal)</option>
                <option value="opus">Opus (WhatsApp / Telegram)</option>
                <option value="wav">WAV (Sem perdas)</option>
                <option value="pcm">PCM</option>
                <option value="flac">FLAC</option>
                <option value="aac">AAC</option>
              </select>
            </div>

            <div className="space-y-1.5">
              <div className="flex justify-between items-center">
                <label className="text-xs font-medium text-muted-foreground">
                  {t("speedLabel")}
                </label>
                <span className="text-xs font-mono text-muted-foreground">
                  {speed.toFixed(2)}x
                </span>
              </div>
              <input
                type="range"
                min="0.5"
                max="2.0"
                step="0.05"
                value={speed}
                onChange={(e) => setSpeed(parseFloat(e.target.value))}
                className="w-full mt-2.5 accent-primary h-1.5 bg-secondary rounded-lg cursor-pointer"
              />
            </div>
          </div>

          {/* Text Sanitizer Filter using Shadcn Field & Switch */}
          <FieldLabel htmlFor="switch-sanitize">
            <Field orientation="horizontal">
              <FieldContent>
                <FieldTitle>{t("smartSanitizerTitle")}</FieldTitle>
                <FieldDescription>
                  {t("smartSanitizerDesc")}
                </FieldDescription>
              </FieldContent>
              <Switch
                id="switch-sanitize"
                checked={sanitize}
                onCheckedChange={setSanitize}
              />
            </Field>
          </FieldLabel>
        </DialogContent>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="h-8 text-xs"
          >
            {t("cancel")}
          </Button>
          <Button
            type="submit"
            disabled={isSubmitting || !name.trim()}
            className="h-8 text-xs font-medium"
          >
            {isSubmitting ? (
              <>
                <IonIcon name="sync-outline" className="text-xs animate-spin mr-1.5" />
                <span>{t("saving")}</span>
              </>
            ) : (
              <>
                <IonIcon name="checkmark-outline" className="text-xs mr-1.5" />
                <span>{editingPersona ? t("saveChangesButton") : t("createPersonaSubmit")}</span>
              </>
            )}
          </Button>
        </DialogFooter>
      </form>
    </Dialog>
  );
}

