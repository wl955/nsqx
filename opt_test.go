package mq

import "testing"

func TestWriter(t *testing.T) {
	if nil == opt.writer {
		t.Log("writer = nil")
	}
	if opt.writer != nil {
		t.Log("writer != nil")
	}
	//    opt_test.go:7: writer = nil

	var u *User
	var us []*User
	for i := 21; i < 26; i++ {
		u = &User{Age: i}
		us = append(us, u)
	}
	for _, v := range us {
		t.Logf("%+v", *v)
	}
	//    opt_test.go:22: {Age:21}
	//    opt_test.go:22: {Age:22}
	//    opt_test.go:22: {Age:23}
	//    opt_test.go:22: {Age:24}
	//    opt_test.go:22: {Age:25}
}

type User struct {
	Age int
}
