package gosc

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func Test_readArguments(t *testing.T) {
	buf := bytes.Buffer{}
	r := bufio.NewReader(&buf)
	w := bufio.NewWriter(&buf)

	t.Run("correctFormat", func(t *testing.T) {
		_ = writeArguments(w, []any{float32(1.0), "Test", int32(2)})
		_ = w.Flush()
		res, err := readArguments(r)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if len(res) != 3 {
			t.Errorf("expected result to contain 3 elements but was got: %d", len(res))
		}
		if reflect.TypeOf(res[0]) != reflect.TypeOf(float32(0)) {
			t.Errorf("expected first argument to be float32 but got: %T", res[0])
		}
		if reflect.TypeOf(res[1]) != reflect.TypeOf("") {
			t.Errorf("expected second argument to be string but got: %T", res[0])
		}
		if reflect.TypeOf(res[2]) != reflect.TypeOf(int32(0)) {
			t.Errorf("expected second argument to be int32 but got: %T", res[0])
		}
		buf.Reset()
	})
}

func Test_readBundle(t *testing.T) {
	buf := bytes.Buffer{}
	r := bufio.NewReader(&buf)
	w := bufio.NewWriter(&buf)

	t.Run("correctBundle", func(t *testing.T) {
		_ = writePackage(&Bundle{
			Timetag: 0,
			Messages: []*Message{
				{
					Address:   "/test",
					Arguments: []any{},
				},
			},
			Bundles: nil,
			Name:    "Test",
		}, w)
		_ = w.Flush()

		res, err := readBundle(r)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if res.Name != "Test" {
			t.Errorf("expected bundle name to be \"Test\" but got: %s", res.Name)
		}
		if len(res.Messages) != 1 {
			t.Errorf("expected bundle to contain 1 message but was: %d", len(res.Messages))
		}

		buf.Reset()
	})
}

func Test_readMessage(t *testing.T) {
	buf := bytes.Buffer{}
	r := bufio.NewReader(&buf)
	w := bufio.NewWriter(&buf)

	t.Run("correctMessage", func(t *testing.T) {
		_ = writePackage(&Message{
			Address:   "/test",
			Arguments: []any{},
		}, w)
		_ = w.Flush()
		res, err := readMessage(r)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if res.Address != "/test" {
			t.Errorf("expected address to be \"/test\" but got: %s", res.Address)
		}
		buf.Reset()
	})
}

func Test_readPackage(t *testing.T) {
	buf := bytes.Buffer{}
	r := bufio.NewReader(&buf)
	w := bufio.NewWriter(&buf)

	t.Run("correctPackage", func(t *testing.T) {
		_ = writePackage(&Message{
			Address:   "/test",
			Arguments: []any{},
		}, w)
		_ = w.Flush()
		res, err := readPackage(r)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if res.GetType() != PackageTypeMessage {
			t.Errorf("expected package type to be PackageTypeMessage but got: %s", res.GetType())
		}
	})
}

func Test_readPaddedString(t *testing.T) {
	buf := bytes.Buffer{}
	r := bufio.NewReader(&buf)
	w := bufio.NewWriter(&buf)

	t.Run("withoutPadding", func(t *testing.T) {
		_ = writePaddedString(w, "abc")
		_ = w.Flush()
		res, err := readPaddedString(r)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if res != "abc" {
			t.Errorf("expected \"abc\" but got: %s", res)
		}
		buf.Reset()
	})
	t.Run("withPadding", func(t *testing.T) {
		_ = writePaddedString(w, "testing")
		_ = w.Flush()
		res, err := readPaddedString(r)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if res != "testing" {
			t.Errorf("expected \"testing\" but got: %s", res)
		}
		buf.Reset()
	})
	t.Run("wrongDelim", func(t *testing.T) {
		_, _ = w.WriteString("test")
		_, err := readPaddedString(r)
		if !errors.Is(err, io.EOF) {
			t.Errorf("expected error io.EOF but got: %v", err)
		}
		buf.Reset()
	})
}
