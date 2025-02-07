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
		f.Set("foobar")
		f.Set("qwerty")
		assertBool(t, f.Check("123456"), false)
		assertBool(t, f.Check("foobar"), true)
		assertBool(t, f.Check("hello"), false)
		assertBool(t, f.Check("qwerty"), true)
		assertBool(t, f.Check("654321"), false)
	})
}
