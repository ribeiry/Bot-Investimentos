package portifolio

import "testing"

func TestValidateNarrative_TickerValido(t *testing.T) {
	if !validateNarrative("BBSE3 subiu 2%", []string{"BBSE3"}) {
		t.Fatal("esperava true")
	}
}

func TestValidateNarrative_TickerInvalido(t *testing.T) {
	if validateNarrative("XYZW9 disparou hoje", []string{"BBSE3"}) {
		t.Fatal("esperava false")
	}
}

func TestValidateNarrative_SoBenchmarks(t *testing.T) {
	if !validateNarrative("IBOV superou S&P nesta semana", nil) {
		t.Fatal("esperava true")
	}
}

func TestValidateNarrative_SemTickers(t *testing.T) {
	if !validateNarrative("Sua carteira está estável", nil) {
		t.Fatal("esperava true")
	}
}

func TestValidateNarrative_TickerMinusculoIgnorado(t *testing.T) {
	if !validateNarrative("bbse3 apareceu em minúsculo", nil) {
		t.Fatal("esperava true (minúsculas não são validadas)")
	}
}

func TestValidateNarrative_MixValidoInvalido(t *testing.T) {
	if validateNarrative("BBSE3 e ABCD1 subiram", []string{"BBSE3"}) {
		t.Fatal("esperava false porque ABCD1 não está permitido")
	}
}
