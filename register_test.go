package spatialite_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/watchedsky-social/go-spatialite"
)

var _ = Describe("Register", func() {
	It("loads spatialite properly", func() {
		db, err := sql.Open("spatialite", ":memory:")
		Expect(err).NotTo(HaveOccurred())
		Expect(db).NotTo(BeNil())

		defer db.Close()

		// If this function completes successfully, it means spatialite is loaded
		_, err = db.Query("SELECT InitSpatialMetadata()")
		Expect(err).NotTo(HaveOccurred())
	})
})
