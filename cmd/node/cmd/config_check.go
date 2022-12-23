package cmd

import (
	"fs.video/blockchain/app/genesis"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


func checkGenesisFile(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	configDir := filepath.Join(path, "config")              
	genesisPath := filepath.Join(configDir, "genesis.json") 
	dataPath := filepath.Join(path, "data")                 
	statePath := filepath.Join(configDir, "priv_validator_state.json")
	upInfoPath := filepath.Join(configDir, "upInfo.txt") 
	upInfoExist, _ := CheckSourceExist(upInfoPath)
	if exist, _ := CheckSourceExist(genesisPath); exist {
		
		genesisBuf, err := ioutil.ReadFile(genesisPath)
		if err != nil {
			log.WithError(err).WithField("path", genesisPath).Error("ioutil.ReadFile")
			return err
		}
		genesisContent := string(genesisBuf)
		if strings.Contains(genesisContent, core.ChainID) && !upInfoExist { 
			return nil
		}
	}
	_, err := os.Create(upInfoPath) 
	if err != nil {
		log.WithError(err).WithField("file", "upInfo.txt").Error("File Create")
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	genesisPathNew := filepath.Join(filepath.Join(pwd, "db"), "genesis.json")
	genesisBuf, err := ioutil.ReadFile(genesisPathNew)
	if err != nil {
		log.WithError(err).WithField("path", genesisPathNew).Error("ioutil.ReadFile")
		return err
	}
	
	err = util.FilePermissionChange(genesisPath)
	if err != nil {
		log.WithError(err).WithField("path", genesisPath).Error("FilePermissionChange")
		return err
	}

	err = ioutil.WriteFile(genesisPath, genesisBuf, 0)
	if err != nil {
		log.WithError(err).WithField("file", genesisPath).Error("ioutil.WriteFile")
		return err
	}

	
	if exist, _ := CheckSourceExist(dataPath); exist {
		
		err = util.RenameBackup(path, "data")
		if err != nil {
			log.WithError(err).WithField("dataPath", dataPath).Error("DirRename")
			
			err = ioutil.WriteFile(upInfoPath, []byte(dataPath), 0)
			if err != nil {
				log.WithError(err).WithField("dataPath", dataPath).Error("WriteString upInfo.txt")
				return err
			}
			return err
		}
		//data
		if _, err := os.Stat(dataPath); os.IsNotExist(err) {
			err = os.MkdirAll(dataPath, os.ModePerm) 
			if err != nil {
				log.WithError(err).Error("Unable to create directory " + dataPath)
				return err
			}
		}
	}

	//priv_validator_state.json
	if exist, _ := CheckSourceExist(statePath); exist {
		err = ioutil.WriteFile(statePath, []byte(genesis.Validator_state), 0)
		if err != nil {
			log.WithError(err).WithField("statePath", statePath).Error("ioutil.WriteFile")
			return err
		}
	}
	
	err = os.Remove(upInfoPath)
	if err != nil {
		log.WithError(err).WithField("upInfoPath", upInfoPath).Error("Remove upInfo.txt")
		return err
	}
	return nil
}

//config.toml,
func checkConfigFile(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	configDir := filepath.Join(path, "config")
	configPath := filepath.Join(configDir, "config.toml")

	
	if exist, _ := CheckSourceExist(configPath); exist {
		return nil
	}
	err := ioutil.WriteFile(configPath, []byte(genesis.ConfigToml), 0)
	if err != nil {
		log.WithError(err).WithField("file", configPath).Error("ioutil.WriteFile")
	}
	return err
}

//client.toml,
func checkClientToml(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	configDir := filepath.Join(path, "config")
	configPath := filepath.Join(configDir, "client.toml")

	
	if exist, _ := CheckSourceExist(configPath); exist {
		return nil
	}
	err := ioutil.WriteFile(configPath, []byte(genesis.ClientToml), 0)
	if err != nil {
		log.WithError(err).WithField("file", configPath).Error("ioutil.WriteFile")
	}
	return err
}

func checkValidatorStateJson(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	configStateJsonFile := filepath.Join(path, "config", "priv_validator_state.json")
	dataStateJsonFile := filepath.Join(path, "data", "priv_validator_state.json")
	
	if exist, _ := CheckSourceExist(configStateJsonFile); exist {
		return nil
	}
	if exist, _ := CheckSourceExist(dataStateJsonFile); exist {
		err := os.Rename(dataStateJsonFile, configStateJsonFile)
		if err != nil {
			log.WithError(err).WithField("dataStateJsonFile", dataStateJsonFile).Error("Rename")
			return err
		}
	}
	return nil
}

//app.toml
func checkAppToml(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	configDir := filepath.Join(path, "config")
	appPath := filepath.Join(configDir, "app.toml")
	if exist, _ := CheckSourceExist(appPath); exist {
		appBuf, err := ioutil.ReadFile(appPath)
		if err != nil {
			log.WithError(err).WithField("path", appPath).Error("ioutil.ReadFile")
			return err
		}
		appContent := string(appBuf)

		if strings.Contains(appContent, "[evm]") {
			return nil
		}
		newAppContent := strings.ReplaceAll(appContent, appContent, genesis.AppToml)
		
		err = ioutil.WriteFile(appPath, []byte(newAppContent), 0)
		if err != nil {
			log.WithError(err).WithField("path", appPath).Error("ioutil.WriteFile")
			return err
		}
		return nil
	}
	return nil
}

func replaceConfig(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	configDir := filepath.Join(path, "config")
	appPath := filepath.Join(configDir, "app.toml")
	if exist, err := CheckSourceExist(appPath); !exist {
		return err
	}
	genesisPath := filepath.Join(configDir, "genesis.json")
	if exist, err := CheckSourceExist(genesisPath); !exist {
		return err
	}
	configPath := filepath.Join(configDir, "config.toml")
	if exist, err := CheckSourceExist(appPath); !exist {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	genesisPathNew := filepath.Join(filepath.Join(pwd, "db"), "genesis.json")
	genesisBuf, err := ioutil.ReadFile(genesisPathNew)
	if err != nil {
		log.WithError(err).WithField("path", genesisPathNew).Error("ioutil.ReadFile")
		return err
	}
	
	err = ioutil.WriteFile(appPath, []byte(genesis.AppToml), 0)
	if err != nil {
		log.WithError(err).WithField("path", appPath).Error("ioutil.WriteFile")
		return err
	}
	err = ioutil.WriteFile(genesisPath, genesisBuf, 0)
	if err != nil {
		log.WithError(err).WithField("path", genesisPath).Error("ioutil.WriteFile")
		return err
	}
	err = ioutil.WriteFile(configPath, []byte(genesis.ConfigToml), 0)
	if err != nil {
		log.WithError(err).WithField("path", configPath).Error("ioutil.WriteFile")
		return err
	}
	return err
}


func CheckSourceExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
