package core

import (
	"strings"
	"testing"
)

func TestAnswerLengthIgnoresLineBreaks(t *testing.T) {
	if got := AnswerLength("あい\r\nう\nx"); got != 4 {
		t.Fatalf("length = %d, want 4", got)
	}
}

func TestValidateAnswerLimits(t *testing.T) {
	problem := Problem{AnswerMaxLength: 5, AllowAttachments: true}
	if err := ValidateAnswer(problem, "12345\n\n", nil); err != nil {
		t.Fatalf("line breaks counted: %v", err)
	}
	if err := ValidateAnswer(problem, "123456", nil); err == nil {
		t.Fatal("over-limit answer accepted")
	}
	if err := ValidateAnswer(Problem{}, strings.Repeat("a", DefaultAnswerMaxLength), nil); err != nil {
		t.Fatalf("default limit: %v", err)
	}
	valid := AnswerAttachment{Name: "syslog.txt", Content: "log"}
	for name, attachments := range map[string][]AnswerAttachment{
		"path":      {{Name: "../etc/passwd.txt", Content: "x"}},
		"hidden":    {{Name: ".bashrc.txt", Content: "x"}},
		"extension": {{Name: "run.sh", Content: "x"}},
		"html":      {{Name: "page.html", Content: "x"}},
		"nul":       {{Name: "bin.log", Content: "a\x00b"}},
		"invalid":   {{Name: "bin.log", Content: "\xff\xfe"}},
		"empty":     {{Name: "empty.log", Content: ""}},
		"too large": {{Name: "big.log", Content: strings.Repeat("a", MaxAnswerAttachmentSize+1)}},
		"duplicate": {valid, valid},
		"too many":  {valid, {Name: "b.txt", Content: "x"}, {Name: "c.txt", Content: "x"}, {Name: "d.txt", Content: "x"}},
	} {
		t.Run(name, func(t *testing.T) {
			if err := ValidateAnswer(problem, "ok", attachments); err == nil {
				t.Fatal("invalid attachment accepted")
			}
		})
	}
	if err := ValidateAnswer(Problem{}, "ok", []AnswerAttachment{valid}); err == nil {
		t.Fatal("attachment accepted for a problem that does not allow them")
	}
	if err := ValidateAnswer(problem, "ok", []AnswerAttachment{valid, {Name: "app.log", Content: "\x1b[31merror\x1b[0m"}}); err != nil {
		t.Fatalf("valid attachments: %v", err)
	}
}
