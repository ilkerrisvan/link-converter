package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	a := App{}
	actual := a.initialize(os.Getenv("CONNECTION_STRING"))
	assert.Nil(t, actual, "APP could not initialize.")
}
