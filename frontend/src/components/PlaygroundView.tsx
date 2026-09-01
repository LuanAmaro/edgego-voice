"use client";

import React, { useState, useRef, useEffect } from "react";
import { Persona, VoiceOption } from "@/types";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { IonIcon } from "@/components/ui/ion-icon";
import { VoiceSelect } from "@/components/ui/voice-select";
import { WaveformPlayer } from "@/components/WaveformPlayer";
import { Textarea } from "@/components/ui/textarea";
import { formatBytes } from "@/lib/utils";
import { useLanguage } from "@/i18n/LanguageContext";

interface PlaygroundViewProps {
  personas: Persona[];
  voices: VoiceOption[];
  serverUrl: string;
  apiKey: string;
  selectedPersonaId?: string;
}

interface ServerTelemetry {
  goroutines: number;
  alloc_mb: number;
  sys_mb: number;
  heap_inuse_mb: number;
  cache_items: number;
  cache_bytes: number;
}

export function PlaygroundView({
  personas,
  voices,
  serverUrl,
  apiKey,
  selectedPersonaId: initialPersonaId,
}: PlaygroundViewProps) {
  const { t } = useLanguage();
  const initialPersona = personas.find((x) => x.id === initialPersonaId);

  const [selectedPersonaId, setSelectedPersonaId] = useState(
    initialPersonaId || ""
  );
  const [voice, setVoice] = useState(
    initialPersona?.voice || voices[0]?.id || "pt-BR-ThalitaMultilingualNeural"
  );
  const [format, setFormat] = useState(initialPersona?.format || "mp3");
  const [speed, setSpeed] = useState(initialPersona?.speed || 1.0);
  const [text, setText] = useState(
    "Olá! Este é um teste no EdgeGo Voice para validar a velocidade e qualidade da voz sintetizada."
  );


  // Sincronizar parâmetros quando a Persona for alterada externamente (ex: botão "Testar")
  useEffect(() => {
    if (initialPersonaId) {
      setSelectedPersonaId(initialPersonaId);
      const p = personas.find((x) => x.id === initialPersonaId);
      if (p) {
        setVoice(p.voice);
        setFormat(p.format || "mp3");
        setSpeed(p.speed || 1.0);
      }
    }
  }, [initialPersonaId, personas]);

  const [isLoading, setIsLoading] = useState(false);
  const [audioUrl, setAudioUrl] = useState<string | null>(null);
  const [telemetry, setTelemetry] = useState<ServerTelemetry | null>(null);

  const [alertInfo, setAlertInfo] = useState<{
    variant: "default" | "destructive";
    title: string;
    message: string;
  } | null>(null);

  const [metrics, setMetrics] = useState<{
    ttfb: number;
    totalTime: number;
    size: number;
    cacheHit?: boolean;
  } | null>(null);

  const abortControllerRef = useRef<AbortController | null>(null);
  const alertTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Buscar telemetria real do runtime Go
  const fetchTelemetry = async () => {
    try {
      const resp = await fetch(`${serverUrl}/health`);
      if (resp.ok) {
        const data = await resp.json();
        setTelemetry(data);
      }
    } catch (e) {}
  };

  useEffect(() => {
    fetchTelemetry();
    const interval = setInterval(fetchTelemetry, 6000);
    return () => clearInterval(interval);
  }, [serverUrl]);

  // Auto-dismiss alert after 4.5s
  useEffect(() => {
    if (alertInfo) {
      if (alertTimeoutRef.current) clearTimeout(alertTimeoutRef.current);
      alertTimeoutRef.current = setTimeout(() => {
        setAlertInfo(null);
      }, 4500);
    }
    return () => {
      if (alertTimeoutRef.current) clearTimeout(alertTimeoutRef.current);
    };
  }, [alertInfo]);

  const handlePersonaChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const pId = e.target.value;
    setSelectedPersonaId(pId);
    if (!pId) return;

    const p = personas.find((x) => x.id === pId);
    if (p) {
      setVoice(p.voice);
      setFormat(p.format);
      setSpeed(p.speed);
    }
  };

  const handleSynthesize = async () => {
    if (!text.trim()) return;

    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    const controller = new AbortController();
    abortControllerRef.current = controller;

    if (audioUrl) {
      URL.revokeObjectURL(audioUrl);
      setAudioUrl(null);
    }

    setIsLoading(true);
    setAlertInfo(null);
    setMetrics(null);

    const startTime = performance.now();
    let firstByteTime = 0;

    try {
      let endpoint = `${serverUrl}/v1/audio/speech`;
      let bodyData: any = {
        model: "tts-1",
        input: text.trim(),
        voice,
        response_format: format,
        speed,
      };

      if (selectedPersonaId) {
        endpoint = `${serverUrl}/v1/persona/${selectedPersonaId}/speech`;
        bodyData = { input: text.trim() };
      }

      const resp = await fetch(endpoint, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${apiKey}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(bodyData),
        signal: controller.signal,
      });

      if (!resp.ok) {
        const errJson = await resp.json().catch(() => ({}));
        throw new Error(
          errJson.error?.message || `Erro ${resp.status} ao sintetizar áudio.`
        );
      }

      const isCacheHit = resp.headers.get("X-Cache") === "HIT";

      if (isCacheHit) {
        const blob = await resp.blob();
        const latency = Math.round(performance.now() - startTime);

        const newUrl = URL.createObjectURL(blob);
        setAudioUrl(newUrl);
        setMetrics({
          ttfb: latency,
          totalTime: latency,
          size: blob.size,
          cacheHit: true,
        });

        setAlertInfo({
          variant: "default",
          title: "Áudio Pronto (Cache HIT)",
          message: `Entregue em ${latency}ms via memória RAM (${formatBytes(blob.size)})`,
        });
        fetchTelemetry();
        return;
      }

      if (!resp.body) {
        throw new Error("Streaming não suportado pelo navegador.");
      }

      const reader = resp.body.getReader();
      const chunks: Uint8Array[] = [];
      let totalBytesReceived = 0;

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        if (value && value.length > 0) {
          if (firstByteTime === 0) {
            firstByteTime = Math.round(performance.now() - startTime);
          }
          chunks.push(value);
          totalBytesReceived += value.length;
        }
      }

      const totalTime = Math.round(performance.now() - startTime);

      if (totalBytesReceived === 0) {
        throw new Error("O áudio retornado está vazio. Tente outra voz neural.");
      }

      const mimeType =
        format === "opus"
          ? "audio/ogg; codecs=opus"
          : format === "mp3"
          ? "audio/mpeg"
          : format === "wav"
          ? "audio/wav"
          : `audio/${format}`;
      const finalBlob = new Blob(chunks as BlobPart[], { type: mimeType });
      const newUrl = URL.createObjectURL(finalBlob);

      setAudioUrl(newUrl);
      setMetrics({
        ttfb: firstByteTime || totalTime,
        totalTime,
        size: totalBytesReceived,
        cacheHit: false,
      });

      setAlertInfo({
        variant: "default",
        title: "Áudio sintetizado",
        message: `1º bloco em ${firstByteTime || totalTime}ms • Total: ${totalTime}ms (${formatBytes(totalBytesReceived)})`,
      });

      fetchTelemetry();
    } catch (err: any) {
      if (err.name === "AbortError") return;
      setAlertInfo({
        variant: "destructive",
        title: "Falha na síntese",
        message: err.message,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6 max-w-5xl mx-auto relative">
      {/* Header */}
      <div className="flex flex-col gap-1">
        <h2 className="text-lg font-semibold tracking-tight text-foreground">
          {t("playgroundTitle")}
        </h2>
        <p className="text-xs text-muted-foreground">
          {t("playgroundSubtitle")}
        </p>
      </div>

      {/* Real-time Telemetry Grid from Go Engine */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <Card className="p-3.5 bg-card/60 border-border/80 shadow-none flex flex-col justify-between">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-medium text-muted-foreground">
              {t("latencyTtfb")}
            </span>
            <IonIcon name="flash-outline" className="text-xs text-primary" />
          </div>
          <div className="text-xl font-bold text-foreground font-mono mt-1">
            {metrics ? `${metrics.ttfb} ms` : "--"}
          </div>
          <span className="text-[10px] text-muted-foreground font-mono">
            {metrics?.cacheHit ? t("cacheHitYes") : t("cacheHitNo")}
          </span>
        </Card>

        <Card className="p-3.5 bg-card/60 border-border/80 shadow-none flex flex-col justify-between">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-medium text-muted-foreground">
              {t("totalDuration")}
            </span>
            <IonIcon name="timer-outline" className="text-xs text-muted-foreground" />
          </div>
          <div className="text-xl font-bold text-foreground font-mono mt-1">
            {metrics ? `${(metrics.totalTime / 1000).toFixed(2)} s` : "--"}
          </div>
          <span className="text-[10px] text-muted-foreground font-mono">
            {metrics ? formatBytes(metrics.size) : t("audioSize")}
          </span>
        </Card>

        <Card className="p-3.5 bg-card/60 border-border/80 shadow-none flex flex-col justify-between">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-medium text-muted-foreground">
              {t("allocatedRam")}
            </span>
            <IonIcon name="hardware-chip-outline" className="text-xs text-muted-foreground" />
          </div>
          <div className="text-xl font-bold text-foreground font-mono mt-1">
            {telemetry ? `${telemetry.alloc_mb.toFixed(2)} MB` : "--"}
          </div>
          <span className="text-[10px] text-muted-foreground font-mono">
            Heap: {telemetry ? `${telemetry.heap_inuse_mb.toFixed(1)} MB` : "--"} • Sys: {telemetry ? `${telemetry.sys_mb.toFixed(1)} MB` : "--"}
          </span>
        </Card>

        <Card className="p-3.5 bg-card/60 border-border/80 shadow-none flex flex-col justify-between">
          <div className="flex items-center justify-between">
            <span className="text-[11px] font-medium text-muted-foreground">
              {t("goroutinesCount")}
            </span>
            <IonIcon name="pulse-outline" className="text-xs text-emerald-400" />
          </div>
          <div className="text-xl font-bold text-foreground font-mono mt-1 flex items-center gap-1.5">
            <span>{telemetry ? telemetry.goroutines : "--"}</span>
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
          </div>
          <span className="text-[10px] text-muted-foreground font-mono">
            Cache: {telemetry ? `${telemetry.cache_items} itens (${formatBytes(telemetry.cache_bytes)})` : "--"}
          </span>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Left Column: Form Controls (6 cols) */}
        <Card className="lg:col-span-6 p-5 space-y-4 bg-card/60 border-border/80 shadow-none">
          {/* Persona selector */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              {t("loadPersonaPreset")}
            </label>
            <select
              value={selectedPersonaId}
              onChange={handlePersonaChange}
              className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1.5 text-xs text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring transition-colors"
            >
              <option value="">{t("noneCustom")}</option>
              {personas.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name} ({p.id})
                </option>
              ))}
            </select>
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

          {/* Text Input Shadcn Component */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-muted-foreground">
              {t("textInputLabel")}
            </label>
            <Textarea
              value={text}
              onChange={(e) => setText(e.target.value)}
              rows={4}
              placeholder={t("textInputPlaceholder")}
              className="resize-none"
            />
          </div>

          <Button
            onClick={handleSynthesize}
            disabled={isLoading || !text.trim()}
            className="w-full h-9 text-xs font-medium"
          >
            {isLoading ? (
              <>
                <IonIcon name="sync-outline" className="text-sm animate-spin mr-2" />
                <span>{t("synthesizingButton")}</span>
              </>
            ) : (
              <>
                <IonIcon name="volume-high-outline" className="text-sm mr-2" />
                <span>{t("synthesizeButton")}</span>
              </>
            )}
          </Button>
        </Card>

        {/* Right Column: Waveform Player Card (6 cols) */}
        <Card className="lg:col-span-6 p-5 flex flex-col items-center justify-center text-center space-y-4 bg-card/60 border-border/80 shadow-none min-h-[300px]">
          {audioUrl && !isLoading ? (
            <div className="w-full">
              {/* Interactive Waveform Player */}
              <WaveformPlayer
                audioUrl={audioUrl}
                format={format}
                personaId={selectedPersonaId}
                metrics={metrics}
              />
            </div>
          ) : (
            <div className="space-y-3 py-6">
              <div className="w-12 h-12 mx-auto rounded-full bg-secondary/80 flex items-center justify-center text-muted-foreground text-xl">
                <IonIcon
                  name={isLoading ? "sync-outline" : "musical-notes-outline"}
                  className={isLoading ? "animate-spin text-foreground" : ""}
                />
              </div>

              <div className="space-y-1">
                <h3 className="text-xs font-medium text-foreground">
                  {isLoading ? "Sintetizando áudio..." : "Waveform de Áudio"}
                </h3>
                <p className="text-[11px] text-muted-foreground max-w-xs mx-auto">
                  {isLoading
                    ? "Recebendo stream de pacotes em tempo real..."
                    : "Digite um texto e clique em 'Gerar Áudio' para visualizar e escutar as ondas sonoras."}
                </p>
              </div>
            </div>
          )}
        </Card>
      </div>

      {/* Subtle Bottom-Right Notification Toast with Solid Opaque Dark Background */}
      {alertInfo && (
        <div className="fixed bottom-6 right-6 z-50 w-full max-w-sm animate-toast pointer-events-auto">
          <div
            role="alert"
            style={{ backgroundColor: "#09090b" }}
            className="border border-zinc-700/90 text-foreground shadow-[0_20px_50px_rgba(0,0,0,0.95),0_0_0_1px_rgba(255,255,255,0.08)] p-3.5 flex items-start gap-3 rounded-xl transition-all"
          >
            <IonIcon
              name={
                alertInfo.variant === "destructive"
                  ? "alert-circle-outline"
                  : "checkmark-circle-outline"
              }
              className={`text-base mt-0.5 shrink-0 ${
                alertInfo.variant === "destructive" ? "text-red-400" : "text-emerald-400"
              }`}
            />
            <div className="flex-1 min-w-0 pr-1">
              <h5 className="text-xs font-semibold text-foreground tracking-tight">
                {alertInfo.title}
              </h5>
              <p className="text-[11px] text-zinc-300 mt-0.5 leading-relaxed">
                {alertInfo.message}
              </p>
            </div>
            <button
              onClick={() => setAlertInfo(null)}
              className="text-zinc-400 hover:text-foreground text-sm shrink-0 -mt-0.5 p-1 rounded-md hover:bg-white/10 transition-colors"
            >
              <IonIcon name="close-outline" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
