package ch13

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestSuccessfullyRead(t *testing.T) {
	properties := gopter.NewProperties(nil)

	skipBytesGen := gen.AnyString()
	readBytesGen := gen.AnyString()
	remainedBytesGen := gen.AnyString()

	properties.Property("Read successfully with SeekStart", prop.ForAll(func(skipBytesSrc, readBytesSrc, remainedBytesSrc string) bool {
		skipBytes := []byte(skipBytesSrc)
		readBytes := []byte(readBytesSrc)
		remainedBytes := []byte(remainedBytesSrc)
		var b bytes.Buffer
		b.Write(skipBytes)
		b.Write(readBytes)
		b.Write(remainedBytes)

		readSeeker := NewSeeker(&b)
		_, err := readSeeker.Seek(int64(len(skipBytes)), io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		var output bytes.Buffer
		n, err := io.CopyN(&output, readSeeker, int64(len(readBytes)))

		if n != int64(len(readBytes)) {
			t.Logf("Read successfully with SeekStart fail: n(%d) != len(readBytes)()%d\n", n, len(readBytes))
			return false
		}

		if output.String() != readBytesSrc {
			t.Logf("Read successfully with SeekStart fail: %s != %s\n", output.String(), readBytesSrc)
			return false
		}

		if err != nil {
			t.Logf("Read successfully with SeekStart fail: err(%v)\n", err)
			return false
		}

		return true
	}, skipBytesGen, readBytesGen, remainedBytesGen))
	properties.TestingRun(t)
}
