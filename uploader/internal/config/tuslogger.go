package config

import (
	"context"
	"log/slog"

	"github.com/phuslu/log"
	expslog "golang.org/x/exp/slog"
)

type LogHandler struct {
	slog.Handler
}

func (l *LogHandler) Enabled(ctx context.Context, level expslog.Level) bool {
	return true
}

func (l *LogHandler) Handle(ctx context.Context, r expslog.Record) error {
	return l.Handler.Handle(ctx, slog.Record{
		Level:   slog.Level(r.Level),
		Time:    r.Time,
		Message: r.Message,
		PC:      r.PC,
	})
}

func (l *LogHandler) WithAttrs(attrs []expslog.Attr) expslog.Handler {
	var attrs2 []slog.Attr
	for _, a := range attrs {
		attrs2 = append(attrs2, slog.Attr{
			Key:   a.Key,
			Value: convertExpSlogValueToSlogValue(a.Value), // Correct conversion
		})
	}
	return &LogHandler{
		Handler: l.Handler.WithAttrs(attrs2),
	}
}

func (l *LogHandler) WithGroup(name string) expslog.Handler {
	return &LogHandler{
		Handler: l.Handler.WithGroup(name),
	}
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		Handler: log.DefaultLogger.Slog().Handler(),
	}
}

// Helper function to convert expslog.Value to slog.Value
func convertExpSlogValueToSlogValue(v expslog.Value) slog.Value {
	switch v.Kind() {
	case expslog.KindString:
		return slog.StringValue(v.String())
	case expslog.KindInt64:
		return slog.Int64Value(v.Int64())
	case expslog.KindFloat64:
		return slog.Float64Value(v.Float64())
	case expslog.KindBool:
		return slog.BoolValue(v.Bool())
	case expslog.KindTime:
		return slog.TimeValue(v.Time())
	case expslog.KindDuration:
		return slog.DurationValue(v.Duration())
	case expslog.KindGroup:
		attrs := v.Group()
		var slogAttrs []slog.Attr
		for _, attr := range attrs {
			slogAttrs = append(slogAttrs, slog.Attr{
				Key:   attr.Key,
				Value: convertExpSlogValueToSlogValue(attr.Value),
			})
		}
		return slog.GroupValue(slogAttrs...)
	default:
		return slog.AnyValue(v.Any()) // Fallback for unknown types
	}
}
