package chesscode

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/notnil/chess"
)

func TestCycle(t *testing.T) {
	src := rand.NewSource(42)
	r := rand.New(src)
	for i := 0; i < 3000; i++ {
		s1 := randomString(r, r.Intn(24))
		s1 = strings.TrimSpace(s1)
		t.Log(s1)
		b, err := Encode(s1)
		if err != nil {
			t.Fatal(err)
		}
		s2, err := Decode(b)
		if err != nil {
			t.Fatal(err)
		}
		s2 = strings.TrimSpace(s2)
		s1 = strings.ToUpper(s1)
		if s1 != s2 {
			t.Log(b.String())
			t.Log(b.Draw())
			t.Fatalf("encode / decode cycle expected %s but got %s", s1, s2)
		}
	}
}

func TestOneCycle(t *testing.T) {
	s1 := "Berlin tomorrow at 2pm."
	s1 = strings.TrimSpace(s1)
	t.Log(s1)
	b, err := Encode(s1)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	s2 = strings.TrimSpace(s2)
	s1 = strings.ToUpper(s1)
	t.Log(b.String())
	t.Log(b.Draw())
	if s1 != s2 {
		t.Fatalf("encode / decode cycle expected %s but got %s", s1, s2)
	}
}

func TestInvalidDecode(t *testing.T) {
	b := &chess.Board{}
	if err := b.UnmarshalText([]byte("7K/8/8/8/8/8/8/8")); err != nil {
		t.Fatal(err)
	}
	s, err := Decode(b)
	if err == nil {
		t.Fatal("expected error but didn't get one")
	}
	if len(s) > 0 {
		t.Fatal("expected empty string but didn't get one")
	}
}

func BenchmarkCycle(b *testing.B) {
	s := ".G2QMYLCKQ569A325HT2H0R"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := Encode(s)
		Decode(b)
	}
}

func randomString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
