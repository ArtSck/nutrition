package analyzer

import (
	"math"
	"slices"
	"testing"

	"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/nutrition"
)

func TestNewSugarSample(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		values  []float64
		wantErr bool
	}{
		{"valid", "breakfast cereals", []float64{2.5, 11.8, 28.4, 61.7, 11.8}, false},
		{"empty subject", "   ", []float64{1, 2, 3, 4, 5}, true},
		{"too few values", "yogurt", []float64{1, 2, 3, 4}, true},
		{"out of range", "candies", []float64{10, 20, 30, 40, 105}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSugarSample(tc.subject, tc.values)
			if (err != nil) != tc.wantErr {
				t.Errorf("очікували помилку=%v, отримали %v", tc.wantErr, err)
			}
		})
	}
}

func TestAnalyze(t *testing.T) {
	s, err := NewSugarSample("breakfast cereals", []float64{2.5, 11.8, 28.4, 61.7, 11.8})
	if err != nil {
		t.Fatalf("несподівана помилка: %v", err)
	}

	rep, err := s.Analyze(0)
	if err != nil {
		t.Fatalf("помилка Analyze: %v", err)
	}

	if rep.Min != 2.5 || rep.Max != 61.7 || math.Abs(rep.Sum-116.2) > 1e-9 || math.Abs(rep.Average-23.24) > 1e-9 {
		t.Errorf("неправильна статистика: %+v", rep)
	}
	if rep.Low != 1 || rep.Moderate != 2 || rep.High != 1 || rep.VeryHigh != 1 {
		t.Errorf("неправильні категорії: %+v", rep)
	}
	wantRepeats := []RepeatItem{{Value: 11.8, Count: 2}}
	if !slices.Equal(rep.Repeats, wantRepeats) {
		t.Errorf("очікували повтори %+v, отримали %+v", wantRepeats, rep.Repeats)
	}
}

func TestGenerateWithSeed(t *testing.T) {
	a, err := nutrition.GenerateWithSeed(5, "breakfast cereals", 42)
	if err != nil {
		t.Fatalf("GenerateWithSeed: %v", err)
	}
	b, _ := nutrition.GenerateWithSeed(5, "breakfast cereals", 42)
	if !slices.Equal(a, b) {
		t.Errorf("однаковий seed дав різні дані")
	}
}
