package overlapfs

import (
	"errors"
	"io/fs"
	"slices"
	"strings"
)

// OverlapFS layers B on top of A
//
// Open and Stat reimplemented to show B unless ErrNotExist is returned then show A
//
// Glob and ReadDir are reimplemented to uniquely merge and sort the returned slices
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

// Open implements fs.FS, the file named will be read from B if the file exists,
// otherwise the file named will be read from A
func (o OverlapFS) Open(name string) (fs.File, error) {
	open, err := o.B.Open(name)
	switch {
	case err == nil:
		return open, nil
	case errors.Is(err, fs.ErrNotExist):
		return o.A.Open(name)
	}
	return nil, err
}

// Stat implements fs.StatFS, the file named will have stats read from B if the
// file exists, otherwise the file named will have stats read from A
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

// Glob implements fs.GlobFS, the pattern will be globed in B and A and the
// result uniquely merged and sorted
func (o OverlapFS) Glob(pattern string) ([]string, error) {
	bSlice, err := fs.Glob(o.B, pattern)
	if err != nil {
		return nil, err
	}
	aSlice, err := fs.Glob(o.A, pattern)
	if err != nil {
		return nil, err
	}

	return mergeUnique(aSlice, bSlice, func(s string) string { return s }), nil
}

// ReadDir implements fs.ReadDirFS, the directory named will be read in B and A
// and the result uniquely merged and sorted
func (o OverlapFS) ReadDir(name string) ([]fs.DirEntry, error) {
	bSlice, err := fs.ReadDir(o.B, name)
	if err != nil {
		return nil, err
	}
	aSlice, err := fs.ReadDir(o.A, name)
	if err != nil {
		return nil, err
	}

	return mergeUnique(aSlice, bSlice, func(entry fs.DirEntry) string { return entry.Name() }), nil
}

// mergeUnique puts unique values from b into a and sorts the resulting merged
// slice
func mergeUnique[T any](a, b []T, key func(T) string) []T {
	m := make(map[string]struct{})
	for _, i := range a {
		m[key(i)] = struct{}{}
	}
	for _, i := range b {
		if _, ok := m[key(i)]; !ok {
			a = append(a, i)
		}
	}
	slices.SortFunc(a, func(a, b T) int {
		return strings.Compare(key(a), key(b))
	})
	return a
}
