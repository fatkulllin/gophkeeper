package app

import (
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/apiclient"
	"github.com/fatkulllin/gophkeeper/internal/client/cryptoutil"
	"github.com/fatkulllin/gophkeeper/internal/client/filemanager"
	"github.com/fatkulllin/gophkeeper/internal/client/fs"
	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/fatkulllin/gophkeeper/internal/client/store"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// InitApp создаёт и инициализирует все зависимости клиентского приложения.
// После успешного выполнения функция формирует экземпляр CliService,
// через который CLI-команды взаимодействуют с API и локальным хранилищем.
func InitApp() (*service.Service, error) {

	appDir, err := fs.PrepareAppDir()

	if err != nil {
		return nil, err
	}

	logger.Log.Debug("config dir", zap.String("dir", appDir))

	apiClient := apiclient.NewApiClient(10)
	fm := filemanager.NewFileManager(appDir)
	cryptoUtil := cryptoutil.NewCryptoUtil()
	boltDB, err := store.NewBoltDB(appDir)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize local storage: %v", err)
	}

	svc := service.NewService(apiClient, fm, boltDB, cryptoUtil)
	return svc, nil
}
