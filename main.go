package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/koltyakov/gosip"
	ondemand "github.com/koltyakov/gosip-sandbox/strategies/ondemand"
	"github.com/koltyakov/gosip/api"
)

var flagSiteURL = flag.String("site-url", "", "The site URL to use.")
var flagFolder = flag.String("folder", "", "The folder to walk, e.g. /sites/MySite/Shared Documents/General")

func main() {
	flag.Parse()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	// Use whatever auth you like.
	spAuth := &ondemand.AuthCnfg{
		SiteURL: *flagSiteURL,
	}
	sp := api.NewSP(&gosip.SPClient{
		AuthCnfg: spAuth,
	})

	// Read folders.
	folders, err := sp.Web().GetFolder(*flagFolder).Folders().Get()
	if err != nil {
		log.Error("failed to get folder", slog.String("folder", *flagFolder), slog.Any("error", err))
		os.Exit(1)
	}
	for _, folder := range folders.Data() {
		log.Info("accessed folder", slog.String("folder", folder.Data().ServerRelativeURL))
	}

	// Read files in the folder.
	files, err := sp.Web().GetFolder(*flagFolder).Files().Get()
	if err != nil {
		log.Error("failed to get files", slog.String("folder", *flagFolder), slog.Any("error", err))
		os.Exit(1)
	}
	for _, file := range files.Data() {
		// Get the file metadata.
		file, err := sp.Web().GetFile(file.Data().ServerRelativeURL).Get()
		if err != nil {
			log.Error("failed to get file", slog.String("file", file.Data().ServerRelativeURL), slog.Any("error", err))
			os.Exit(1)
		}
		log.Info("accessed file metadata", slog.String("file", file.Data().ServerRelativeURL))
		// Get the file data.
		rc, err := sp.Web().GetFile(file.Data().ServerRelativeURL).GetReader()
		if err != nil {
			log.Error("failed to get file data", slog.String("file", file.Data().ServerRelativeURL), slog.Any("error", err))
			os.Exit(1)
		}
		// Create a hash of the data to prove we can read it.
		hash := sha256.New()
		_, err = io.Copy(hash, rc)
		if err != nil {
			log.Error("failed to hash file data", slog.String("file", file.Data().ServerRelativeURL), slog.Any("error", err))
			os.Exit(1)
		}
		log.Info("read file data", slog.String("file", file.Data().ServerRelativeURL), slog.String("hash", fmt.Sprintf("%x", hash.Sum(nil))))
		if err = rc.Close(); err != nil {
			log.Error("failed to close file data", slog.String("file", file.Data().ServerRelativeURL), slog.Any("error", err))
			os.Exit(1)
		}
	}
	log.Info("done")
}
