package daysteps

import (
	"errors"
	"fmt"
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
	// Разделяем строку на части
	parts := strings.Split(data, ",")

	// Проверяем, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования шагов: %w", err)
	}

	// Проверяем, что шагов больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о шагах и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	// Проверяем, что шагов больше 0
	if steps <= 0 {
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
