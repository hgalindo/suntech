package st600

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/larixsource/suntech/st"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAllST600Spec(t *testing.T) {
	f, err := os.Open("ascii_spec.txt")
	require.Nil(t, err)
	defer f.Close()

	// save each example frame from the spec
	var specFrames [][]byte
	// buffer with all content to parse
	var buf bytes.Buffer
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// skip empty lines
		if len(scanner.Bytes()) > 0 && scanner.Bytes()[0] != '#' {
			b := scanner.Bytes()
			b = bytes.TrimPrefix(b, []byte("[command] "))
			b = bytes.TrimPrefix(b, []byte("[response] "))
			b = bytes.TrimPrefix(b, []byte("[report] "))
			frame := make([]byte, 0, len(b)+1)
			frame = append(frame, b...)
			frame = append(frame, st.EndOfFrame)
			specFrames = append(specFrames, frame)
			buf.Write(frame)
		}
	}
	require.Nil(t, scanner.Err())

	i := 0
	p := ParseBytes(buf.Bytes(), ParserOpts{
		SkipUnknownFrames: true,
	})
	// parse and compare raw frame with the spec frames loaded before
	for p.Next() {
		assert.Equal(t, specFrames[i], p.Msg().Frame, "not equals:\n%s\n%s", specFrames[i], p.Msg().Frame)
		i++
	}
	require.Nil(t, p.Error())
}
