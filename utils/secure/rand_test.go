package secure

import "testing"

func TestRandomStr(t *testing.T) {
	for i := 0; i <= 36; i++ {
		t.Logf("rand %d: %s", i, RandomStr(i))
	}
}

func BenchmarkRandomStrLen5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomStr(10)
	}
}

func BenchmarkRandomStrLen10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomStr(10)
	}
}

func BenchmarkRandomStrLen15(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomStr(15)
	}
}

func BenchmarkRandomStrLen20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomStr(20)
	}
}

func BenchmarkRandomStrLen36(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomStr(36)
	}
}
