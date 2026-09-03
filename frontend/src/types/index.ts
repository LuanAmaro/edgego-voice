export interface Persona {
  id: string;
  name: string;
  description: string;
  voice: string;
  format: string;
  speed: number;
  pitch: string;
  break_comma?: string;
  break_period?: string;
  telephony?: boolean;
  auto_breath?: boolean;
  remove_filter: boolean;
  api_key?: string;
  created_at?: string;
  updated_at?: string;
}

export interface VoiceOption {
  id: string;
  name: string;
  language: string;
  gender: string;
  country_code?: string;
  flag?: string;
  locale?: string;
  description: string;
}

export type TabType = "personas" | "playground" | "apidocs";

