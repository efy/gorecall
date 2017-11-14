package main

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password := "test"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	match := CheckPasswordHash(password, hash)
	if match != true {
		t.Error("expected", true, "got", "match")
	}
}

func TestCheckPasswordHashPrecomputed(t *testing.T) {
	password := "test"
	hash := "$2a$10$StmyKkvbf4f/XJ9cQDv/IOMfU.7Q8Kzns/fwqIOutOW.52lAwAnAq"
	match := CheckPasswordHash(password, hash)
	if match != true {
		t.Error("expected", true, "got", match)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword("test")
	}
}
