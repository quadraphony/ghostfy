package logging

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

func New(out io.Writer) *Logger {
	return &Logger{out: out}
}

func (l *Logger) Log(level, msg string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var b strings.Builder
	b.WriteString("time=")
	b.WriteString(time.Now().UTC().Format(time.RFC3339))
	b.WriteString(" level=")
	b.WriteString(level)
	b.WriteString(" msg=")
	b.WriteString(strconvQuote(msg))

	for key, value := range fields {
		b.WriteByte(' ')
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(strconvQuote(fmt.Sprint(value)))
	}

	b.WriteByte('\n')
	_, _ = io.WriteString(l.out, b.String())
}

func strconvQuote(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	return `"` + replacer.Replace(s) + `"`
}
