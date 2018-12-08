package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoplay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goplay Suite")
}
