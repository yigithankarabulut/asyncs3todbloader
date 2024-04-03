package customerror_test

import (
	"errors"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/customerror"
	"testing"
)

func TestError(t *testing.T) {
	err := customerror.New("some error", true)
	var customErr *customerror.Error

	if !errors.As(err, &customErr) {
		t.Errorf("error does not match the target type, want: %T, got: %v", customErr, err)
	}
	shouldEqual := "some error"
	if customErr.Error() != shouldEqual {
		t.Errorf("error message does not match, want: %s, got: %s", shouldEqual, customErr.Error())
	}
	shouldLoggable := true
	if customErr.Loggable != shouldLoggable {
		t.Errorf("error should be loggable, want: %t, got: %t", shouldLoggable, customErr.Loggable)
	}
}

func TestErrorWrap(t *testing.T) {
	err := customerror.New("some error", true)
	err = err.Wrap(errors.New("wrapped error"))
	var customErr *customerror.Error

	if !errors.As(err, &customErr) {
		t.Errorf("error does not match the target type, want: %T, got: %v", customErr, err)
	}
	shouldEqual := "wrapped error, some error"
	if customErr.Error() != shouldEqual {
		t.Errorf("error message does not match, want: %s, got: %s", shouldEqual, customErr.Error())
	}
	shouldLoggable := true
	if customErr.Loggable != shouldLoggable {
		t.Errorf("error should be loggable, want: %t, got: %t", shouldLoggable, customErr.Loggable)
	}

	if customErr.Err == nil {
		t.Errorf("wrapped error can not be nil, want: %v, got: nil", customErr.Err)
	}
}

func TestUnwrap(t *testing.T) {
	err := customerror.New("some error", false)
	wrappedErr := err.Wrap(errors.New("inner")) // nolint

	var customErr *customerror.Error

	if !errors.As(wrappedErr, &customErr) {
		t.Errorf("error does not match the target type, want: %T, got: %v", customErr, err)
	}

	shouldEqual := "inner"
	unwrappedErr := customErr.Unwrap()
	if unwrappedErr.Error() != shouldEqual {
		t.Errorf("unwrapped error does not match, want: %s, got: %s", shouldEqual, unwrappedErr.Error())
	}
}

func TestAddDataDestroyData(t *testing.T) {
	err := customerror.New("some error", false).AddData("hello")

	var customErr *customerror.Error

	if !errors.As(err, &customErr) {
		t.Errorf("error does not match the target type, want: %T, got: %v", customErr, err)
	}

	if customErr.Data == nil {
		t.Errorf("data should not be nil, want: %v, got: nil", customErr.Data)
	}

	shouldEqual := "hello"
	data, ok := customErr.Data.(string)
	if !ok {
		t.Error("data should be assertable to string")
	}

	if data != shouldEqual {
		t.Errorf("data does not match, want: %s, got: %s", shouldEqual, data)
	}

	shouldEqual = "some error"
	if err.Error() != shouldEqual {
		t.Errorf("error does not match, want: %s, got: %s", shouldEqual, err.Error())
	}

	err = err.DestroyData()
	if !errors.As(err, &customErr) {
		t.Errorf("error does not match the target type, want: %T, got: %v", customErr, err)
	}

	if customErr.Data != nil {
		t.Errorf("data should be nil, want: nil, got: %v", customErr.Data)
	}
}
