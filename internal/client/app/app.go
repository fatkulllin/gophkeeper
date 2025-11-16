package app

import (
	"github.com/fatkulllin/gophkeeper/internal/client/apiclient"
	"github.com/fatkulllin/gophkeeper/internal/client/filemanager"
	"github.com/fatkulllin/gophkeeper/internal/client/service"
)

var CliService *service.Service

func InitApp() {
	apiClient := apiclient.NewApiClient(10)
	filemanager := filemanager.NewFileManager()
	CliService = service.NewService(apiClient, filemanager)

}
