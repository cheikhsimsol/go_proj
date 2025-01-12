package main

import "testing"

func BenchmarkFuncTypeA(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		FuncTypeA()
	}

}

func BenchmarkFuncTypeB(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		FuncTypeB()
	}
}
