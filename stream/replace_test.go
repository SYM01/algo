//go:build go1.16
// +build go1.16

package stream

import (
	"bytes"
	"io"
	"testing"
)

func runSimpleReplacerTest(t *testing.T, content, old, new []byte, cnt int) {
	pseudoReder := bytes.NewReader(content)
	expectedResult := content
	if len(old) > 0 {
		expectedResult = bytes.ReplaceAll(content, old, new)
	}

	r := NewSimpleReplacer(pseudoReder, old, new)
	got, err := io.ReadAll(r)
	if err != nil && err != io.EOF {
		t.Errorf("unpexted error: %v", err)
	}

	if !bytes.Equal(expectedResult, got) {
		t.Errorf("unpexted result. Expect: %s, Got: %s", expectedResult, got)
	}

	if cnt > -1 && cnt != r.Count() {
		t.Errorf("unpexted result. Expect replacement times: %d, Got: %d", cnt, r.Count())
	}
}

func Test_simpleReplacer_Read(t *testing.T) {
	type args struct {
		content []byte
		old     []byte
		new     []byte
		cnt     int
	}
	tests := []args{
		{
			[]byte(`aaaaaaaaaaaaaaoldbbbbbbbbbbboldbbbbbbb`),
			[]byte(`old`),
			[]byte(`new`),
			2,
		},
		{
			[]byte(`111111aaaaaaaaaaaaaaoldbbbbbbbbbbboldbbbbbbol`),
			[]byte(`old`),
			[]byte(`new`),
			2,
		},
		{
			[]byte(`aaaaaaaaaa`),
			[]byte(`a`),
			bytes.Repeat([]byte(`long`), 100),
			10,
		},
		{
			[]byte(`aaaaaaaaaa`),
			[]byte(``),
			[]byte(`long long long long long long long long long string`),
			0,
		},
		{
			[]byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			[]byte("00"),
			[]byte("0"),
			256,
		},
	}
	for _, tt := range tests {
		runSimpleReplacerTest(t, tt.content, tt.old, tt.new, tt.cnt)
	}
}
