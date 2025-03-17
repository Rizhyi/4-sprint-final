package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rizhyi/4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {

	// Parsing data string to slice
	splittedData := strings.Split(data, ",")
	if len(splittedData) != 2 {
		return 0, 0, errors.New("неверные данные: шаги и продолжительность прогулки")
	}

	// Converting steps info into int
	steps, err := strconv.Atoi(splittedData[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("колличество шагов не может быть равно 0")
	}

	// Converting walking duration into time.Duration format
	walkDuration, err := time.ParseDuration(splittedData[1])
	if err != nil {
		return 0, 0, err
	}

	return steps, walkDuration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {

	// Getting steps amd walking time data
	steps, walkDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if steps <= 0 {
		return ""
	}

	// Calc distance in kilometers
	distance := (float64(steps) * StepLength) / 1000

	// Calc wasted calories
	wastedCalories := spentcalories.WalkingSpentCalories(steps, weight, height, walkDuration)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distance, wastedCalories)
}
