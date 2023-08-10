package test

import (
	"testing"
)

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Error("Add(1, 2) should be equal to 3")
	}
	t.Log("test add success")
}
