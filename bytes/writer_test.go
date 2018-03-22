package bytes

import (
	"testing"
)

func TestWriter(t *testing.T) {
	w := NewWriterSize(10)
	w.WriteString("my writer test")
	w.WriteString(" new+")

	t.Logf("print:%s\n", w.Buffer())
}
