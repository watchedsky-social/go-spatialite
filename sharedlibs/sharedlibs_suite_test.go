package sharedlibs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSharedlibs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sharedlibs Suite")
}
