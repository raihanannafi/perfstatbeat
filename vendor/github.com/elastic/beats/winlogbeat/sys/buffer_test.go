package sys

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ io.Writer = &ByteBuffer{}

func TestByteBuffer(t *testing.T) {
	input := "hello"
	length := len(input)
	buf := NewByteBuffer(1024)

	n, err := buf.Write([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, length, n)

	assert.Equal(t, input, string(buf.Bytes()))
	assert.Equal(t, length, len(buf.Bytes()))
	assert.Equal(t, length, buf.Len())
}

func TestByteBufferGrow(t *testing.T) {
	input := "hello"
	length := len(input)
	buf := NewByteBuffer(0)

	n, err := buf.Write([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, length, n)

	assert.Equal(t, input, string(buf.Bytes()))
	assert.Equal(t, length, len(buf.Bytes()))
	assert.Equal(t, length, buf.Len())
	assert.Equal(t, length, len(buf.buf))

	n, err = buf.Write([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, length, n)

	assert.Equal(t, input+input, string(buf.Bytes()))
	assert.Equal(t, 2*length, len(buf.Bytes()))
	assert.Equal(t, 2*length, buf.Len())
}

func BenchmarkByteBuffer(b *testing.B) {
	input := []byte("test writing this sentence to a buffer")

	b.Run("byteBuffer", func(b *testing.B) {
		buf := NewByteBuffer(1024)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			buf.Write(input)
			buf.Bytes()
			buf.Reset()
		}
	})

	b.Run("bytes.Buffer", func(b *testing.B) {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			buf.Write(input)
			buf.Bytes()
			buf.Reset()
		}
	})
}

func BenchmarkByteBufferGrow(b *testing.B) {
	b.Run("byteBuffer", func(b *testing.B) {
		buf := NewByteBuffer(0)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			buf.Write([]byte("a"))
			buf.Bytes()
		}
	})

	b.Run("bytes.Buffer", func(b *testing.B) {
		buf := bytes.NewBuffer(make([]byte, 0))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			buf.Write([]byte("a"))
			buf.Bytes()
		}
	})
}
