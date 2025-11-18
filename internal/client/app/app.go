package app

import (
	"github.com/fatkulllin/gophkeeper/internal/client/apiclient"
	"github.com/fatkulllin/gophkeeper/internal/client/cryptoutil"
	"github.com/fatkulllin/gophkeeper/internal/client/filemanager"
	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/fatkulllin/gophkeeper/internal/client/store"
)

var CliService *service.Service

func InitApp() {
	apiClient := apiclient.NewApiClient(10)
	filemanager := filemanager.NewFileManager()
	cryptoUtil := cryptoutil.NewCryptoUtil()
	boltDB, _ := store.NewBoltDB()
	CliService = service.NewService(apiClient, filemanager, boltDB, cryptoUtil)
}
