package edgetts

import (
	"strings"
	"testing"
)

func TestSpeedToRate(t *testing.T) {
	tests := []struct {
		name      string
		speed     float64
		expected  string
		expectErr bool
	}{
		{name: "Velocidade Normal 1.0", speed: 1.0, expected: "+0%", expectErr: false},
		{name: "Velocidade Acelerada 1.5", speed: 1.5, expected: "+50%", expectErr: false},
		{name: "Velocidade Dobrada 2.0", speed: 2.0, expected: "+100%", expectErr: false},
		{name: "Velocidade Lenta 0.75", speed: 0.75, expected: "-25%", expectErr: false},
		{name: "Velocidade Mínima 0.25", speed: 0.25, expected: "-75%", expectErr: false},
		{name: "Abaixo do Mínimo", speed: 0.1, expected: "", expectErr: true},
		{name: "Acima do Máximo", speed: 2.5, expected: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SpeedToRate(tt.speed)
			if (err != nil) != tt.expectErr {
				t.Fatalf("SpeedToRate(%f) erro inesperado: %v", tt.speed, err)
			}
			if got != tt.expected {
				t.Errorf("SpeedToRate(%f) = %q, esperado %q", tt.speed, got, tt.expected)
			}
		})
	}
}

func TestBuildSSML(t *testing.T) {
	ssml := BuildSSML("Olá & <Mundo>", "pt-BR-ThalitaMultilingualNeural", "+20%", "+0Hz", "+0%", "pt-BR")

	if !strings.Contains(ssml, "Olá &amp; &lt;Mundo&gt;") {
		t.Errorf("BuildSSML não escapou entidades XML corretamente: %s", ssml)
	}

	if !strings.Contains(ssml, "name='pt-BR-ThalitaMultilingualNeural'") {
		t.Errorf("BuildSSML voz incorreta: %s", ssml)
	}

	if !strings.Contains(ssml, "rate='+20%'") {
		t.Errorf("BuildSSML taxa de velocidade incorreta: %s", ssml)
	}
}
