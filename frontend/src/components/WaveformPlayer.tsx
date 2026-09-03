"use client";

import React, { useState, useRef, useEffect, useMemo } from "react";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";

interface WaveformPlayerProps {
  audioUrl: string;
  format: string;
  personaId?: string;
  metrics: {
    ttfb: number;
    totalTime: number;
    size: number;
    cacheHit?: boolean;
  } | null;
  onPlayStateChange?: (isPlaying: boolean) => void;
}

export function WaveformPlayer({ audioUrl, format, personaId, onPlayStateChange }: WaveformPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(false);

  useEffect(() => {
    onPlayStateChange?.(isPlaying);
  }, [isPlaying, onPlayStateChange]);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [playbackRate, setPlaybackRate] = useState(1.0);
  const [isMuted, setIsMuted] = useState(false);
  const [hoverPct, setHoverPct] = useState<number | null>(null);
  const [needsReload, setNeedsReload] = useState(false);

  const audioRef = useRef<HTMLAudioElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  const peaks = useMemo(() => {
    const arr: number[] = [];
    for (let i = 0; i < 54; i++) {
      const s1 = Math.sin((i / 54) * Math.PI) * 0.65;
      const s2 = Math.sin(i * 0.75) * 0.22;
      const s3 = Math.cos(i * 1.4) * 0.13;
      arr.push(Math.max(0.18, Math.min(0.95, Math.abs(s1 + s2 + s3))));
    }
    return arr;
  }, []);

  // Toca automaticamente quando audioUrl muda
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !audioUrl) return;

    setCurrentTime(0);
    setDuration(0);
    setIsPlaying(false);
    setNeedsReload(false);

    // load() + play() na nova URL
    audio.load();
    audio.play()
      .then(() => setIsPlaying(true))
      .catch(() => setIsPlaying(false));
  }, [audioUrl]);

  const handleLoadedMetadata = () => {
    const audio = audioRef.current;
    if (audio && isFinite(audio.duration) && audio.duration > 0) {
      setDuration(audio.duration);
    }
  };

  const handleTimeUpdate = () => {
    const audio = audioRef.current;
    if (!audio) return;
    setCurrentTime(audio.currentTime || 0);
    if (!duration && audio.duration && isFinite(audio.duration) && audio.duration > 0) {
      setDuration(audio.duration);
    }
  };

  const handleEnded = () => {
    setIsPlaying(false);
    setCurrentTime(0);
    // Marcar que precisa de reload antes do próximo play
    setNeedsReload(true);
  };

  const doPlay = async (offset?: number) => {
    const audio = audioRef.current;
    if (!audio) return;

    if (needsReload) {
      // Ao terminar o audio, o browser expira a stream do Blob.
      // load() reinicializa o decoder interno sem re-baixar o arquivo (o Blob ainda existe).
      audio.load();
      setNeedsReload(false);
      // Aguardar o browser estar pronto (canplay)
      await new Promise<void>((resolve) => {
        const onReady = () => {
          audio.removeEventListener("canplay", onReady);
          resolve();
        };
        audio.addEventListener("canplay", onReady);
      });
    }

    if (offset !== undefined) {
      audio.currentTime = offset;
    } else if (audio.ended || audio.currentTime >= (audio.duration || duration)) {
      audio.currentTime = 0;
    }

    try {
      await audio.play();
      setIsPlaying(true);
    } catch {
      // Último fallback
      audio.currentTime = 0;
      audio.play().then(() => setIsPlaying(true)).catch(() => {});
    }
  };

  const togglePlayPause = () => {
    const audio = audioRef.current;
    if (!audio) return;

    if (isPlaying) {
      audio.pause();
      setIsPlaying(false);
    } else {
      doPlay();
    }
  };

  const restartAudio = async () => {
    const audio = audioRef.current;
    if (!audio) return;

    audio.pause();
    setCurrentTime(0);
    setNeedsReload(false);

    // Sempre reinicializa o decoder para garantir que o restart funcione
    // tanto durante a reprodução quanto após o término
    audio.load();

    await new Promise<void>((resolve) => {
      const onReady = () => {
        audio.removeEventListener("canplay", onReady);
        resolve();
      };
      audio.addEventListener("canplay", onReady);
    });

    audio.play().then(() => setIsPlaying(true)).catch(() => {});
  };

  const seekTo = (pct: number) => {
    const audio = audioRef.current;
    if (!audio) return;
    const totalDur = duration > 0 ? duration : (audio.duration && isFinite(audio.duration) ? audio.duration : 1);
    const targetTime = Math.max(0, Math.min(1, pct)) * totalDur;
    setCurrentTime(targetTime);
    doPlay(targetTime);
  };

  const handleContainerClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!containerRef.current) return;
    const rect = containerRef.current.getBoundingClientRect();
    seekTo((e.clientX - rect.left) / rect.width);
  };

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!containerRef.current) return;
    const rect = containerRef.current.getBoundingClientRect();
    setHoverPct(Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width)));
  };

  const toggleMute = () => {
    const audio = audioRef.current;
    if (!audio) return;
    audio.muted = !isMuted;
    setIsMuted(!isMuted);
  };

  const changeSpeed = () => {
    const audio = audioRef.current;
    if (!audio) return;
    const speeds = [1.0, 1.25, 1.5, 2.0];
    const next = speeds[(speeds.indexOf(playbackRate) + 1) % speeds.length];
    audio.playbackRate = next;
    setPlaybackRate(next);
  };

  const formatTime = (secs: number) => {
    if (isNaN(secs) || !isFinite(secs) || secs < 0) return "0:00";
    const m = Math.floor(secs / 60);
    const s = Math.floor(secs % 60);
    return `${m}:${s < 10 ? "0" : ""}${s}`;
  };

  const currentDuration = duration > 0 ? duration : (audioRef.current?.duration && isFinite(audioRef.current.duration) ? audioRef.current.duration : 0);
  const progressPct = currentDuration > 0 ? Math.min(100, (currentTime / currentDuration) * 100) : 0;

  return (
    <div className="w-full space-y-3 pt-1 animate-card-enter">
      <audio
        ref={audioRef}
        src={audioUrl}
        preload="auto"
        onTimeUpdate={handleTimeUpdate}
        onLoadedMetadata={handleLoadedMetadata}
        onEnded={handleEnded}
        onPlay={() => setIsPlaying(true)}
        onPause={() => setIsPlaying(false)}
        className="hidden"
      />

      {/* Waveform Visualizer */}
      <div
        ref={containerRef}
        onClick={handleContainerClick}
        onMouseMove={handleMouseMove}
        onMouseLeave={() => setHoverPct(null)}
        className="relative w-full h-24 bg-[#09090b]/80 border border-border/70 rounded-xl p-3.5 flex items-center justify-between gap-[3px] cursor-pointer select-none group transition-colors hover:border-zinc-700 overflow-hidden"
      >
        {hoverPct !== null && currentDuration > 0 && (
          <div
            className="absolute top-1 bottom-1 w-[1.5px] bg-foreground/60 pointer-events-none z-10"
            style={{ left: `${hoverPct * 100}%` }}
          >
            <span className="absolute -top-1 left-1/2 -translate-x-1/2 -translate-y-full text-[9px] font-mono bg-zinc-900 border border-border px-1.5 py-0.5 rounded text-foreground shadow-md">
              {formatTime(hoverPct * currentDuration)}
            </span>
          </div>
        )}

        <div
          className="absolute top-0 bottom-0 w-[2px] bg-primary z-10 pointer-events-none transition-all duration-75 shadow-[0_0_8px_rgba(0,223,129,0.5)]"
          style={{ left: `${progressPct}%` }}
        />

        {peaks.map((amplitude, idx) => {
          const barPct = (idx / peaks.length) * 100;
          const isPassed = barPct <= progressPct;
          return (
            <div key={idx} className="flex-1 flex items-center justify-center h-full pointer-events-none">
              <div
                className={`w-full max-w-[3px] rounded-full transition-all duration-100 ${
                  isPassed
                    ? "bg-primary shadow-[0_0_4px_rgba(0,223,129,0.4)]"
                    : "bg-[#27272a] group-hover:bg-[#323238]"
                }`}
                style={{ height: `${Math.max(12, amplitude * 85)}%` }}
              />
            </div>
          );
        })}
      </div>

      {/* Controls */}
      <div className="flex items-center justify-between px-1">
        <div className="flex items-center gap-2">
          <Button variant="primary" size="icon" onClick={togglePlayPause} className="h-8 w-8 rounded-full shadow-sm" title={isPlaying ? "Pausar" : "Reproduzir"}>
            <IonIcon name={isPlaying ? "pause-outline" : "play-outline"} className="text-base" />
          </Button>

          <Button variant="ghost" size="icon" onClick={restartAudio} className="h-8 w-8 rounded-full text-muted-foreground hover:text-foreground" title="Reiniciar">
            <IonIcon name="reload-outline" className="text-sm" />
          </Button>

          <button onClick={changeSpeed} className="text-[11px] font-mono font-medium text-muted-foreground hover:text-foreground bg-secondary px-2 py-0.5 rounded transition-colors">
            {playbackRate}x
          </button>
        </div>

        <div className="text-xs font-mono text-muted-foreground select-none">
          <span className="text-foreground font-semibold">{formatTime(currentTime)}</span>
          <span className="mx-1">/</span>
          <span>{formatTime(currentDuration)}</span>
        </div>

        <div className="flex items-center gap-1.5">
          <Button variant="ghost" size="icon" onClick={toggleMute} className="h-8 w-8 rounded-full text-muted-foreground hover:text-foreground" title={isMuted ? "Desmutar" : "Mutar"}>
            <IonIcon name={isMuted ? "volume-mute-outline" : "volume-high-outline"} className="text-sm" />
          </Button>

          <a
            href={audioUrl}
            download={`audio_${personaId || "speech"}.${format}`}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md border border-input bg-background hover:bg-accent text-xs font-medium text-foreground transition-colors shadow-sm"
          >
            <IonIcon name="download-outline" className="text-sm" />
            <span>Baixar {format.toUpperCase()}</span>
          </a>
        </div>
      </div>
    </div>
  );
}
