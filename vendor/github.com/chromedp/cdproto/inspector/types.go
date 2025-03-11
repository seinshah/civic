package inspector

import (
	"fmt"
	"strings"
)

// Code generated by cdproto-gen. DO NOT EDIT.

// DetachReason detach reason.
//
// See: ( -- none -- )
type DetachReason string

// String returns the DetachReason as string value.
func (t DetachReason) String() string {
	return string(t)
}

// DetachReason values.
const (
	DetachReasonTargetClosed         DetachReason = "target_closed"
	DetachReasonCanceledByUser       DetachReason = "canceled_by_user"
	DetachReasonReplacedWithDevtools DetachReason = "replaced_with_devtools"
	DetachReasonRenderProcessGone    DetachReason = "Render process gone."
)

// UnmarshalJSON satisfies [json.Unmarshaler].
func (t *DetachReason) UnmarshalJSON(buf []byte) error {
	s := string(buf)
	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)

	switch DetachReason(s) {
	case DetachReasonTargetClosed:
		*t = DetachReasonTargetClosed
	case DetachReasonCanceledByUser:
		*t = DetachReasonCanceledByUser
	case DetachReasonReplacedWithDevtools:
		*t = DetachReasonReplacedWithDevtools
	case DetachReasonRenderProcessGone:
		*t = DetachReasonRenderProcessGone
	default:
		return fmt.Errorf("unknown DetachReason value: %v", s)
	}
	return nil
}
