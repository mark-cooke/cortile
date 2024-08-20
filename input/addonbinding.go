package input

import (
	"os"

	"os/exec"
	"path/filepath"

	"github.com/leukipp/cortile/v2/common"
	"github.com/leukipp/cortile/v2/desktop"

	log "github.com/sirupsen/logrus"
)

func BindAddons(tr *desktop.Tracker) {
	if common.HasFlag("disable-addons-folder") {
		return
	}

	// check if addons folder exists
	configFolderPath := common.ConfigFolderPath(common.Build.Name)
	addonsFolderPath := filepath.Join(configFolderPath, "addons")
	if _, err := os.Stat(addonsFolderPath); os.IsNotExist(err) {
		return
	}

	// read files in addons folder
	files, err := os.ReadDir(addonsFolderPath)
	if err != nil {
		log.Warn("Error reading addons: ", addonsFolderPath)
		return
	}

	// run files in addons folder
	for _, file := range files {
		addonFilePath := filepath.Join(addonsFolderPath, file.Name())
		log.Info("Execute addon ", addonFilePath)

		// execute addon scripts
		addon := exec.Command(addonFilePath)
		addon.Stdout = os.Stdout
		addon.Stderr = os.Stderr

		if err = addon.Start(); err != nil {
			log.Warn("Error executing addon: ", err)
		}
	}
}