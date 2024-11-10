package humanslog

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/cschleiden/humanslog/humanlog/stdiosink"
	typesv1 "github.com/humanlogio/api/go/types/v1"
	"github.com/mattn/go-isatty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HumanSlog struct {
	handler slog.Handler
	enabled bool
	stdio   *stdiosink.Stdio

	attrs []slog.Attr
	group string
}

func New(h slog.Handler) *HumanSlog {
	return &HumanSlog{
		handler: h,
		enabled: isatty.IsTerminal(os.Stderr.Fd()),
		stdio:   stdiosink.NewStdio(os.Stderr, stdiosink.DefaultStdioOpts),
	}
}

// Enabled implements slog.Handler.
func (h *HumanSlog) Enabled(ctx context.Context, lvl slog.Level) bool {
	if h.enabled {
		return true
	} else {
		return h.handler.Enabled(ctx, lvl)
	}
}

// Handle implements slog.Handler.
func (h *HumanSlog) Handle(ctx context.Context, r slog.Record) error {
	if !h.enabled {
		return h.handler.Handle(ctx, r)
	}

	msg := r.Message

	attrs := mapRecordAttrs(r, h.group)

	if len(h.attrs) > 0 {
		for _, a := range h.attrs {
			attrs = append(attrs, mapAttr(a, h.group))
		}
	}

	return h.stdio.Receive(ctx, &typesv1.LogEvent{
		ParsedAt: timestamppb.New(time.Now()),
		Structured: &typesv1.StructuredLogEvent{
			Timestamp: timestamppb.New(r.Time),
			Lvl:       r.Level.String(),
			Msg:       msg,
			Kvs:       attrs,
		},
	})
}

func mapRecordAttrs(r slog.Record, qualifier string) []*typesv1.KV {
	kvs := make([]*typesv1.KV, 0)
	r.Attrs(func(a slog.Attr) bool {
		kvs = append(kvs, mapAttr(a, qualifier))

		return true
	})

	return kvs
}

func mapAttr(a slog.Attr, qualifier string) *typesv1.KV {
	key := a.Key
	if qualifier != "" {
		key = qualifier + "." + key
	}
	return &typesv1.KV{
		Key:   key,
		Value: a.Value.String(),
	}
}

// WithAttrs implements slog.Handler.
func (h *HumanSlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HumanSlog{
		handler: h.handler.WithAttrs(attrs),
		group:   h.group,
		attrs:   append(h.attrs, attrs...),
		enabled: h.enabled,
		stdio:   h.stdio,
	}
}

// WithGroup implements slog.Handler.
func (h *HumanSlog) WithGroup(name string) slog.Handler {
	return &HumanSlog{
		handler: h.handler.WithGroup(name),
		group:   name,
		attrs:   h.attrs,
		enabled: h.enabled,
		stdio:   h.stdio,
	}
}

var _ slog.Handler = (*HumanSlog)(nil)
