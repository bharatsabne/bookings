package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/some", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("Got Invalid should have got valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/some", nil)
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")
	r.PostForm = postData
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows Valid when reqiured fileds missing")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	if form.Has("x") {
		t.Error("X iis not present but showing")
	}

	postData.Add("a", "a")
	form = New(postData)
	if !form.Has("a") {
		t.Error("Expected field 'a' to be present and not empty")
	}
}

func TestFormMinimumLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	form.MinimumLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length of Non existed filed")
	}
	ifError := form.Errors.Get("x")
	if ifError == "" {
		t.Error("should have an error")
	}
	postData = url.Values{}
	postData.Add("a", "111r1")
	form = New(postData)
	form.MinimumLength("a", 15)
	if form.Valid() {
		t.Error("Shows min length of 15 met when data is shorter")
	}

	postData = url.Values{}
	postData.Add("b", "111r1")
	form = New(postData)
	form.MinimumLength("b", 1)
	if !form.Valid() {
		t.Error("Shows min length of 1 is not met when it is")
	}
	ifError = form.Errors.Get("b")
	if ifError != "" {
		t.Error("Error should be empty")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("Forms shows valid email for non exited filed")
	}
	postData = url.Values{}
	postData.Add("email", "bharat@sabne.com")
	form = New(postData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Got an invalid when we should not have")
	}
	postData = url.Values{}
	postData.Add("email", "bharat@sab@ne.com")
	form = New(postData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got an valid when we should not have")
	}
}
