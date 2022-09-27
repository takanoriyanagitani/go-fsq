package fsq

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

type NextQueue func(ctx context.Context, previous string) (next string, err error)

type NextCheck func(next string) (checked string, err error)

type opener func(name string) (fs.File, error)

func (o opener) Open(name string) (fs.File, error) { return o(name) }

var openerDefault opener = ComposeErr(
	os.Open, // string -> *os.File, error
	func(f *os.File) (fs.File, error) { return f, nil },
)

func openerNew(f fs.FS) (opener, error) {
	return ErrFromBool(
		nil != f,
		func() opener { return f.Open },
		func() error { return fmt.Errorf("Invalid filesystem") },
	)
}

func openerNewOr(alt opener, f fs.FS) opener {
	return ErrUnwrapOrElse(
		openerNew,
		func(_ error) opener { return alt },
	)(f)
}

var openerNewOrDefault func(fs.FS) opener = Curry(openerNewOr)(openerDefault)

var fsDefault fs.FS = openerNewOrDefault(nil)

// NextCheckBuilder creates NextCheck.
// If "next" file exists, it returns fs.ErrExist
func NextCheckBuilder(f fs.FS) NextCheck {
	var _f fs.FS = openerNewOrDefault(f)
	return func(next string) (checked string, err error) {
		file, e := _f.Open(next)
		if nil == e {
			_ = file.Close()
			return "", fs.ErrExist
		}
		return ErrFromBool(
			errors.Is(e, fs.ErrNotExist),
			func() string { return next },
			func() error { return fmt.Errorf("Unexpected error: %v", e) },
		)
	}
}

func (n NextQueue) UnwrapOrElse(ctx context.Context, prev string, alt func() string) (next string) {
	return ErrUnwrapOrElse(
		func(p string) (string, error) { return n(ctx, p) },
		func(_ error) string { return alt() },
	)(prev)
}

func (n NextQueue) ToChecked(checker NextCheck) NextQueue {
	return ComposeContext(
		func(ctx context.Context, prev string) (next string, err error) { return n(ctx, prev) },
		func(ctx context.Context, next string) (string, error) { return checker(next) },
	)
}

func (n NextQueue) WithoutDir() NextQueue {
	return func(ctx context.Context, previous string) (next string, err error) {
		var basename string = filepath.Base(previous)
		var dirname string = filepath.Dir(previous)
		next, e := n(ctx, basename)
		return filepath.Join(dirname, next), e
	}
}

var NextQueueI64 NextQueue = ComposeContext(
	func(_ context.Context, prev string) (int64, error) { return strconv.ParseInt(prev, 16, 64) },
	func(_ context.Context, p int64) (string, error) { return fmt.Sprintf("%016x", p+1), nil },
)
