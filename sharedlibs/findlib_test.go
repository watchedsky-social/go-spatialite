package sharedlibs_test

import (
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/watchedsky-social/go-spatialite/sharedlibs"
)

var _ = Describe("Findlib", func() {
	It("Finds a common library in a common location for linux", func() {
		if runtime.GOOS != "linux" {
			Skip("This test is linux only")
		}

		files, err := sharedlibs.FindSharedLibraryFiles("libm", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(files)).To(BeNumerically(">", 0))
	})

	It("Finds a common library in a common location for darwin", func() {
		if runtime.GOOS != "darwin" {
			Skip("This test is darwin only")
		}

		files, err := sharedlibs.FindSharedLibraryFiles("libclang", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(files)).To(BeNumerically(">", 0))
	})

	It("Finds a common library in a common location for windows", func() {
		if runtime.GOOS != "windows" {
			Skip("This test is windows only")
		}

		files, err := sharedlibs.FindSharedLibraryFiles("kernel32", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(files)).To(BeNumerically(">", 0))
	})

	It("Finds a library in a specific search path", func() {
		path, _ := filepath.Abs("testdata")

		os.Setenv("GO_SPATIALITE_SEARCH_PATH", path)
		files, err := sharedlibs.FindSharedLibraryFiles("libfoo", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(files)).To(BeNumerically(">", 0))
	})
})
