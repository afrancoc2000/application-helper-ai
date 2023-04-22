package models

import (
	"testing"

	"github.com/afrancoc2000/application-helper-ai/internal/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDeployment(t *testing.T) {
	Convey("Deployment", t, func() {

		deployment := models.Gpt4_0314

		Convey("String", func() {
			result := deployment.String()
			So(result, ShouldEqual, "gpt-4-0314")
		})

		Convey("MaxTokens", func() {
			result := deployment.MaxTokens()
			So(result, ShouldEqual, 8192)
		})

		Convey("IsChat", func() {
			result := deployment.IsChat()
			So(result, ShouldEqual, true)
		})

		Convey("DeploymentFromName", func() {
			testDeployment, err := models.DeploymentFromName("gpt-4-0314")
			So(err, ShouldBeNil)
			So(testDeployment, ShouldEqual, models.Gpt4_0314)
		})
	})

}
