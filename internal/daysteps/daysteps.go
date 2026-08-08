package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Проверяем, что строка не пустая
	if data == "" {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Разделяем строку на части
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Удаляем пробелы в начале и конце
	stepsStr := strings.TrimSpace(parts[0])
	durationStr := strings.TrimSpace(parts[1])

	// Проверяем, что строка с шагами не пустая
	if stepsStr == "" || stepsStr == "+" || stepsStr == "-" {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования шагов: %w", err)
	}

	// Проверяем, что шагов больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// Проверяем, что продолжительность не пустая
	if durationStr == "" {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о шагах и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка парсинга данных:", err)
		return ""
	}

	// Проверяем, что шагов больше 0
	if steps <= 0 {
		return ""
	}

	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength
	// Переводим в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем калории (используем функцию из пакета spentcalories)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка при расчете калорий:", err)
		return ""
	}

	// Формируем строку результата
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories)
}
