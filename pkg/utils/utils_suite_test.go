package utils

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubevirtHook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KubevirtHook Suite")
}
