//go:build go1.18
// +build go1.18

package stream

import (
	"bytes"
	"testing"
)

func Fuzz_simpleReplacer_Read(f *testing.F) {
	type args struct {
		content []byte
		old     []byte
		new     []byte
	}
	tests := []args{
		{
			[]byte(`aaaaaaaaaaaaaaoldbbbbbbbbbbboldbbbbbbb`),
			[]byte(`old`),
			[]byte(`new`),
		},
		{
			[]byte(`111111aaaaaaaaaaaaaaoldbbbbbbbbbbboldbbbbbbol`),
			[]byte(`old`),
			[]byte(`new`),
		},
		{
			[]byte(`aaaaaaaaaa`),
			[]byte(`a`),
			bytes.Repeat([]byte(`long`), 100),
		},
		{
			[]byte(`aaaaaaaaaa`),
			[]byte(``),
			[]byte(`long long long long long long long long long string`),
		},
		{
			[]byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			[]byte("00"),
			[]byte("0"),
		},
	}
	for _, tt := range tests {
		f.Add(tt.content, tt.old, tt.new)
	}

	f.Fuzz(func(t *testing.T, content, old, new []byte) {
		runSimpleReplacerTest(t, content, old, new, -1)
	})
}
