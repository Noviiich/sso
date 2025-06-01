package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

// PrettyHandlerOptions содержит настройки для PrettyHandler
type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions // Стандартные опции slog (уровень, добавление источника и т.д.)
}

// PrettyHandler - кастомный обработчик логов с красивым цветным выводом
type PrettyHandler struct {
	opts         PrettyHandlerOptions // Настройки обработчика
	slog.Handler                      // Встроенный JSONHandler для базовой функциональности
	l            *stdLog.Logger       // Стандартный логгер для вывода форматированных сообщений
	attrs        []slog.Attr          // Атрибуты, сохраненные через WithAttrs()
}

// NewPrettyHandler создает новый экземпляр PrettyHandler
func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer, // Поток вывода (обычно os.Stdout или файл)
) *PrettyHandler {
	h := &PrettyHandler{
		// Используем JSONHandler как базовый обработчик для совместимости
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		// Создаем стандартный логгер без префикса и флагов для чистого вывода
		l: stdLog.New(out, "", 0),
	}

	return h
}

// Handle обрабатывает одну запись лога и выводит её в красивом формате
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// Получаем строковое представление уровня лога
	level := r.Level.String() + ":"

	// Раскрашиваем уровень в зависимости от его типа
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level) // DEBUG: фиолетовый
	case slog.LevelInfo:
		level = color.BlueString(level) // INFO: синий
	case slog.LevelWarn:
		level = color.YellowString(level) // WARN: желтый
	case slog.LevelError:
		level = color.RedString(level) // ERROR: красный
	}

	// Создаем map для сбора всех атрибутов
	fields := make(map[string]interface{}, r.NumAttrs())

	// Собираем атрибуты из текущей записи лога
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any() // Извлекаем значение атрибута
		return true                   // Продолжаем итерацию
	})

	// Добавляем сохраненные атрибуты из WithAttrs()
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	// Если есть атрибуты, форматируем их как JSON с отступами
	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err // Возвращаем ошибку если не удалось сериализовать
		}
	}

	// Форматируем время в читаемом виде [15:05:05.000]
	timeStr := r.Time.Format("[15:05:05.000]")
	// Раскрашиваем сообщение в голубой цвет
	msg := color.CyanString(r.Message)

	// Выводим финальную строку лога:
	// [время] УРОВЕНЬ: сообщение {json-атрибуты}
	h.l.Println(
		timeStr,                      // Время
		level,                        // Цветной уровень
		msg,                          // Цветное сообщение
		color.WhiteString(string(b)), // JSON атрибуты белым цветом
	)

	return nil
}

// WithAttrs создает новый обработчик с дополнительными атрибутами
// Эти атрибуты будут добавлены ко всем последующим записям лога
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler, // Сохраняем базовый обработчик
		l:       h.l,       // Сохраняем логгер
		attrs:   attrs,     // ВНИМАНИЕ: здесь потенциальная проблема - перезаписываем attrs
	}
}

// WithGroup создает новый обработчик с группировкой атрибутов
// TODO: Полная реализация группировки еще не готова
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// TODO: implement - нужно добавить поддержку группировки в pretty-формате
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name), // Делегируем группировку базовому обработчику
		l:       h.l,                       // Сохраняем логгер
	}
}
