package validator

import "regexp"

const (
	FailLenMin = 2 << iota
	FailLenMax
	FailValMin
	FailValMax
	FailRegExp
	FailEmail
	FailReq
	FailType
)

// pre‑compiled e‑mail regexp (RFC‑5322‑ish, good enough for most cases)
var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
