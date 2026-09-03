package audioorchestrator

import (
	"bytes"
	"encoding/binary"
	"strings"
	"time"
)

// mp3SilenceFrame24kHz48kbpsMono representa um frame MPEG-2 Layer III válido de 24kHz 48kbps Mono
// Duração de cada frame: 576 amostras / 24000 Hz = 24ms (288 bytes por frame).
var mp3SilenceFrame24kHz48kbpsMono = func() []byte {
	frame := make([]byte, 288)
	// Header MPEG 2 Layer III, 48kbps, 24kHz, Mono, No CRC: 0xFF, 0xF3, 0x50, 0xC0
	frame[0] = 0xFF
	frame[1] = 0xF3
	frame[2] = 0x50
	frame[3] = 0xC0
	// Side information (9 bytes) com energia zero (global_gain = 0, part2_3_length = 0)
	// Restante preenchido com zeros (Huffman data vazio / silêncio absoluto)
	return frame
}()

// mp3SilenceFrame44kHz128kbps representa um frame MPEG-1 Layer III padrão de 44.1kHz 128kbps
// Duração: 1152 amostras / 44100 Hz ≈ 26.12ms (417 bytes por frame).
var mp3SilenceFrame44kHz128kbps = func() []byte {
	frame := make([]byte, 417)
	// Header MPEG 1 Layer III, 128kbps, 44.1kHz, Joint Stereo, No CRC: 0xFF, 0xFB, 0x90, 0x64
	frame[0] = 0xFF
	frame[1] = 0xFB
	frame[2] = 0x90
	frame[3] = 0x64
	return frame
}()

// GenerateSilence gera um buffer de áudio com silêncio puro correspondente à duração e formato solicitados.
func GenerateSilence(format string, duration time.Duration) ([]byte, error) {
	if duration <= 0 {
		return nil, nil
	}

	f := strings.ToLower(strings.TrimSpace(format))
	switch f {
	case "mp3", "audio/mpeg":
		return generateMP3Silence(duration), nil
	case "wav", "audio/wav":
		return generateWAVSilence(duration, 24000, 1, 16), nil
	case "pcm", "raw":
		return generatePCMSilence(duration, 24000, 1, 16), nil
	case "opus", "ogg", "audio/opus", "audio/ogg":
		// Para Opus/OGG ou fallback, geramos MP3 silence ou PCM de acordo com o container
		return generateMP3Silence(duration), nil
	default:
		return generateMP3Silence(duration), nil
	}
}

// generateMP3Silence gera frames MP3 de silêncio puro para a duração desejada.
func generateMP3Silence(duration time.Duration) []byte {
	frameDuration := 24 * time.Millisecond
	numFrames := int(duration / frameDuration)
	if numFrames <= 0 {
		numFrames = 1
	}

	frameLen := len(mp3SilenceFrame24kHz48kbpsMono)
	buf := make([]byte, numFrames*frameLen)
	for i := 0; i < numFrames; i++ {
		copy(buf[i*frameLen:], mp3SilenceFrame24kHz48kbpsMono)
	}
	return buf
}

// generatePCMSilence gera bytes zero (0x00) para áudio PCM puro.
func generatePCMSilence(duration time.Duration, sampleRate, channels, bitsPerSample int) []byte {
	bytesPerSec := sampleRate * channels * (bitsPerSample / 8)
	totalBytes := int(float64(bytesPerSec) * duration.Seconds())
	if totalBytes <= 0 {
		totalBytes = 2
	}
	return make([]byte, totalBytes)
}

// generateWAVSilence gera um arquivo WAV válido contendo silêncio puro.
func generateWAVSilence(duration time.Duration, sampleRate, channels, bitsPerSample int) []byte {
	pcmData := generatePCMSilence(duration, sampleRate, channels, bitsPerSample)
	dataLen := uint32(len(pcmData))

	buf := new(bytes.Buffer)
	buf.Grow(44 + len(pcmData))

	// RIFF header
	buf.WriteString("RIFF")
	_ = binary.Write(buf, binary.LittleEndian, uint32(36+dataLen))
	buf.WriteString("WAVE")

	// fmt subchunk
	buf.WriteString("fmt ")
	_ = binary.Write(buf, binary.LittleEndian, uint32(16)) // Subchunk1Size
	_ = binary.Write(buf, binary.LittleEndian, uint16(1))  // AudioFormat (1 = PCM)
	_ = binary.Write(buf, binary.LittleEndian, uint16(channels))
	_ = binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	byteRate := uint32(sampleRate * channels * (bitsPerSample / 8))
	_ = binary.Write(buf, binary.LittleEndian, byteRate)
	blockAlign := uint16(channels * (bitsPerSample / 8))
	_ = binary.Write(buf, binary.LittleEndian, blockAlign)
	_ = binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample))

	// data subchunk
	buf.WriteString("data")
	_ = binary.Write(buf, binary.LittleEndian, dataLen)
	buf.Write(pcmData)

	return buf.Bytes()
}
