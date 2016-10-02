package csp

import (
	//"fmt"
	"gopkg.in/tylerb/is.v1"
	"io/ioutil"
	"os"
	"testing"
)

func TestCmsDecoder(t *testing.T) {
	is := is.New(t)

	f, err := os.Open("/tmp/logical.cms")
	is.NotErr(err)
	defer f.Close()

	msg, err := NewCmsDecoder()
	is.NotErr(err)
	o, err := os.Create("/tmp/logical.bin")
	is.NotErr(err)
	defer o.Close()

	n, err := msg.Decode(o, f)
	is.NotErr(err)
	is.NotZero(n)

	store, err := msg.CertStore()
	is.NotErr(err)
	is.NotZero(store)

	for _, c := range store.Certs() {
		is.Lax().NotErr(msg.Verify(c)) // XXX can not verify RSA on Linux
	}
	is.NotErr(msg.Close())
}

func TestCmsDetached(t *testing.T) {
	is := is.New(t)

	sig, err := ioutil.ReadFile("/tmp/data1.p7s")
	is.NotErr(err)
	data, err := os.Open("/tmp/data1.bin")
	is.NotErr(err)
	msg, err := NewCmsDecoder(sig)
	is.NotErr(err)
	_, err = msg.Decode(nil, data)
	is.NotErr(err)

	store, err := msg.CertStore()
	is.NotErr(err)
	is.NotZero(store)

	for _, c := range store.Certs() {
		is.Lax().NotErr(msg.Verify(c))
	}
	is.NotErr(msg.Close())
}
