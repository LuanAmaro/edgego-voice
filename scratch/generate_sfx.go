package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
)

func writeWAV(filename string, samples []int16, sampleRate int) error {
	buf := new(bytes.Buffer)
	numChannels := uint16(1)
	bitsPerSample := uint16(16)
	byteRate := uint32(sampleRate * int(numChannels) * int(bitsPerSample/8))
	blockAlign := uint16(numChannels * (bitsPerSample / 8))
	dataSize := uint32(len(samples) * 2)

	buf.WriteString("RIFF")
	_ = binary.Write(buf, binary.LittleEndian, uint32(36+dataSize))
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	_ = binary.Write(buf, binary.LittleEndian, uint32(16))
	_ = binary.Write(buf, binary.LittleEndian, uint16(1)) // PCM
	_ = binary.Write(buf, binary.LittleEndian, numChannels)
	_ = binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	_ = binary.Write(buf, binary.LittleEndian, byteRate)
	_ = binary.Write(buf, binary.LittleEndian, blockAlign)
	_ = binary.Write(buf, binary.LittleEndian, bitsPerSample)
	buf.WriteString("data")
	_ = binary.Write(buf, binary.LittleEndian, dataSize)

	for _, s := range samples {
		_ = binary.Write(buf, binary.LittleEndian, s)
	}

	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func main() {
	sampleRate := 24000
	sfxDir := "sfx"

	// 1. Pigarro Feminino (~0.45s, 2 impulsos guturais em 270Hz)
	{
		dur := 0.45
		total := int(float64(sampleRate) * dur)
		pcm := make([]int16, total)
		r := rand.New(rand.NewSource(42))
		for i := 0; i < total; i++ {
			t := float64(i) / float64(sampleRate)
			// Dois pulsos centrados em t=0.10s e t=0.28s
			env := math.Exp(-math.Pow((t-0.10)/0.035, 2)) + 0.85*math.Exp(-math.Pow((t-0.28)/0.045, 2))
			noise := r.Float64()*2.0 - 1.0
			tone := math.Sin(2*math.Pi*270*t) + 0.4*math.Sin(2*math.Pi*540*t)
			val := (noise*0.65 + tone*0.35) * env * 16000
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}
		_ = writeWAV(filepath.Join(sfxDir, "pigarro-feminino.wav"), pcm, sampleRate)
		_ = writeWAV(filepath.Join(sfxDir, "pigarro.wav"), pcm, sampleRate)
	}

	// 2. Pigarro Masculino (~0.50s, tom mais grave em 170Hz)
	{
		dur := 0.50
		total := int(float64(sampleRate) * dur)
		pcm := make([]int16, total)
		r := rand.New(rand.NewSource(84))
		for i := 0; i < total; i++ {
			t := float64(i) / float64(sampleRate)
			env := math.Exp(-math.Pow((t-0.12)/0.04, 2)) + 0.9*math.Exp(-math.Pow((t-0.32)/0.05, 2))
			noise := r.Float64()*2.0 - 1.0
			tone := math.Sin(2*math.Pi*170*t) + 0.4*math.Sin(2*math.Pi*340*t)
			val := (noise*0.60 + tone*0.40) * env * 18000
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}
		_ = writeWAV(filepath.Join(sfxDir, "pigarro-masculino.wav"), pcm, sampleRate)
	}

	// 3. Risada Feminina (~0.85s, 3 gargalhadas "ha-ha-ha" decrescentes em 420Hz)
	{
		dur := 0.85
		total := int(float64(sampleRate) * dur)
		pcm := make([]int16, total)
		r := rand.New(rand.NewSource(105))
		for i := 0; i < total; i++ {
			t := float64(i) / float64(sampleRate)
			// Três pulsos de risada: t=0.08s, t=0.30s, t=0.55s com atenuação gradual
			env := math.Exp(-math.Pow((t-0.08)/0.05, 2)) +
				0.75*math.Exp(-math.Pow((t-0.30)/0.06, 2)) +
				0.50*math.Exp(-math.Pow((t-0.55)/0.07, 2))
			noise := r.Float64()*2.0 - 1.0
			// Modulação de pitch da risada
			freq := 420.0 - 40.0*t
			tone := math.Sin(2*math.Pi*freq*t) + 0.3*math.Sin(2*math.Pi*freq*2*t)
			val := (tone*0.75 + noise*0.25) * env * 17000
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}
		_ = writeWAV(filepath.Join(sfxDir, "risada-feminina.wav"), pcm, sampleRate)
		_ = writeWAV(filepath.Join(sfxDir, "risada.wav"), pcm, sampleRate)
		_ = writeWAV(filepath.Join(sfxDir, "risos.wav"), pcm, sampleRate)
		_ = writeWAV(filepath.Join(sfxDir, "risos-feminina.wav"), pcm, sampleRate)
	}

	// 4. Risada Masculina (~0.90s, tom mais grave em 220Hz)
	{
		dur := 0.90
		total := int(float64(sampleRate) * dur)
		pcm := make([]int16, total)
		r := rand.New(rand.NewSource(210))
		for i := 0; i < total; i++ {
			t := float64(i) / float64(sampleRate)
			env := math.Exp(-math.Pow((t-0.10)/0.06, 2)) +
				0.75*math.Exp(-math.Pow((t-0.34)/0.07, 2)) +
				0.50*math.Exp(-math.Pow((t-0.60)/0.08, 2))
			noise := r.Float64()*2.0 - 1.0
			freq := 220.0 - 30.0*t
			tone := math.Sin(2*math.Pi*freq*t) + 0.35*math.Sin(2*math.Pi*freq*2*t)
			val := (tone*0.70 + noise*0.30) * env * 19000
			pcm[i] = int16(math.Max(-32768, math.Min(32767, val)))
		}
		_ = writeWAV(filepath.Join(sfxDir, "risada-masculina.wav"), pcm, sampleRate)
		_ = writeWAV(filepath.Join(sfxDir, "risos-masculino.wav"), pcm, sampleRate)
	}

	// 5. Garantir cópia genérica de respiracao.wav caso respiracao-feminina.wav exista
	if data, err := os.ReadFile(filepath.Join(sfxDir, "respiracao-feminina.wav")); err == nil && len(data) > 0 {
		_ = os.WriteFile(filepath.Join(sfxDir, "respiracao.wav"), data, 0644)
	}

	fmt.Println("Todos os arquivos de SFX foram gerados com sucesso!")
}
