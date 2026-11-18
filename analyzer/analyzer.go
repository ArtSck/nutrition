package analyzer

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	MinSampleCount   = 5
	MaxSampleCount   = 100
	MaxSubjectLength = 100
	MinSugarValue    = 0.0
	MaxSugarValue    = 100.0
)

type SugarContent float64

type CategoryCount struct {
	Name  string
	Count int
}

type RepeatFrequency struct {
	Value SugarContent
	Count int
}

type SugarSample struct {
	subject string
	values  []SugarContent
}

type AnalysisReport struct {
	Subject    string
	Count      int
	Min        SugarContent
	Max        SugarContent
	Sum        float64
	Average    float64
	Categories []CategoryCount
	Repeats    []RepeatFrequency
}

func NewSugarSample(subject string, rawValues []float64) (SugarSample, error) {
	trimmed := strings.TrimSpace(subject)
	if trimmed == "" {
		return SugarSample{}, errors.New("предметний параметр не може бути порожнім")
	}
	if len([]rune(trimmed)) > MaxSubjectLength {
		return SugarSample{}, fmt.Errorf("предметний параметр перевищує %d символів", MaxSubjectLength)
	}
	if len(rawValues) < MinSampleCount || len(rawValues) > MaxSampleCount {
		return SugarSample{}, fmt.Errorf("кількість значень має бути від %d до %d, отримано %d", MinSampleCount, MaxSampleCount, len(rawValues))
	}
	copied := make([]SugarContent, len(rawValues))
	for i, v := range rawValues {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return SugarSample{}, fmt.Errorf("значення [%d] не є скінченним числом", i)
		}
		if v < MinSugarValue || v > MaxSugarValue {
			return SugarSample{}, fmt.Errorf("значення %.2f виходить за допустимий діапазон [%.1f, %.1f]", v, MinSugarValue, MaxSugarValue)
		}
		copied[i] = SugarContent(v)
	}
	return SugarSample{subject: trimmed, values: copied}, nil
}

func (s SugarSample) Subject() string { return s.subject }

func (s SugarSample) Values() []SugarContent {
	out := make([]SugarContent, len(s.values))
	copy(out, s.values)
	return out
}

func (s SugarSample) Analyze() (AnalysisReport, error) {
	if len(s.values) < MinSampleCount {
		return AnalysisReport{}, errors.New("набір даних не ініціалізовано або має недостатньо елементів")
	}
	minVal, maxVal := s.values[0], s.values[0]
	sum := 0.0
	low, moderate, high, veryHigh := 0, 0, 0, 0
	freqMap := make(map[SugarContent]int)

	for _, v := range s.values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sum += float64(v)
		switch {
		case v <= 5.0:
			low++
		case v <= 15.0:
			moderate++
		case v <= 30.0:
			high++
		default:
			veryHigh++
		}
		freqMap[v]++
	}

	categories := []CategoryCount{
		{Name: "Низький (до 5)", Count: low},
		{Name: "Помірний (5.1-15)", Count: moderate},
		{Name: "Високий (15.1-30)", Count: high},
		{Name: "Дуже високий (понад 30)", Count: veryHigh},
	}

	repeats := make([]RepeatFrequency, 0)
	for val, freq := range freqMap {
		if freq > 1 {
			repeats = append(repeats, RepeatFrequency{Value: val, Count: freq})
		}
	}
	sort.Slice(repeats, func(i, j int) bool { return repeats[i].Value < repeats[j].Value })

	return AnalysisReport{
		Subject:    s.subject,
		Count:      len(s.values),
		Min:        minVal,
		Max:        maxVal,
		Sum:        sum,
		Average:    sum / float64(len(s.values)),
		Categories: categories,
		Repeats:    repeats,
	}, nil
}
