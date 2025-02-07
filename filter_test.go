package bloom

import "testing"

func assertBool(tb testing.TB, value, expected bool) {
	if value != expected {
		tb.Errorf("expected %v, got %v", expected, value)
	}
}

func TestFilter(t *testing.T) {
	t.Run("", func(t *testing.T) {
		f, err := NewFilter(&Config{
			Size:            1000,
			Hasher:          &hasherStringCRC64{},
			HashChecksLimit: 3,
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = f.Set("foobar")
		_ = f.Set("qwerty")
		assertBool(t, f.Check("123456"), false)
		assertBool(t, f.Check("foobar"), true)
		assertBool(t, f.Check("hello"), false)
		assertBool(t, f.Check("qwerty"), true)
		assertBool(t, f.Check("654321"), false)
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("", func(b *testing.B) {
		b.ReportAllocs()
		f, _ := NewFilter(&Config{
			Size:            1000,
			Hasher:          &hasherStringCRC64{},
			HashChecksLimit: 3,
		})
		_ = f.Set("foobar")
		_ = f.Set("qwerty")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f.Check("foobar")
		}
	})
}
