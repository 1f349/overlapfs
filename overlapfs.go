package overlapfs

import (
	"errors"
	"fmt"
	"io/fs"
)

// OverlapFS layers B on top of A
//
// Open and Stat reimplemented to show B unless ErrNotExist is returned then show A.
//
// Glob and ReadDir are reimplemented to uniquely merge the returned lists.
//
// For ReadFile and Sub use the standard fs.ReadFile and fs.Sub implementations
type OverlapFS struct {
	A, B fs.FS
}

var _ interface {
	fs.FS
	fs.StatFS
	fs.GlobFS
	fs.ReadDirFS
} = &OverlapFS{}

func (o OverlapFS) Open(name string) (fs.File, error) {
	fmt.Println("Open", name)
	open, err := o.B.Open(name)
	switch {
	case err == nil:
		return open, nil
	case errors.Is(err, fs.ErrNotExist):
		return o.A.Open(name)
	}
	return nil, err
}

func (o OverlapFS) Stat(name string) (fs.FileInfo, error) {
	stat, err := fs.Stat(o.B, name)
	switch {
	case err == nil:
		return stat, nil
	case errors.Is(err, fs.ErrNotExist):
		return fs.Stat(o.A, name)
	}
	return nil, err
}

func (o OverlapFS) Glob(pattern string) ([]string, error) {
	aSlice, err := fs.Glob(o.B, pattern)
	if err != nil {
		return nil, err
	}
	bSlice, err := fs.Glob(o.A, pattern)
	if err != nil {
		return nil, err
	}

	mergeUnique(aSlice, bSlice, func(s string) string { return s })
	return aSlice, nil
}

func (o OverlapFS) ReadDir(name string) ([]fs.DirEntry, error) {
	aSlice, err := fs.ReadDir(o.B, name)
	if err != nil {
		return nil, err
	}
	bSlice, err := fs.ReadDir(o.A, name)
	if err != nil {
		return nil, err
	}

	mergeUnique(aSlice, bSlice, func(entry fs.DirEntry) string { return entry.Name() })
	return aSlice, nil
}

// mergeUnique puts unique values from b into a
func mergeUnique[T any](a, b []T, key func(T) string) {
	m := make(map[string]struct{})
	for _, i := range a {
		m[key(i)] = struct{}{}
	}
	for _, i := range b {
		if _, ok := m[key(i)]; !ok {
			a = append(a, i)
		}
	}
}
