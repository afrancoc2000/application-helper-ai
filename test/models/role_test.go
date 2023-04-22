package models

import (
	"testing"

	"github.com/afrancoc2000/application-helper-ai/internal/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRole(t *testing.T) {
	Convey("Role", t, func() {

		Convey("user", func() {
			result := models.User.String()
			So(result, ShouldEqual, "user")
		})

		Convey("assistant", func() {
			result := models.Assistant.String()
			So(result, ShouldEqual, "assistant")
		})

		Convey("system", func() {
			result := models.System.String()
			So(result, ShouldEqual, "system")
		})
	})

}
