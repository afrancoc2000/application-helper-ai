package models

import (
	"testing"

	"github.com/afrancoc2000/application-helper-ai/internal/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppFile(t *testing.T) {
	Convey("AppFile", t, func() {

		Convey("AppFileFromString", func() {
			openaiResult := `[{"fileName":"main.tf","filePath":"./","fileContent":"\n# Configure the Azure provider\nprovider \"azurerm\" {\n\tfeatures {}\n}\n\n# Create a resource group\nresource \"azurerm_resource_group\" \"aks\" {  \n\tname     = var.resource_group_name\n\tlocation = var.resource_group_location\n}\n"}]`
			files, err := models.AppFileFromString(openaiResult)
			So(err, ShouldBeNil)
			So(len(files), ShouldEqual, 1)
			So(files[0].Name, ShouldEqual, "main.tf")
			So(files[0].Path, ShouldEqual, "./")
			So(files[0].Content, ShouldEqual, "\n# Configure the Azure provider\nprovider \"azurerm\" {\n\tfeatures {}\n}\n\n# Create a resource group\nresource \"azurerm_resource_group\" \"aks\" {  \n\tname     = var.resource_group_name\n\tlocation = var.resource_group_location\n}\n")
		})
	})

}
