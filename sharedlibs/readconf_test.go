package sharedlibs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/watchedsky-social/go-spatialite/sharedlibs"
)

var _ = Describe("Readconf", func() {
	It("reads a conf file", func() {
		files, err := sharedlibs.ReadConfFile("testdata/conf/ld.so.conf")
		Expect(err).NotTo(HaveOccurred())
		Expect(files).To(HaveExactElements("/opt/homebrew/lib", "/opt/homebrew/Cellar", "/Applications/Xcode.app"))
	})
})
