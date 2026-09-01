"use client";

import React, { useState, useEffect } from "react";
import { Persona, VoiceOption, TabType } from "@/types";
import { LoginView } from "@/components/LoginView";
import { Navbar } from "@/components/Navbar";
import { PersonasList } from "@/components/PersonasList";
import { PersonaModal } from "@/components/PersonaModal";
import { DeletePersonaModal } from "@/components/DeletePersonaModal";
import { SnippetModal } from "@/components/SnippetModal";
import { PlaygroundView } from "@/components/PlaygroundView";
import { ApiDocsView } from "@/components/ApiDocsView";
import { IonIcon } from "@/components/ui/ion-icon";
import { LanguageProvider, useLanguage } from "@/i18n/LanguageContext";

function MainDashboard() {
  const { t } = useLanguage();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [serverUrl, setServerUrl] = useState("");
  const [apiKey, setApiKey] = useState("");
  const [isLoadingAuth, setIsLoadingAuth] = useState(false);

  const [activeTab, setActiveTab] = useState<TabType>("personas");
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [voices, setVoices] = useState<VoiceOption[]>([]);

  // Modals state
  const [isPersonaModalOpen, setIsPersonaModalOpen] = useState(false);
  const [editingPersona, setEditingPersona] = useState<Persona | null>(null);

  const [isSnippetModalOpen, setIsSnippetModalOpen] = useState(false);
  const [snippetPersona, setSnippetPersona] = useState<Persona | null>(null);

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [deletingPersona, setDeletingPersona] = useState<Persona | null>(null);

  const [playgroundPersonaId, setPlaygroundPersonaId] = useState<string>("");

  useEffect(() => {
    const savedUrl =
      localStorage.getItem("edgego_server_url") || window.location.origin;
    const savedKey = localStorage.getItem("edgego_api_key") || "";

    setServerUrl(savedUrl);
    setApiKey(savedKey);

    if (savedKey) {
      handleLogin(savedKey, savedUrl, true);
    }
  }, []);

  const handleLogin = async (key: string, url: string, isSilent = false) => {
    setIsLoadingAuth(true);
    try {
      const resp = await fetch(`${url}/api/personas`, {
        headers: { Authorization: `Bearer ${key}` },
      });

      if (!resp.ok) {
        throw new Error(t("invalidAuthError"));
      }

      setServerUrl(url);
      setApiKey(key);
      localStorage.setItem("edgego_server_url", url);
      localStorage.setItem("edgego_api_key", key);
      setIsAuthenticated(true);

      // Carregar vozes e personas
      await Promise.all([fetchVoices(url), fetchPersonas(url, key)]);
    } catch (err: any) {
      if (!isSilent) {
        alert(err.message);
      } else {
        localStorage.removeItem("edgego_api_key");
      }
    } finally {
      setIsLoadingAuth(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("edgego_api_key");
    setIsAuthenticated(false);
    setApiKey("");
  };

  const fetchVoices = async (url: string) => {
    try {
      const resp = await fetch(`${url}/api/voices`);
      const data = await resp.json();
      setVoices(data.data || []);
    } catch (err) {
      console.error("Erro ao carregar vozes:", err);
    }
  };

  const fetchPersonas = async (url: string, key: string) => {
    try {
      const resp = await fetch(`${url}/api/personas`, {
        headers: { Authorization: `Bearer ${key}` },
      });
      const data = await resp.json();
      setPersonas(data.data || []);
    } catch (err) {
      console.error("Erro ao carregar personas:", err);
    }
  };

  const handleSavePersona = async (personaData: Partial<Persona>) => {
    try {
      let resp;
      if (editingPersona) {
        resp = await fetch(`${serverUrl}/api/personas/${editingPersona.id}`, {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${apiKey}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify(personaData),
        });
      } else {
        resp = await fetch(`${serverUrl}/api/personas`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${apiKey}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify(personaData),
        });
      }

      if (!resp.ok) {
        const err = await resp.json().catch(() => ({}));
        throw new Error(err.error?.message || "Falha ao salvar persona.");
      }

      await fetchPersonas(serverUrl, apiKey);
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleDeletePersona = (persona: Persona) => {
    setDeletingPersona(persona);
    setIsDeleteModalOpen(true);
  };

  const handleConfirmDelete = async (persona: Persona) => {
    const resp = await fetch(`${serverUrl}/api/personas/${persona.id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${apiKey}` },
    });

    if (!resp.ok) {
      throw new Error("Erro ao excluir persona.");
    }

    await fetchPersonas(serverUrl, apiKey);
  };

  if (!isAuthenticated) {
    return (
      <LoginView
        serverUrl={serverUrl}
        setServerUrl={setServerUrl}
        onLogin={(key, url) => handleLogin(key, url, false)}
        isLoading={isLoadingAuth}
      />
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-background text-foreground">
      <Navbar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        personasCount={personas.length}
        onLogout={handleLogout}
      />

      <main className="flex-1 max-w-7xl w-full mx-auto px-6 py-8">
        {activeTab === "personas" && (
          <PersonasList
            personas={personas}
            voices={voices}
            onOpenNewModal={() => {
              setEditingPersona(null);
              setIsPersonaModalOpen(true);
            }}
            onTest={(p) => {
              setPlaygroundPersonaId(p.id);
              setActiveTab("playground");
            }}
            onSnippet={(p) => {
              setSnippetPersona(p);
              setIsSnippetModalOpen(true);
            }}
            onEdit={(p) => {
              setEditingPersona(p);
              setIsPersonaModalOpen(true);
            }}
            onDelete={handleDeletePersona}
          />
        )}

        {activeTab === "playground" && (
          <PlaygroundView
            personas={personas}
            voices={voices}
            serverUrl={serverUrl}
            apiKey={apiKey}
            selectedPersonaId={playgroundPersonaId}
          />
        )}

        {activeTab === "apidocs" && (
          <ApiDocsView
            serverUrl={serverUrl}
            apiKey={apiKey}
            personas={personas}
          />
        )}
      </main>

      {/* Modals */}
      <PersonaModal
        open={isPersonaModalOpen}
        onOpenChange={setIsPersonaModalOpen}
        editingPersona={editingPersona}
        voices={voices}
        onSave={handleSavePersona}
      />

      <SnippetModal
        open={isSnippetModalOpen}
        onOpenChange={setIsSnippetModalOpen}
        persona={snippetPersona}
        serverUrl={serverUrl}
        apiKey={apiKey}
      />

      <DeletePersonaModal
        open={isDeleteModalOpen}
        onOpenChange={setIsDeleteModalOpen}
        persona={deletingPersona}
        onConfirm={handleConfirmDelete}
      />

      <footer className="border-t border-border px-6 py-4 flex flex-col sm:flex-row items-center justify-between gap-2 text-xs text-muted-foreground">
        <div className="flex items-center gap-1.5">
          <span>{t("developedBy")}</span>
          <a
            href="https://wa.me/5584981277917"
            target="_blank"
            rel="noopener noreferrer"
            className="font-medium text-foreground hover:text-primary transition-colors inline-flex items-center gap-1 hover:underline underline-offset-4"
          >
            <span>Rodolfo Bandeira</span>
            <IonIcon name="logo-whatsapp" className="text-sm text-emerald-400" />
          </a>
        </div>
        <div className="font-mono text-[11px]">{t("version")}</div>
      </footer>
    </div>
  );
}

export default function Home() {
  return (
    <LanguageProvider>
      <MainDashboard />
    </LanguageProvider>
  );
}

