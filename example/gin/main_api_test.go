package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIndexUsers(t *testing.T) {
	ts := httptest.NewServer(getEngine())
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/users?limit=30")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestShowUsers(t *testing.T) {
	ts := httptest.NewServer(getEngine())
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/users/1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestCreateUsers(t *testing.T) {
	ts := httptest.NewServer(getEngine())
	defer ts.Close()
	values := url.Values{}
	values.Add("name", "gotokatsuya")
	resp, err := http.PostForm(ts.URL+"/users", values)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestUpdateUsers(t *testing.T) {
	ts := httptest.NewServer(getEngine())
	defer ts.Close()
	type req struct {
		Name string `json:"name"`
	}
	r := req{Name: "gotokatsuya"}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		t.Fatal(err)
	}
	newReq, err := http.NewRequest("PUT", ts.URL+"/users", buf)
	if err != nil {
		t.Fatal(err)
	}
	newReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
