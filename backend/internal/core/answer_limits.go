package core

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	// DefaultAnswerMaxLength applies to problems that do not set
	// answer_max_length.
	DefaultAnswerMaxLength = 3000
	// MaxAnswerLength is the largest answer_max_length a problem may set.
	MaxAnswerLength = 10000
	// MaxAnswerBodyLength bounds the body with line breaks, which the
	// per-problem limit does not count; it matches the OpenAPI maxLength.
	MaxAnswerBodyLength = 30000

	MaxAnswerAttachments    = 3
	MaxAnswerAttachmentSize = 1 << 20
)

// attachmentNamePattern keeps names usable as plain download file names:
// no paths, no leading dot, and only .txt or .log.
var attachmentNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._-]{0,95}\.(txt|log)$`)

func (p Problem) EffectiveAnswerMaxLength() int32 {
	if p.AnswerMaxLength <= 0 {
		return DefaultAnswerMaxLength
	}
	return p.AnswerMaxLength
}

// AnswerLength counts characters (Unicode code points) without line breaks.
func AnswerLength(body string) int {
	return utf8.RuneCountInString(body) - strings.Count(body, "\n") - strings.Count(body, "\r")
}

// ValidateAnswer checks an answer body and its attachments against the
// problem's limits before it is stored.
func ValidateAnswer(problem Problem, body string, attachments []AnswerAttachment) error {
	if utf8.RuneCountInString(body) > MaxAnswerBodyLength {
		return NewError(http.StatusUnprocessableEntity, "validation_error",
			fmt.Sprintf("Answer is over %d characters including line breaks", MaxAnswerBodyLength))
	}
	if length, limit := AnswerLength(body), int(problem.EffectiveAnswerMaxLength()); length > limit {
		return NewError(http.StatusUnprocessableEntity, "validation_error",
			fmt.Sprintf("Answer is %d characters, over the limit of %d (line breaks are not counted)", length, limit))
	}
	if len(attachments) == 0 {
		return nil
	}
	if !problem.AllowAttachments {
		return NewError(http.StatusUnprocessableEntity, "validation_error", "This problem does not accept attachments")
	}
	if len(attachments) > MaxAnswerAttachments {
		return NewError(http.StatusUnprocessableEntity, "validation_error", fmt.Sprintf("At most %d attachments are allowed", MaxAnswerAttachments))
	}
	names := make(map[string]struct{}, len(attachments))
	for _, attachment := range attachments {
		if !attachmentNamePattern.MatchString(attachment.Name) {
			return NewError(http.StatusUnprocessableEntity, "validation_error",
				"Attachment names must be letters, digits, '.', '_' or '-' and end with .txt or .log")
		}
		if _, duplicate := names[attachment.Name]; duplicate {
			return NewError(http.StatusUnprocessableEntity, "validation_error", "Attachment names must be unique")
		}
		names[attachment.Name] = struct{}{}
		if len(attachment.Content) == 0 || len(attachment.Content) > MaxAnswerAttachmentSize {
			return NewError(http.StatusUnprocessableEntity, "validation_error", "Attachments must be between 1 byte and 1 MiB")
		}
		if !utf8.ValidString(attachment.Content) || strings.ContainsRune(attachment.Content, 0) {
			return NewError(http.StatusUnprocessableEntity, "validation_error", "Attachments must be UTF-8 text")
		}
	}
	return nil
}
