import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: "EdgeGo Voice - Painel Administrativo de Personas",
  description: "Gerenciador de Personas de Áudio e Síntese de Voz com Microsoft Edge TTS e Go",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR" className={`dark ${inter.variable}`}>
      <head>
        <script
          type="module"
          src="https://unpkg.com/ionicons@7.1.0/dist/ionicons/ionicons.esm.js"
        />
        <script
          noModule
          src="https://unpkg.com/ionicons@7.1.0/dist/ionicons/ionicons.js"
        />
      </head>
      <body className="font-sans antialiased selection:bg-primary/20 selection:text-primary">
        {children}
      </body>
    </html>
  );
}
