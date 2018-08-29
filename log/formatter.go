package astroflow

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/astroflow/astroflow-go"
	"github.com/json-iterator/go"
)

const (
	cReset    = 0
	cBold     = 1
	cRed      = 31
	cGreen    = 32
	cYellow   = 33
	cBlue     = 34
	cMagenta  = 35
	cCyan     = 36
	cGray     = 37
	cDarkGray = 90
)

type CLIFormatter struct {
	NoColor            bool
	TimestampFieldName string
	LevelFieldName     string
	MessageFieldName   string
}

func NewCLIFormatter() CLIFormatter {
	return CLIFormatter{
		TimestampFieldName: astroflow.TimestampFieldName,
		MessageFieldName:   astroflow.MessageFieldName,
		LevelFieldName:     astroflow.LevelFieldName,
		NoColor:            false,
	}
}

func (formatter CLIFormatter) Format(event astroflow.Event) []byte {
	var ret = new(bytes.Buffer)

	if t, ok := event[formatter.TimestampFieldName].(time.Time); ok {
		ret.WriteString(t.Format(time.RFC3339))
	}

	lvlColor := cReset
	level := ""
	if l, ok := event[formatter.LevelFieldName].(string); ok {
		level = l
		if !formatter.NoColor {
			lvlColor = levelColor(l)
			ret.WriteString(colorize(levelSymbol(l), levelColor(l), !formatter.NoColor))
		}
	}

	if m, ok := event[formatter.MessageFieldName].(string); ok {
		ret.WriteString(m)
	}

	// do not display additional fields when level == info
	if level == "info" {
		return ret.Bytes()
	}

	fields := make([]string, 0, len(event))
	for field := range event {
		switch field {
		case formatter.TimestampFieldName, formatter.MessageFieldName, formatter.LevelFieldName:
			continue
		}

		fields = append(fields, field)
	}

	sort.Strings(fields)
	for _, field := range fields {
		if needsQuote(field) {
			field = strconv.Quote(field)
		}
		fmt.Fprintf(ret, " %s=", colorize(field, lvlColor, !formatter.NoColor))

		switch value := event[field].(type) {
		case string:
			if len(value) == 0 {
				ret.WriteString("\"\"")
			} else if needsQuote(value) {
				ret.WriteString(strconv.Quote(value))
			} else {
				ret.WriteString(value)
			}
		case time.Time:
			ret.WriteString(value.Format(time.RFC3339))
		default:
			b, err := jsoniter.Marshal(value)
			if err != nil {
				fmt.Fprintf(ret, "[error: %v]", err)
			} else {
				fmt.Fprint(ret, string(b))
			}
		}

	}

	return ret.Bytes()
}

func levelColor(level string) int {
	switch level {
	case "debug":
		return cMagenta
	case "info":
		return cGreen
	case "warning":
		return cYellow
	case "error", "fatal":
		return cRed
	default:
		return cReset
	}
}

func levelSymbol(level string) string {
	switch level {
	case "info":
		return " ✔ "
	case "warning":
		return " ⚠ "
	case "error", "fatal":
		return " ✘ "
	default:
		return " • "
	}
}

func colorize(s interface{}, color int, enabled bool) string {
	if !enabled {
		return fmt.Sprintf("%v", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, s)
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}
