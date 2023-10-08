package overlapfs

import (
	"github.com/1f349/overlapfs/a"
	"github.com/1f349/overlapfs/b"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestOverlapFS_Open(t *testing.T) {
	o := OverlapFS{A: a.Embed, B: b.Embed}
	open, err := o.Open("hello.txt")
	assert.NoError(t, err)
	stat, err := open.Stat()
	assert.NoError(t, err)
	assert.Equal(t, "hello.txt", stat.Name())
	bAll, err := io.ReadAll(open)
	assert.Equal(t, "Hello World!\n", string(bAll))
}

func TestOverlapFS_Stat(t *testing.T) {
	o := OverlapFS{A: a.Embed, B: b.Embed}
	stat, err := o.Stat("hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello.txt", stat.Name())
	assert.Equal(t, int64(13), stat.Size())
	assert.False(t, stat.IsDir())
}

func TestOverlapFS_Glob(t *testing.T) {
	o := OverlapFS{A: a.Embed, B: b.Embed}
	glob, err := o.Glob("*.txt")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a-test.txt", "hello.txt", "just-a.txt", "just-b.txt"}, glob)
}

func TestOverlapFS_ReadDir(t *testing.T) {
	o := OverlapFS{A: a.Embed, B: b.Embed}
	dir, err := o.ReadDir(".")
	assert.NoError(t, err)
	z := make([]string, 0, len(dir))
	for _, i := range dir {
		z = append(z, i.Name())
	}
	assert.Equal(t, []string{"a-test.txt", "a.go", "b.go", "example.md", "hello.txt", "just-a.txt", "just-b.txt"}, z)
}
