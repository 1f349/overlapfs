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
	assert.Equal(t, "Hello World!", string(bAll))
}

func TestOverlapFS_Stat(t *testing.T) {
	o := OverlapFS{A: a.Embed, B: b.Embed}
	open, err := o.Stat("hello.txt")
	assert.NoError(t, err)
	stat, err := open.Stat()
	assert.NoError(t, err)
	assert.Equal(t, "hello.txt", stat.Name())
	bAll, err := io.ReadAll(open)
	assert.Equal(t, "Hello World!", string(bAll))
}
