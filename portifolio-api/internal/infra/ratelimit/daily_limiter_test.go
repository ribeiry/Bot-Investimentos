package ratelimit

import "testing"

func TestDailyLimiter_AllowUntilLimit(t *testing.T) {
	l := NewDailyLimiter(3)
	for i := 0; i < 3; i++ {
		if !l.Allow(1) {
			t.Fatalf("chamada %d deveria ser permitida", i+1)
		}
	}
	if l.Allow(1) {
		t.Fatal("chamada acima do limite deveria ser negada")
	}
}

func TestDailyLimiter_IndependentePorUser(t *testing.T) {
	l := NewDailyLimiter(1)
	if !l.Allow(1) {
		t.Fatal("user 1 primeira chamada")
	}
	if l.Allow(1) {
		t.Fatal("user 1 segunda chamada deveria negar")
	}
	if !l.Allow(2) {
		t.Fatal("user 2 deveria permitir")
	}
}
