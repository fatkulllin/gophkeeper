/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/fatkulllin/gophkeeper/internal/client/cmd"
	"github.com/fatkulllin/gophkeeper/logger"
)

func main() {
	defer logger.Log.Sync()
	cmd.Execute()
}
