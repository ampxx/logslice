package parser

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTimestamp_KnownLayouts(t *testing.T) {
	cases := []struct {
		input  string
		wantOK bool
	}{
		{"2024-03-15T12:30:00Z", true},
		{"2024-03-15 12:30:00", true},
		{"2024/03/15 12:30:00", true},
		{"not-a-timestamp", false},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			ts, _, err := ParseTimestamp(tc.input)
			if tc.wantOK {
				require.NoError(t, err)
				assert.False(t, ts.IsZero())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestLineParser_Parse(t *testing.T) {
	input := `2024-03-15T08:00:01Z INFO  server started
2024-03-15T08:00:05Z ERROR failed to connect to db
2024-03-15T08:01:00Z DEBUG retrying connection
`

	p := NewLineParser(strings.NewReader(input))
	entries, err := p.Parse()
	require.NoError(t, err)
	assert.Len(t, entries, 3)

	assert.Equal(t, "INFO", entries[0].Level)
	assert.Equal(t, "ERROR", entries[1].Level)
	assert.Equal(t, "DEBUG", entries[2].Level)

	expected, _ := time.Parse(time.RFC3339, "2024-03-15T08:00:01Z")
	assert.True(t, entries[0].Timestamp.Equal(expected))
}

func TestLineParser_SkipsBlankLines(t *testing.T) {
	input := "\n2024-03-15T09:00:00Z INFO hello\n\n"
	p := NewLineParser(strings.NewReader(input))
	entries, err := p.Parse()
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestLineParser_NoTimestamp(t *testing.T) {
	input := "plain log line without timestamp\n"
	p := NewLineParser(strings.NewReader(input))
	entries, err := p.Parse()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.True(t, entries[0].Timestamp.IsZero(), "expected zero timestamp for unparseable line")
}
