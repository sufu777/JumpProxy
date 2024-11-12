package main

import "testing"

func BenchmarkNextAAZ345(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NextAAZ345()
	}
}

func TestNextAAZ345(t *testing.T) {
	println(NextAAZ345())
	println(len(NextAAZ345()))
}
