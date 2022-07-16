package ch12

import "testing"

// go test -v ch12/*
func Test_standardLib(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			standardLib()
			t.Log("失敗時やgo test -vのときだけ表示されます")
			t.Fatal("メッセージとともにテストを失敗させます")
		})
	}
}
