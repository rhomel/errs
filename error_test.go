package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rhomel/errs"
)

func Test_NewError_Is(t *testing.T) {
	t.Run("outer error wraps inner error", func(t *testing.T) {
		errInner := fmt.Errorf("inner error")
		errOuter := fmt.Errorf("outer error")

		combined := errs.NewError(errOuter, errInner)

		assertb(t, "error should be errInner",
			true,
			errors.Is(combined, errInner))
		assertb(t, "error should be errOuter",
			true,
			errors.Is(combined, errOuter))
		asserts(t, "error message should match",
			"outer error: inner error",
			fmt.Sprintf("%v", combined))
	})
	t.Run("triple layered error", func(t *testing.T) {
		err1 := fmt.Errorf("first")
		err2 := fmt.Errorf("second")
		err3 := fmt.Errorf("third")

		combined := errs.NewError(err3, errs.NewError(err2, err1))

		assertb(t, "error should be err1",
			true,
			errors.Is(combined, err1))
		assertb(t, "error should be err2",
			true,
			errors.Is(combined, err2))
		assertb(t, "error should be err3",
			true,
			errors.Is(combined, err3))
		asserts(t, "error message should match",
			"third: second: first",
			fmt.Sprintf("%v", combined))
	})
	t.Run("constant error wraps a source error", func(t *testing.T) {
		sourceErr := fmt.Errorf("source error")
		constantErr := errs.Const("const error")

		combined := errs.NewError(constantErr, sourceErr)

		assertb(t, "error should be constant error",
			true,
			errors.Is(combined, constantErr))
		assertb(t, "error should be source error",
			true,
			errors.Is(combined, sourceErr))
		asserts(t, "error message should match",
			"const error: source error",
			fmt.Sprintf("%v", combined))
	})
}

type customError struct {
	value string
}

func (e *customError) Error() string {
	return e.value
}

func Test_NewError_As(t *testing.T) {
	t.Run("can unwrap inner error", func(t *testing.T) {
		sourceError := customError{"source error"}
		wrappingError := fmt.Errorf("wrapping error")

		combined := errs.NewError(wrappingError, &sourceError)
		var custom *customError
		ok := errors.As(combined, &custom)

		assertb(t, "errors.As should be ok",
			true,
			ok)
		asserts(t, "unwrapped error message should match",
			sourceError.value,
			custom.value)
	})
	t.Run("can unwrap outer error", func(t *testing.T) {
		sourceError := fmt.Errorf("source error")
		wrappingError := customError{"wrapping error"}

		combined := errs.NewError(&wrappingError, sourceError)
		var custom *customError
		ok := errors.As(combined, &custom)

		assertb(t, "errors.As should be ok",
			true,
			ok)
		asserts(t, "unwrapped error message should match",
			wrappingError.value,
			custom.value)
	})
	t.Run("does not unwrap error", func(t *testing.T) {
		sourceError := errs.Const("source")
		wrappingError := errs.Const("wrapper")

		combined := errs.NewError(wrappingError, sourceError)
		var custom *customError
		ok := errors.As(combined, &custom)

		assertb(t, "errors.As should not be ok",
			false,
			ok)
		if custom != nil {
			t.Errorf("custom should be nil, got: %v", custom)
		}
	})
	t.Run("unwraps Error type", func(t *testing.T) {
		sourceError := errs.Const("source")
		wrappingError := errs.Const("wrapper")

		combined := errs.NewError(wrappingError, sourceError)
		var err *errs.Error
		ok := errors.As(combined, &err)

		assertb(t, "errors.As should be ok",
			true,
			ok)
		asserts(t, "unwrapped Error.Current matches",
			combined.Current.Error(),
			err.Current.Error())
		asserts(t, "unwrapped Error.Next matches",
			combined.Next.Error(),
			err.Next.Error())
	})
}

func assertb(t *testing.T, label string, want, got bool) {
	t.Helper()
	if want != got {
		fail(t, label, want, got)
	}
}

// TODO: maybe generics?
func asserts(t *testing.T, label string, want, got string) {
	t.Helper()
	if want != got {
		fail(t, label, want, got)
	}
}

func fail(t *testing.T, label string, want, got interface{}) {
	t.Helper()
	t.Errorf("%s:\nWant: %v\n Got: %v", label, want, got)
}
