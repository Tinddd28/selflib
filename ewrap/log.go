package ewrap

import "github.com/Tinddd28/selflib/types"

type EwrapLogger func(msg string, err error, fs ...types.Field)

func Log(err error, f EwrapLogger) {
	if err == nil {
		return
	}

	if e, ok := err.(*E); ok {
		var wrapped error
		if len(e.errs) > 1 {
			wrapped = e.errs[1]
		}

		f(e.Reason(), wrapped, e.fields...)

		return
	}

	f(err.Error(), nil)
}
