package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ArtSck/nutrition/analyzer"
	"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/nutrition"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stderr)
		return 2
	}

	switch args[0] {
	case "analyze":
		return runAnalyze(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		printHelp(stdout)
		return 0
	case "version":
		fmt.Fprintln(stdout, "nutrition v1.0 (Харчова цінність)")
		return 0
	default:
		fmt.Fprintf(stderr, "Помилка: невідома підкоманда %q\n", args[0])
		printHelp(stderr)
		return 2
	}
}

func runAnalyze(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(stderr)

	count := fs.Int("count", 0, "Кількість елементів (5..100)")
	subject := fs.String("subject", "", "Категорія продуктів (до 100 символів)")
	seed := fs.Int64("seed", 0, "Необов'язковий seed для відтворюваності")
	noColor := fs.Bool("no-color", false, "Вимкнути кольори")
	limit := fs.Float64("limit", 0, "ліміт вмісту цукру (0..100)")
	fs.Usage = func() { printHelp(stderr) }

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printHelp(stdout)
			return 0
		}
		return 2
	}

	if fs.NArg() > 0 {
		fmt.Fprintln(stderr, "Помилка: зайві позиційні аргументи")
		return 2
	}
	if *count < analyzer.MinCount || *count > analyzer.MaxCount {
		fmt.Fprintln(stderr, "Помилка: --count має бути від 5 до 100")
		return 2
	}
	trimmedSubject := strings.TrimSpace(*subject)
	if trimmedSubject == "" || len([]rune(trimmedSubject)) > 100 {
		fmt.Fprintln(stderr, "Помилка: --subject має містити від 1 до 100 символів")
		return 2
	}

	hasSeed := false
	hasLimit := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			hasSeed = true
		}
		if f.Name == "limit" {
			hasLimit = true
		}
	})

	if hasLimit && (*limit < 0 || *limit > 100) {
		fmt.Fprintln(stderr, "Помилка: --limit має бути від 0 до 100")
		return 2
	}

	var raw []float64
	var err error
	if hasSeed {
		raw, err = nutrition.GenerateWithSeed(*count, trimmedSubject, *seed)
	} else {
		raw, err = nutrition.Generate(*count, trimmedSubject)
	}
	if err != nil {
		fmt.Fprintf(stderr, "Помилка генерації: %v\n", err)
		return 1
	}

	sample, err := analyzer.NewSugarSample(trimmedSubject, raw)
	if err != nil {
		fmt.Fprintf(stderr, "Помилка даних: %v\n", err)
		return 1
	}

	rep, err := sample.Analyze(*limit)
	if err != nil {
		fmt.Fprintf(stderr, "Помилка аналізу: %v\n", err)
		return 1
	}

	printReport(stdout, rep, sample.Values(), !*noColor, hasLimit, *limit)
	return 0
}

func printReport(w io.Writer, rep analyzer.AnalysisReport, values []float64, color bool, hasLimit bool, limit float64) {
	green, yellow, reset := "", "", ""
	if color {
		green, yellow, reset = "\x1b[32m", "\x1b[33m", "\x1b[0m"
	}

	fmt.Fprintf(w, "%sКатегорія:%s %s\n", yellow, reset, rep.Subject)
	fmt.Fprintf(w, "%sДані (г цукру на 100 г):%s %v\n", green, reset, values)
	fmt.Fprintf(w, "%sКількість:%s %d\n", yellow, reset, rep.Count)
	fmt.Fprintf(w, "%sМінімум:%s %.2f\n", yellow, reset, rep.Min)
	fmt.Fprintf(w, "%sМаксимум:%s %.2f\n", yellow, reset, rep.Max)
	fmt.Fprintf(w, "%sСума:%s %.2f\n", yellow, reset, rep.Sum)
	fmt.Fprintf(w, "%sСереднє:%s %.2f\n", yellow, reset, rep.Average)
	fmt.Fprintf(w, "%sКатегорії:%s низький=%d, помірний=%d, високий=%d, дуже високий=%d\n",
		yellow, reset, rep.Low, rep.Moderate, rep.High, rep.VeryHigh)

	fmt.Fprintf(w, "%sПовтори:%s\n", yellow, reset)
	if len(rep.Repeats) == 0 {
		fmt.Fprintln(w, "  Повторів немає")
	} else {
		for _, r := range rep.Repeats {
			fmt.Fprintf(w, "  %.2f зустрічається %d раз(и)\n", r.Value, r.Count)
		}
	}
	if hasLimit {
		fmt.Fprintf(w, "Перевищують ліміт (%.2f): %d\n", limit, rep.AboveLimit)
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "nutrition — CLI-аналізатор харчової цінності")
	fmt.Fprintln(w, "Використання:")
	fmt.Fprintln(w, "  nutrition analyze --count <5..100> --subject <назва> [--seed <число>] [--no-color] [--limit <0..100>]")
	fmt.Fprintln(w, "  nutrition help | version")
}
