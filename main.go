package main

import (
	"os"

	"github.com/gunturaf/sukab-property/domain/property"
	"github.com/gunturaf/sukab-property/infra"
	"github.com/gunturaf/sukab-property/server"
)

func isDBEnableSSL() bool {
	if os.Getenv("DB_ENABLE_SSL") != "" {
		return true
	}
	return false
}

func main() {
	db := infra.NewDB(infra.DBAttributes{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		EnableSSL:    isDBEnableSSL(),
	})

	repository := property.NewRepo(db)
	importer := property.NewImporter(repository)
	lister := property.NewLister(repository)

	httpServer := server.New(importer, lister)
	httpServer.Run(os.Getenv("SERVER_LISTEN_ADDR"))
}
