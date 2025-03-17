package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

func parseTraining(data string) (int, string, time.Duration, error) {

	// Parse data string to slice
	splittedData := strings.Split(data, ",")
	if len(splittedData) != 3 {
		return 0, "", 0, errors.New("неверные данные: шаги и продолжительность прогулки")
	}

	// Converting steps info into int
	steps, err := strconv.Atoi(splittedData[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("колличество шагов не может быть равно 0")
	}

	// Converting walking duration into time.Duration format
	walkDuration, err := time.ParseDuration(splittedData[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, splittedData[1], walkDuration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {

	// Calc distance in kilometers
	distance := (float64(steps) * lenStep) / mInKm

	return distance
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {

	// Calc average speed
	if duration <= 0 {
		return 0
	}
	d := distance(steps)
	avgSpeed := d / duration.Hours()

	return avgSpeed
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {

	// Returning formated data about training
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}
	var result string
	switch activityType {
	case "Бег":
		distance := distance(steps)
		avgSpeed := meanSpeed(steps, duration)
		cal := RunningSpentCalories(steps, weight, duration)
		result = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activityType, duration.Hours(), distance, avgSpeed, cal)

	case "Ходьба":
		distance := distance(steps)
		avgSpeed := meanSpeed(steps, duration)
		cal := WalkingSpentCalories(steps, weight, height, duration)
		result = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activityType, duration.Hours(), distance, avgSpeed, cal)

	default:
		result = "неизвестный тип тренировки"
	}
	return result
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {

	// Calc cal spent when running
	avgSpeed := meanSpeed(steps, duration)
	calWasted := ((runningCaloriesMeanSpeedMultiplier * avgSpeed) - runningCaloriesMeanSpeedShift) * weight

	return calWasted
}

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {

	// Calc cal spent when running
	avgSpeed := meanSpeed(steps, duration)
	calWasted := ((walkingCaloriesWeightMultiplier * weight) + (avgSpeed*avgSpeed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
	return calWasted
}
