import React from "react";
import * as Flags from "country-flag-icons/react/3x2";

interface CountryFlagProps {
  countryCode?: string; // "BR", "US", "PT", "ES", "MX", "FR", "IT", "DE", "GB", "JP"
  className?: string;
  title?: string;
}

export function CountryFlag({
  countryCode = "BR",
  className = "w-5 h-3.5 rounded-[2px] shadow-sm inline-block shrink-0 align-middle",
  title,
}: CountryFlagProps) {
  const code = (countryCode || "BR").toUpperCase();
  // @ts-expect-error dynamic flag indexing from country-flag-icons
  const FlagComponent = Flags[code];

  if (!FlagComponent) {
    return <Flags.BR className={className} title={title || "Brasil"} />;
  }

  return <FlagComponent className={className} title={title || code} />;
}
