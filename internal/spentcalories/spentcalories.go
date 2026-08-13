package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных тренировки")
	}

	steps, err := strconv.Atoi(parts[0])
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования шагов: %w", err)
	}

	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	hours := duration.Hours()

	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга тренировки:", err)
		return "", err
	}

	dist := distance(steps, height)

	speed := meanSpeed(steps, height, duration)

	var calories float64
	var errCal error

	switch activityType {
	case "Бег":
		calories, errCal = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, errCal = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if errCal != nil {
		log.Println("Ошибка при расчете калорий:", errCal)
		return "", errCal
	}

	result := fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		activityType, duration.Hours(), dist, speed, calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)

	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)

	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
