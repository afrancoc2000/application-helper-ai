package models

import (
	"testing"

	"github.com/afrancoc2000/application-helper-ai/internal/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMessage(t *testing.T) {
	Convey("Message", t, func() {

		message := models.Message{
			Role:    models.User,
			Content: "Create a hello world html application",
		}

		Convey("ConvertToOpenAIMessages", func() {
			result := models.ConvertToOpenAIMessages([]models.Message{message})
			So(len(result), ShouldEqual, 1)
			So(result[0].Role, ShouldEqual, "user")
			So(result[0].Content, ShouldEqual, "Create a hello world html application")
		})

		Convey("ConvertToAzureOpenAIMessages", func() {
			result := models.ConvertToAzureOpenAIMessages([]models.Message{message})
			So(len(result), ShouldEqual, 1)
			So(result[0].Role, ShouldEqual, "user")
			So(result[0].Content, ShouldEqual, "Create a hello world html application")
		})
	})

}
