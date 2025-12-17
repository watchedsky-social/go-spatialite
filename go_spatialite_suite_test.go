package spatialite_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoSpatialite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoSpatialite Suite")
}
