# 🔊 Pasta de Efeitos Sonoros Customizados (SFX)

Coloque seus arquivos de áudio de efeitos sonoros nesta pasta para que o **EdgeGo Voice** os utilize automaticamente nas suas tags de texto.

---

## 📁 Nomenclatura dos Arquivos

O sistema busca os arquivos na seguinte ordem de prioridade:

1. **Específico por Voz**: `{NomeDaVoz}_{efeito}.mp3` ou `{NomeDaVoz}_{efeito}.wav`
   - Exemplo: `pt-BR-FranciscaNeural_pigarro.mp3`
   - Exemplo: `Francisca_risada.mp3`
   - Exemplo: `Thalita_tosse.mp3`

2. **Genérico (Qualquer Voz)**: `{efeito}.mp3` ou `{efeito}.wav`
   - Exemplo: `pigarro.mp3`
   - Exemplo: `tosse.mp3`
   - Exemplo: `risada.mp3`
   - Exemplo: `suspiro.mp3`
   - Exemplo: `respiracao.mp3`
   - Exemplo: `aplausos.mp3`
   - Exemplo: `notificacao.mp3`

---

## 🏷️ Como Usar no Texto / API:

Qualquer arquivo que você colocar aqui fica disponível imediatamente via tag de texto:

- Se você tiver `sfx/pigarro.mp3` $\rightarrow$ use no texto: `Olá [som:pigarro] como vai?`
- Se você tiver `sfx/risada.mp3` $\rightarrow$ use no texto: `Isso foi muito bom [som:risada] mesmo!`
- Se você tiver `sfx/notificacao.wav` $\rightarrow$ use no texto: `[som:notificacao] Você recebeu uma mensagem.`
- Se você tiver `sfx/meu_efeito.mp3` $\rightarrow$ use no texto: `[som:meu_efeito] texto aqui.`

---

## ⚡ Formatos Suportados:
- `.mp3` (Recomendado: 24kHz ou 44.1kHz mono/stereo)
- `.wav`
- `.opus` / `.ogg`
- `.aac`
