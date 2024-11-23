package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	systemAttr = "sys"
)

var (
	output io.Writer = os.Stderr
	writer io.Writer
	closer io.Closer
	deflog = New("")
)

func init() {
	slog.SetDefault(deflog.Logger)
}

type global struct{}

func (global) Write(p []byte) (int, error) {
	if output == nil {
		return len(p), nil
	}
	n, err := output.Write(p)
	if writer != nil {
		n, err = writer.Write(p)
	}
	return n, err
}

type handler struct {
	w     io.Writer
	sys   string
	attrs []slog.Attr
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return true // TODO
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}
	n := len(h.attrs)
	attrs := h.attrs[:n:n]
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	sys := h.sys
	var abuf bytes.Buffer
	for _, a := range attrs {
		switch a.Key {
		case systemAttr:
			sys = a.Value.String()
		default:
			abuf.WriteString(a.Key)
			abuf.WriteString("=")
			abuf.WriteString(a.Value.String())
			abuf.WriteString(" ")
		}
	}
	if sys != "" {
		sys = "[" + sys + "] "
	}
	pref := "INFO"
	switch r.Level {
	case slog.LevelDebug:
		pref = "DEBUG"
	case slog.LevelInfo:
		pref = "INFO"
	case slog.LevelWarn:
		pref = "WARN"
	case slog.LevelError:
		pref = "ERROR"
	}
	_, err := fmt.Fprintf(h.w, "%s %s %s%s %s\n",
		r.Time.Format(time.DateTime),
		pref, sys,
		strings.TrimSpace(r.Message),
		strings.TrimSpace(abuf.String()),
	)
	return err
}

func (h *handler) clone() *handler {
	n := len(h.attrs)
	return &handler{
		w:     h.w,
		sys:   h.sys,
		attrs: h.attrs[:n:n],
	}
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()
	for _, a := range attrs {
		switch a.Key {
		case systemAttr:
			h2.sys = a.Value.String()
		default:
			h2.attrs = append(h2.attrs, a)
		}
	}
	return h2
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h // TODO
}

type Logger struct {
	*slog.Logger
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Print(args ...interface{}) {
	l.Info(fmt.Sprintln(args...))
}

func (l *Logger) Println(args ...interface{}) {
	l.Info(fmt.Sprintln(args...))
}

func NewHandler() slog.Handler {
	return &handler{w: global{}}
}

func New(name string) *Logger {
	if name != "" {
		name = "main"
	}
	log := slog.New(NewHandler())
	return &Logger{log}
}

func WithSystem(log *slog.Logger, name string) *slog.Logger {
	return log.With(systemAttr, name)
}

func Printf(format string, args ...interface{}) {
	deflog.Printf(format, args...)
}

func Println(args ...interface{}) {
	deflog.Println(args...)
}

func WriteToFile(path string) error {
	Close()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	writer = f
	closer = f
	return nil
}

func SetOutput(w io.Writer) {
	output = w
}

func Close() {
	if closer != nil {
		_ = closer.Close()
		closer = nil
	}
}
