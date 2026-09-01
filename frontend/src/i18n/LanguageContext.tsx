"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import {
  LanguageCode,
  LanguageOption,
  SUPPORTED_LANGUAGES,
  translations,
  TranslationKeys,
} from "./translations";

interface LanguageContextType {
  language: LanguageCode;
  setLanguage: (lang: LanguageCode) => void;
  currentLanguage: LanguageOption;
  languages: LanguageOption[];
  t: (key: TranslationKeys) => string;
}

const LanguageContext = createContext<LanguageContextType | undefined>(
  undefined
);

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const [language, setLanguageState] = useState<LanguageCode>("pt-BR");

  useEffect(() => {
    const saved = localStorage.getItem("edgego_lang") as LanguageCode;
    if (saved && (saved === "pt-BR" || saved === "en-US" || saved === "es-ES")) {
      setLanguageState(saved);
    }
  }, []);

  const setLanguage = (lang: LanguageCode) => {
    setLanguageState(lang);
    localStorage.setItem("edgego_lang", lang);
  };

  const currentLanguage =
    SUPPORTED_LANGUAGES.find((l) => l.code === language) ||
    SUPPORTED_LANGUAGES[0];

  const t = (key: TranslationKeys): string => {
    const langDict = translations[language] || translations["pt-BR"];
    return langDict[key] || translations["pt-BR"][key] || (key as string);
  };

  return (
    <LanguageContext.Provider
      value={{
        language,
        setLanguage,
        currentLanguage,
        languages: SUPPORTED_LANGUAGES,
        t,
      }}
    >
      {children}
    </LanguageContext.Provider>
  );
}

export function useLanguage() {
  const context = useContext(LanguageContext);
  if (!context) {
    throw new Error("useLanguage must be used within a LanguageProvider");
  }
  return context;
}
