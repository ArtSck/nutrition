package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/nutrition"
)

func main() {
	var count int
	var category string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введіть кількість елементів (від 5 до 100): ")
	if !scanner.Scan() {
		fmt.Println("\x1b[31mПомилка: кількість елементів не може бути порожньою\x1b[0m")
		return
	}
	countText := strings.TrimSpace(scanner.Text())
	count, err := strconv.Atoi(countText)
	if err != nil || count < 5 || count > 100 {
		fmt.Println("\x1b[31mПомилка: кількість елементів має бути цілим числом від 5 до 100\x1b[0m")
		return
	}

	fmt.Print("Введіть категорію (наприклад, breakfast cereals): ")
	if !scanner.Scan() {
		fmt.Println("\x1b[31mПомилка: категорія не може бути порожньою\x1b[0m")
		return
	}

	category = strings.TrimSpace(scanner.Text())
	if category == "" {
		fmt.Println("\x1b[31mПомилка: категорія не може бути порожньою\x1b[0m")
		return
	}

	const seed int64 = 42
	values, err := nutrition.GenerateWithSeed(count, category, seed)
	if err != nil {
		fmt.Println("\x1b[31mПомилка:", err, "\x1b[0m")
		return
	}
	secondValues, err := nutrition.GenerateWithSeed(count, category, seed)
	if err != nil {
		fmt.Println("\x1b[31mПомилка:", err, "\x1b[0m")
		return
	}

	fmt.Println("\x1b[32mЗгенеровані дані (г цукру на 100 г): \x1b[0m", values)
	fmt.Println("\x1b[32mПовторний виклик дав такий самий зріз:", slices.Equal(values, secondValues), "\x1b[0m")

	min := values[0]
	max := values[0]
	sum := 0.0

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	average := sum / float64(len(values))

	fmt.Println("\x1b[33mМінімум: ", min)
	fmt.Println("\x1b[33mМаксимум: ", max)
	fmt.Printf("\x1b[33mСума: %.2f\x1b[0m\n", sum)
	fmt.Printf("\x1b[33mСереднє: %.2f\x1b[0m\n", average)

	low, moderate, high, veryHigh := 0, 0, 0, 0
	for _, v := range values {
		switch {
		case v <= 5:
			low++
		case v <= 15:
			moderate++
		case v <= 30:
			high++
		default:
			veryHigh++
		}
	}

	fmt.Println("\x1b[33m- Категоризація вмісту цукру:\x1b[0m")
	fmt.Println("\x1b[33m- Низький (до 5): \x1b[0m", low)
	fmt.Println("\x1b[33m- Помірний (5.1-15): \x1b[0m", moderate)
	fmt.Println("\x1b[33m- Високий (15.1-30): \x1b[0m", high)
	fmt.Println("\x1b[33m- Дуже високий (понад 30): \x1b[0m", veryHigh)

	countsMap := make(map[float64]int)
	for _, v := range values {
		countsMap[v]++
	}

	fmt.Println("\x1b[33m- Повтори значень через map:\x1b[0m")
	for val, freq := range countsMap {
		if freq > 1 {
			fmt.Println("\x1b[33mЗначення \x1b[0m", val, "\x1b[33mзустрічається\x1b[0m", freq, "\x1b[33mраз\x1b[0m")
		}
	}

	first, err := nutrition.GenerateWithSeed(count, category, 42)
	if err != nil {
		fmt.Println("\x1b[31mПомилка:", err, "\x1b[0m")
		return
	}

	second, err := nutrition.GenerateWithSeed(count, category, 43)
	if err != nil {
		fmt.Println("\x1b[31mПомилка:", err, "\x1b[0m")
		return
	}

	fmt.Println("\x1b[36mSeed 42:\x1b[0m", first)
	fmt.Println("\x1b[36mSeed 43:\x1b[0m", second)
	fmt.Println("\x1b[36mЗрізи різні:\x1b[0m", !slices.Equal(first, second))

	fmt.Print("Введіть бажаний ліміт цукру у продукті (від 0.1 до 100 г на 100 г): ")
	if !scanner.Scan() {
		fmt.Println("\x1b[31mПомилка: ліміт цукру не може бути порожнім\x1b[0m")
		return
	}
	limitText := strings.TrimSpace(scanner.Text())
	sugarLimit, err := strconv.ParseFloat(limitText, 64)
	if err != nil || math.IsNaN(sugarLimit) || math.IsInf(sugarLimit, 0) || sugarLimit < 0.1 || sugarLimit > 100 {
		fmt.Println("\x1b[31mПомилка: ліміт цукру має бути числом від 0.1 до 100\x1b[0m")
		return
	}

	satisfying := make([]float64, 0, len(values))
	for _, value := range values {
		if value <= sugarLimit {
			satisfying = append(satisfying, value)
		}
	}

	percentage := float64(len(satisfying)) / float64(len(values)) * 100
	fmt.Println("\x1b[33mЗначення, що відповідають ліміту:\x1b[0m", satisfying)
	fmt.Printf("\x1b[33mЗадовільняючих результатів: %d із %d (%.2f%%)\x1b[0m\n", len(satisfying), len(values), percentage)
}
