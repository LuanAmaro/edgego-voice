export type LanguageCode = "pt-BR" | "en-US" | "es-ES";

export interface LanguageOption {
  code: LanguageCode;
  name: string;
  nativeName: string;
  countryCode: string;
  flagTitle: string;
}

export const SUPPORTED_LANGUAGES: LanguageOption[] = [
  {
    code: "pt-BR",
    name: "Português",
    nativeName: "Português (Brasil)",
    countryCode: "BR",
    flagTitle: "Brasil",
  },
  {
    code: "en-US",
    name: "English",
    nativeName: "English (US)",
    countryCode: "US",
    flagTitle: "United States",
  },
  {
    code: "es-ES",
    name: "Español",
    nativeName: "Español",
    countryCode: "ES",
    flagTitle: "España",
  },
];

export const translations = {
  "pt-BR": {
    // Navbar & Common
    statusOnline: "Online",
    logout: "Desconectar",
    selectLanguage: "Selecionar Idioma",
    changeLanguageDescription: "Escolha o idioma preferido para a interface do painel.",
    personas: "Personas",
    playground: "Playground TTS",
    apiDocs: "API",
    version: "v1.0 (Go Native + Next.js)",
    developedBy: "Desenvolvido por",
    cancel: "Cancelar",
    save: "Salvar",
    saving: "Salvando...",
    close: "Fechar",
    copied: "Copiado!",
    copy: "Copiar",
    error: "Erro",
    success: "Sucesso",

    // Login View
    loginTitle: "Painel Administrativo",
    loginSubtitle: "Gerenciador de Personas de Áudio e Síntese de Voz em Alta Performance",
    serverUrlLabel: "URL DO SERVIDOR",
    apiKeyLabel: "API KEY",
    loginButton: "Entrar no Painel",
    loggingIn: "Conectando...",
    invalidAuthError: "Chave de API inválida ou servidor inacessível.",

    // Stats Grid
    activePersonasStat: "Personas Ativas",
    audioFormatsStat: "Formatos de Áudio",
    neuralVoicesStat: "Vozes Neurais",
    engineCacheStat: "Motor & Cache",
    engineCacheValue: "Golang + RAM LRU",

    // Personas List
    personasTitle: "Personas Registradas",
    personasSubtitle: "Gerencie conexões e configurações de áudio para cada caso de uso",
    newPersonaButton: "Nova Persona",
    emptyPersonasTitle: "Nenhuma Persona criada ainda",
    emptyPersonasDesc: "Crie sua primeira persona de áudio para começar a sintetizar com parâmetros pré-definidos.",
    createPersonaButton: "Criar Persona",

    // Persona Card
    activeStatus: "Ativa",
    noDescription: "Sem descrição informada.",
    voiceLabel: "Voz:",
    formatLabel: "Formato:",
    speedLabel: "Velocidade:",
    testAction: "Testar",
    apiAction: "API",
    editPersonaTitle: "Editar Persona",
    deletePersonaTitle: "Excluir Persona",

    // Persona Modal
    modalNewTitle: "Nova Persona de Áudio",
    modalEditTitle: "Editar Persona de Áudio",
    modalDescription: "Configure a voz neural, formato e parâmetros de síntese",
    displayNameLabel: "Nome de Exibição",
    displayNamePlaceholder: "ex: Atendente Fernanda",
    uniqueIdLabel: "Identificador Único (ID/Slug)",
    uniqueIdHint: "Apenas letras, números e traços",
    uniqueIdPlaceholder: "ex: atendente-fernanda",
    endpointPreview: "Endpoint",
    descriptionLabel: "Descrição (Finalidade)",
    descriptionPlaceholder: "Descreva o contexto ou aplicação desta persona...",
    neuralVoiceLabel: "Voz Neural",
    outputFormatLabel: "Formato de Saída",
    smartSanitizerTitle: "Higienizador Inteligente",
    smartSanitizerDesc: "Remove emojis, markdown, URLs e pontuação excessiva para fala fluida",
    saveChangesButton: "Salvar Alterações",
    createPersonaSubmit: "Criar Persona",

    // Delete Persona Modal
    deleteModalTitle: "Excluir persona",
    deleteModalDescription: "Esta ação é irreversível e excluirá permanentemente a persona",
    deleteModalConfirmLabel: "Para confirmar, digite",
    deleteModalInField: "no campo abaixo:",
    deleteModalPlaceholder: "Digite o slug para confirmar",
    deletePermanentButton: "Excluir persona",
    deletingButton: "Excluindo...",

    // Snippet Modal
    snippetModalTitle: "Integração da Persona",
    curlExampleTitle: "Exemplo cURL",

    // Playground View
    playgroundTitle: "Playground de Síntese",
    playgroundSubtitle: "Teste vozes neurais, ajuste parâmetros de áudio e analise a telemetria em tempo real",
    configCardTitle: "Configuração da Voz",
    loadPersonaPreset: "Carregar Persona",
    noneCustom: "Nenhuma (Personalizado)",
    textInputLabel: "Texto para Síntese",
    textInputPlaceholder: "Digite o texto que deseja sintetizar em voz neural...",
    synthesizeButton: "Sintetizar Áudio",
    synthesizingButton: "Sintetizando...",
    outputAndStatsTitle: "Resultado & Métricas",
    latencyTtfb: "Latência TTFB",
    totalDuration: "Tempo Total",
    audioSize: "Tamanho do Áudio",
    cacheHitStatus: "Cache Status",
    cacheHitYes: "HIT (RAM <1ms)",
    cacheHitNo: "MISS (Streaming)",
    serverHealthTitle: "Telemetria do Servidor Go",
    goroutinesCount: "Goroutines Ativas",
    allocatedRam: "Memória Alocada",
    lruCacheItems: "Itens em Cache LRU",
    lruCacheMemory: "Memória do Cache",

    // API Docs View
    apiDocsTitle: "Integrações & Exemplos de Código",
    apiDocsSubtitle: "Exemplos prontos para cURL, Python SDK, Node.js e ferramentas no-code (n8n, Typebot)",
    curlSectionTitle: "cURL - Endpoint Dedicado da Persona",
    pythonSectionTitle: "Python - SDK Oficial OpenAI",
    nodeSectionTitle: "Node.js / JavaScript Fetch",
  },

  "en-US": {
    // Navbar & Common
    statusOnline: "Online",
    logout: "Sign Out",
    selectLanguage: "Select Language",
    changeLanguageDescription: "Choose your preferred language for the dashboard interface.",
    personas: "Personas",
    playground: "TTS Playground",
    apiDocs: "API",
    version: "v1.0 (Go Native + Next.js)",
    developedBy: "Developed by",
    cancel: "Cancel",
    save: "Save",
    saving: "Saving...",
    close: "Close",
    copied: "Copied!",
    copy: "Copy",
    error: "Error",
    success: "Success",

    // Login View
    loginTitle: "Admin Dashboard",
    loginSubtitle: "High-Performance Neural Voice & Audio Persona Manager",
    serverUrlLabel: "SERVER URL",
    apiKeyLabel: "API KEY",
    loginButton: "Sign In",
    loggingIn: "Connecting...",
    invalidAuthError: "Invalid API Key or server unreachable.",

    // Stats Grid
    activePersonasStat: "Active Personas",
    audioFormatsStat: "Audio Formats",
    neuralVoicesStat: "Neural Voices",
    engineCacheStat: "Engine & Cache",
    engineCacheValue: "Golang + RAM LRU",

    // Personas List
    personasTitle: "Registered Personas",
    personasSubtitle: "Manage connections and audio settings for each use case",
    newPersonaButton: "New Persona",
    emptyPersonasTitle: "No Personas created yet",
    emptyPersonasDesc: "Create your first audio persona to start synthesizing with preset parameters.",
    createPersonaButton: "Create Persona",

    // Persona Card
    activeStatus: "Active",
    noDescription: "No description provided.",
    voiceLabel: "Voice:",
    formatLabel: "Format:",
    speedLabel: "Speed:",
    testAction: "Test",
    apiAction: "API",
    editPersonaTitle: "Edit Persona",
    deletePersonaTitle: "Delete Persona",

    // Persona Modal
    modalNewTitle: "New Audio Persona",
    modalEditTitle: "Edit Audio Persona",
    modalDescription: "Configure neural voice, output format, and synthesis parameters",
    displayNameLabel: "Display Name",
    displayNamePlaceholder: "e.g.: Support Assistant",
    uniqueIdLabel: "Unique Identifier (ID/Slug)",
    uniqueIdHint: "Letters, numbers and dashes only",
    uniqueIdPlaceholder: "e.g.: support-assistant",
    endpointPreview: "Endpoint",
    descriptionLabel: "Description (Purpose)",
    descriptionPlaceholder: "Describe the context or purpose of this persona...",
    neuralVoiceLabel: "Neural Voice",
    outputFormatLabel: "Output Format",
    smartSanitizerTitle: "Smart Sanitizer",
    smartSanitizerDesc: "Removes emojis, markdown, URLs and excessive punctuation for fluid speech",
    saveChangesButton: "Save Changes",
    createPersonaSubmit: "Create Persona",

    // Delete Persona Modal
    deleteModalTitle: "Delete persona",
    deleteModalDescription: "This action is irreversible and will permanently delete the persona",
    deleteModalConfirmLabel: "To confirm, type",
    deleteModalInField: "in the field below:",
    deleteModalPlaceholder: "Type the slug to confirm",
    deletePermanentButton: "Delete persona",
    deletingButton: "Deleting...",

    // Snippet Modal
    snippetModalTitle: "Persona Integration",
    curlExampleTitle: "cURL Example",

    // Playground View
    playgroundTitle: "Synthesis Playground",
    playgroundSubtitle: "Test neural voices, adjust audio parameters, and analyze real-time telemetry",
    configCardTitle: "Voice Configuration",
    loadPersonaPreset: "Load Persona",
    noneCustom: "None (Custom)",
    textInputLabel: "Text to Synthesize",
    textInputPlaceholder: "Type the text you want to synthesize into neural voice...",
    synthesizeButton: "Synthesize Audio",
    synthesizingButton: "Synthesizing...",
    outputAndStatsTitle: "Output & Metrics",
    latencyTtfb: "TTFB Latency",
    totalDuration: "Total Time",
    audioSize: "Audio Size",
    cacheHitStatus: "Cache Status",
    cacheHitYes: "HIT (RAM <1ms)",
    cacheHitNo: "MISS (Streaming)",
    serverHealthTitle: "Go Server Telemetry",
    goroutinesCount: "Active Goroutines",
    allocatedRam: "Allocated RAM",
    lruCacheItems: "Cached Items (LRU)",
    lruCacheMemory: "Cache Memory",

    // API Docs View
    apiDocsTitle: "Integrations & Code Examples",
    apiDocsSubtitle: "Ready-to-use examples for cURL, Python SDK, Node.js, and No-Code tools (n8n, Typebot)",
    curlSectionTitle: "cURL - Dedicated Persona Endpoint",
    pythonSectionTitle: "Python - Official OpenAI SDK",
    nodeSectionTitle: "Node.js / JavaScript Fetch",
  },

  "es-ES": {
    // Navbar & Common
    statusOnline: "En línea",
    logout: "Desconectar",
    selectLanguage: "Seleccionar Idioma",
    changeLanguageDescription: "Elija el idioma preferido para la interfaz del panel.",
    personas: "Personas",
    playground: "Playground TTS",
    apiDocs: "API",
    version: "v1.0 (Go Native + Next.js)",
    developedBy: "Desarrollado por",
    cancel: "Cancelar",
    save: "Guardar",
    saving: "Guardando...",
    close: "Cerrar",
    copied: "¡Copiado!",
    copy: "Copiar",
    error: "Error",
    success: "Éxito",

    // Login View
    loginTitle: "Panel Administrativo",
    loginSubtitle: "Administrador de Personas de Audio y Síntesis de Voz de Alto Rendimiento",
    serverUrlLabel: "URL DEL SERVIDOR",
    apiKeyLabel: "CLAVE API",
    loginButton: "Iniciar Sesión",
    loggingIn: "Conectando...",
    invalidAuthError: "Clave de API no válida o servidor inaccesible.",

    // Stats Grid
    activePersonasStat: "Personas Activas",
    audioFormatsStat: "Formatos de Audio",
    neuralVoicesStat: "Voces Neurales",
    engineCacheStat: "Motor & Caché",
    engineCacheValue: "Golang + RAM LRU",

    // Personas List
    personasTitle: "Personas Registradas",
    personasSubtitle: "Gestione conexiones y configuraciones de audio para cada caso de uso",
    newPersonaButton: "Nueva Persona",
    emptyPersonasTitle: "Aún no hay Personas creadas",
    emptyPersonasDesc: "Cree su primera persona de audio para comenzar a sintetizar con parámetros predefinidos.",
    createPersonaButton: "Crear Persona",

    // Persona Card
    activeStatus: "Activa",
    noDescription: "Sin descripción informada.",
    voiceLabel: "Voz:",
    formatLabel: "Formato:",
    speedLabel: "Velocidad:",
    testAction: "Probar",
    apiAction: "API",
    editPersonaTitle: "Editar Persona",
    deletePersonaTitle: "Eliminar Persona",

    // Persona Modal
    modalNewTitle: "Nueva Persona de Audio",
    modalEditTitle: "Editar Persona de Audio",
    modalDescription: "Configure la voz neural, formato y parámetros de síntesis",
    displayNameLabel: "Nombre para Mostrar",
    displayNamePlaceholder: "ej: Asistente Fernanda",
    uniqueIdLabel: "Identificador Único (ID/Slug)",
    uniqueIdHint: "Solo letras, números y guiones",
    uniqueIdPlaceholder: "ej: asistente-fernanda",
    endpointPreview: "Endpoint",
    descriptionLabel: "Descripción (Finalidad)",
    descriptionPlaceholder: "Describa el contexto o aplicación de esta persona...",
    neuralVoiceLabel: "Voz Neural",
    outputFormatLabel: "Formato de Salida",
    smartSanitizerTitle: "Higienizador Inteligente",
    smartSanitizerDesc: "Elimina emojis, markdown, URLs y puntuación excesiva para un habla fluida",
    saveChangesButton: "Guardar Cambios",
    createPersonaSubmit: "Crear Persona",

    // Delete Persona Modal
    deleteModalTitle: "Eliminar persona",
    deleteModalDescription: "Esta acción es irreversible y eliminará permanentemente la persona",
    deleteModalConfirmLabel: "Para confirmar, escriba",
    deleteModalInField: "en el campo inferior:",
    deleteModalPlaceholder: "Escriba el slug para confirmar",
    deletePermanentButton: "Eliminar persona",
    deletingButton: "Eliminando...",

    // Snippet Modal
    snippetModalTitle: "Integración de la Persona",
    curlExampleTitle: "Ejemplo cURL",

    // Playground View
    playgroundTitle: "Playground de Síntesis",
    playgroundSubtitle: "Pruebe voces neurales, ajuste parámetros de audio y analice la telemetría en tiempo real",
    configCardTitle: "Configuración de Voz",
    loadPersonaPreset: "Cargar Persona",
    noneCustom: "Ninguna (Personalizado)",
    textInputLabel: "Texto para Sintetizar",
    textInputPlaceholder: "Escriba el texto que desea sintetizar en voz neural...",
    synthesizeButton: "Sintetizar Audio",
    synthesizingButton: "Sintetizando...",
    outputAndStatsTitle: "Resultado & Métricas",
    latencyTtfb: "Latencia TTFB",
    totalDuration: "Tiempo Total",
    audioSize: "Tamaño de Audio",
    cacheHitStatus: "Estado de Caché",
    cacheHitYes: "HIT (RAM <1ms)",
    cacheHitNo: "MISS (Streaming)",
    serverHealthTitle: "Telemetría del Servidor Go",
    goroutinesCount: "Goroutines Activas",
    allocatedRam: "Memoria Asignada",
    lruCacheItems: "Elementos en Caché LRU",
    lruCacheMemory: "Memoria del Caché",

    // API Docs View
    apiDocsTitle: "Integraciones & Ejemplos de Código",
    apiDocsSubtitle: "Ejemplos listos para cURL, Python SDK, Node.js y herramientas no-code (n8n, Typebot)",
    curlSectionTitle: "cURL - Endpoint Dedicado de la Persona",
    pythonSectionTitle: "Python - SDK Oficial OpenAI",
    nodeSectionTitle: "Node.js / JavaScript Fetch",
  },
};

export type TranslationKeys = keyof typeof translations["pt-BR"];
