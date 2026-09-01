"use client";

import React, { useState } from "react";
import { Persona } from "@/types";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { IonIcon } from "@/components/ui/ion-icon";
import { useLanguage } from "@/i18n/LanguageContext";

interface ApiDocsViewProps {
  serverUrl: string;
  apiKey: string;
  personas: Persona[];
}

export function ApiDocsView({ serverUrl, apiKey, personas }: ApiDocsViewProps) {
  const { t } = useLanguage();
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  const samplePersona = personas[0] || { id: "whatsapp-suporte", format: "opus" };

  const copyToClipboard = (text: string, key: string) => {
    navigator.clipboard.writeText(text);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2000);
  };

  const curlSnippet = `curl -X POST "${serverUrl}/v1/persona/${samplePersona.id}/speech" \\
  -H "Authorization: Bearer ${apiKey}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "input": "Olá! Esta é uma mensagem de voz gerada via API."
  }' \\
  --output audio_${samplePersona.id}.${samplePersona.format}`;

  const pythonSnippet = `from openai import OpenAI

client = OpenAI(
    base_url="${serverUrl}/v1",
    api_key="${apiKey}"
)

# Chamada usando a API OpenAI compatível
response = client.audio.speech.create(
    model="tts-1",
    voice="alloy", # Mapeado em voices.json ou use a persona
    input="Olá! Teste de síntese com o SDK oficial da OpenAI e EdgeGo Voice."
)

response.stream_to_file("output.mp3")`;

  const nodeSnippet = `const fs = require('fs');

async function gerarAudio() {
  const response = await fetch('${serverUrl}/v1/persona/${samplePersona.id}/speech', {
    method: 'POST',
    headers: {
      'Authorization': 'Bearer ${apiKey}',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      input: 'Mensagem de voz rápida para o cliente.'
    })
  });

  const buffer = await response.arrayBuffer();
  fs.writeFileSync('output.${samplePersona.format}', Buffer.from(buffer));
  console.log('Áudio salvo com sucesso!');
}

gerarAudio();`;

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold text-foreground">
          {t("apiDocsTitle")}
        </h2>
        <p className="text-sm text-muted-foreground">
          {t("apiDocsSubtitle")}
        </p>
      </div>

      <div className="space-y-4">
        {/* cURL */}
        <Card className="overflow-hidden border-border bg-[#0C0A09]">
          <div className="flex items-center justify-between px-4 py-2.5 bg-white/[0.03] border-b border-border">
            <div className="flex items-center gap-2 text-xs font-semibold text-foreground">
              <IonIcon name="terminal-outline" className="text-primary text-base" />
              <span>{t("curlSectionTitle")}</span>
            </div>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => copyToClipboard(curlSnippet, "curl")}
              className="h-7 text-xs"
            >
              {copiedKey === "curl" ? (
                <IonIcon name="checkmark-outline" className="text-primary text-sm" />
              ) : (
                <IonIcon name="copy-outline" className="text-sm" />
              )}
              <span>{copiedKey === "curl" ? t("copied") : t("copy")}</span>
            </Button>
          </div>
          <pre className="p-4 font-mono text-xs text-zinc-300 overflow-x-auto leading-relaxed">
            {curlSnippet}
          </pre>
        </Card>

        {/* Python */}
        <Card className="overflow-hidden border-border bg-[#0C0A09]">
          <div className="flex items-center justify-between px-4 py-2.5 bg-white/[0.03] border-b border-border">
            <div className="flex items-center gap-2 text-xs font-semibold text-foreground">
              <IonIcon name="logo-python" className="text-primary text-base" />
              <span>{t("pythonSectionTitle")}</span>
            </div>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => copyToClipboard(pythonSnippet, "python")}
              className="h-7 text-xs"
            >
              {copiedKey === "python" ? (
                <IonIcon name="checkmark-outline" className="text-primary text-sm" />
              ) : (
                <IonIcon name="copy-outline" className="text-sm" />
              )}
              <span>{copiedKey === "python" ? t("copied") : t("copy")}</span>
            </Button>
          </div>
          <pre className="p-4 font-mono text-xs text-zinc-300 overflow-x-auto leading-relaxed">
            {pythonSnippet}
          </pre>
        </Card>

        {/* Node.js */}
        <Card className="overflow-hidden border-border bg-[#0C0A09]">
          <div className="flex items-center justify-between px-4 py-2.5 bg-white/[0.03] border-b border-border">
            <div className="flex items-center gap-2 text-xs font-semibold text-foreground">
              <IonIcon name="logo-nodejs" className="text-primary text-base" />
              <span>{t("nodeSectionTitle")}</span>
            </div>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => copyToClipboard(nodeSnippet, "node")}
              className="h-7 text-xs"
            >
              {copiedKey === "node" ? (
                <IonIcon name="checkmark-outline" className="text-primary text-sm" />
              ) : (
                <IonIcon name="copy-outline" className="text-sm" />
              )}
              <span>{copiedKey === "node" ? t("copied") : t("copy")}</span>
            </Button>
          </div>
          <pre className="p-4 font-mono text-xs text-zinc-300 overflow-x-auto leading-relaxed">
            {nodeSnippet}
          </pre>
        </Card>
      </div>
    </div>
  );
}

