package log

import (
	"context"
	"io"
	"log/slog"
	"strconv"
)

func InitLog(w io.Writer) {
	slog.SetDefault(slog.New(NewDefaultContextParsedJSONHandler(w)))
}

var (
	TraceHeader        = "exermon-trace"
	DefaultTraceParser = func(ctx context.Context) (string, string) {
		traceId := ctx.Value(TraceHeader)
		if traceId == nil {
			return "log_id", ""
		}
		return "log_id", traceId.(string)
	}
)

type ContextParser func(ctx context.Context) (string, string)

type ContextParsedJSONHandler struct {
	definedCtxHandlers []ContextParser
	jsonHandler        *slog.JSONHandler
}

func NewDefaultContextParsedJSONHandler(w io.Writer) *ContextParsedJSONHandler {
	return NewContextParsedJSONHandler(w, []ContextParser{DefaultTraceParser})
}

func NewContextParsedJSONHandler(w io.Writer, definedCtxHandlers []ContextParser) *ContextParsedJSONHandler {
	return &ContextParsedJSONHandler{
		definedCtxHandlers: definedCtxHandlers,
		jsonHandler: slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: true,
			Level:     nil,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.TimeKey:
					t := a.Value.Time()
					return slog.Attr{
						Key:   slog.TimeKey,
						Value: slog.StringValue(t.Format("2006-01-02T15:04:05.999")),
					}
				case slog.SourceKey:
					source := a.Value.Any().(*slog.Source)
					return slog.Attr{
						Key:   slog.SourceKey,
						Value: slog.StringValue(source.File + ":" + strconv.FormatInt(int64(source.Line), 10)),
					}
				}
				return a
			},
		}),
	}
}

func (c ContextParsedJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.jsonHandler.Enabled(ctx, level)
}

func (c ContextParsedJSONHandler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(c.handleCtx(ctx)...)
	return c.jsonHandler.Handle(ctx, record)
}

func (c ContextParsedJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c.jsonHandler.WithAttrs(attrs)
}

func (c ContextParsedJSONHandler) WithGroup(name string) slog.Handler {
	return c.jsonHandler.WithGroup(name)
}

func (c ContextParsedJSONHandler) handleCtx(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	for _, hand := range c.definedCtxHandlers {
		key, value := hand(ctx)
		if value == "" {
			continue
		}
		attrs = append(attrs, slog.String(key, value))
	}
	return attrs
}

func (c ContextParsedJSONHandler) AddContextParse(parse ContextParser) {
	c.definedCtxHandlers = append(c.definedCtxHandlers, parse)
}
