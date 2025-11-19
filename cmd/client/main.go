// Command gophkeeper-client запускает клиентское CLI-приложение GophKeeper.
//
// Обеспечивает пользовательский интерфейс для работы с хранилищем (регистрация, вход,
// управление записями, синхронизация и т.д.).
package main

import (
	"github.com/fatkulllin/gophkeeper/internal/client/cmd"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
)

// main является точкой входа клиентского приложения.
// Инициализация производится внутри пакета cmd, после чего запускается
// обработка CLI-команд.
func main() {
	defer logger.Log.Sync()
	cmd.Execute()
}
