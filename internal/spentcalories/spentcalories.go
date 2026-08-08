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
	// Разделяем строку на части
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 3
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных тренировки")
	}

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования шагов: %w", err)
	}

	// Вид активности (проверяем корректность)
	activityType := parts[1]
	if activityType != "Бег" && activityType != "Ходьба" {
		return 0, "", 0, errors.New("неизвестный тип активности")
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// Длина шага = рост * коэффициент
	stepLength := height * stepLengthCoefficient
	// Дистанция в метрах = шаги * длина шага
	distanceMeters := float64(steps) * stepLength
	// Переводим в километры
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	// Переводим продолжительность в часы
	hours := duration.Hours()

	// Вычисляем и возвращаем среднюю скорость
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем данные из строки
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга тренировки:", err)
		return "", err
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)

	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Вычисляем калории в зависимости от типа тренировки
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

	// Формируем строку результата
	result := fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f",
		activityType, duration.Hours(), dist, speed, calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность параметров
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

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	minutes := duration.Minutes()

	// Рассчитываем калории: (вес * скорость * минуты) / 60
	calories := (weight * speed * minutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем корректность параметров
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

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	minutes := duration.Minutes()

	// Рассчитываем калории с корректирующим коэффициентом
	calories := (weight * speed * minutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
