package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var maH myHandler
	h := NoSurf(&maH)

	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Sprintf("Type is not http.Handlr, but is %T", v))
	}
}
func TestSessionLoad(t *testing.T) {
	var maH myHandler
	h := SessionLoad(&maH)

	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Sprintf("Type is not http.Handlr, but is %T", v))
	}
}
