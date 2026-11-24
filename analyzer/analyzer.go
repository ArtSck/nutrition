package analyzer

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	MinCount = 5
	MaxCount = 100
)

type SugarSample struct {
	subject string
	values  []float64
}

type RepeatItem struct {
	Value float64
	Count int
}

type AnalysisReport struct {
	Subject    string
	Count      int
	Min        float64
	Max        float64
	Sum        float64
	Average    float64
	Low        int
	Moderate   int
	High       int
	VeryHigh   int
	Repeats    []RepeatItem
	AboveLimit int
}

func NewSugarSample(subject string, raw []float64) (SugarSample, error) {
	trimmed := strings.TrimSpace(subject)
	if trimmed == "" || len([]rune(trimmed)) > 100 {
		return SugarSample{}, errors.New("категорія має містити від 1 до 100 символів")
	}
	if len(raw) < MinCount || len(raw) > MaxCount {
		return SugarSample{}, fmt.Errorf("кількість значень має бути від %d до %d", MinCount, MaxCount)
	}

	copied := make([]float64, len(raw))
	for i, v := range raw {
		if v < 0 || v > 100 {
			return SugarSample{}, fmt.Errorf("значення %.2f поза межами 0..100", v)
		}
		copied[i] = v
	}

	return SugarSample{subject: trimmed, values: copied}, nil
}

func (s SugarSample) Values() []float64 {
	out := make([]float64, len(s.values))
	copy(out, s.values)
	return out
}

func (s SugarSample) Analyze(limit float64) (AnalysisReport, error) {
	if len(s.values) < MinCount {
		return AnalysisReport{}, errors.New("некоректний стан вибірки")
	}

	rep := AnalysisReport{
		Subject: s.subject,
		Count:   len(s.values),
		Min:     s.values[0],
		Max:     s.values[0],
	}

	freq := make(map[float64]int)
	for _, v := range s.values {
		if v < rep.Min {
			rep.Min = v
		}
		if v > rep.Max {
			rep.Max = v
		}
		rep.Sum += v

		if v > limit {
			rep.AboveLimit++
		}

		switch {
		case v <= 5:
			rep.Low++
		case v <= 15:
			rep.Moderate++
		case v <= 30:
			rep.High++
		default:
			rep.VeryHigh++
		}
		freq[v]++
	}

	rep.Average = rep.Sum / float64(rep.Count)

	for val, cnt := range freq {
		if cnt > 1 {
			rep.Repeats = append(rep.Repeats, RepeatItem{Value: val, Count: cnt})
		}
	}
	sort.Slice(rep.Repeats, func(i, j int) bool { return rep.Repeats[i].Value < rep.Repeats[j].Value })

	return rep, nil
}
