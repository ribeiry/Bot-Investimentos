package ratelimit

import "testing"

func TestGlobalLimiter_MinuteCap(t *testing.T) {
	l := NewGlobalLimiter(2, 100)
	if !l.Allow() || !l.Allow() {
		t.Fatal("primeiras 2 deveriam passar")
	}
	if l.Allow() {
		t.Fatal("terceira no mesmo minuto deve negar")
	}
}

func TestGlobalLimiter_DayCap(t *testing.T) {
	l := NewGlobalLimiter(100, 2)
	if !l.Allow() || !l.Allow() {
		t.Fatal("primeiras 2 deveriam passar")
	}
	if l.Allow() {
		t.Fatal("acima do cap diário deve negar")
	}
}
