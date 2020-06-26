package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

//NewTestServer return *Server test object
func NewTestServer() *Server {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	return New(":8080", logger)
}
