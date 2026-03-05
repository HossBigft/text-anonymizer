package maskmanager

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type (
	MaskManager struct {
		maskMap map[string]string
	}
)

var configDir = filepath.Join(os.Getenv("HOME"), ".config", "anonymizer")
var mapFileName = "map.json"
var mapFilePath = filepath.Join(configDir, mapFileName)

func NewMaskManager() *MaskManager {
	newManager := MaskManager{}
	newManager.maskMap, _ = readMasks(mapFilePath)
	return &newManager
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readMasks(path string) (map[string]string, error) {
	masks := make(map[string]string)
	masksFileHandle, err := os.ReadFile(path)
	if err != nil {
		return masks, err
	}
	err = json.Unmarshal(masksFileHandle, &masks)
	return masks, err
}

func (self *MaskManager) GetMaskMap() map[string]string {
	return self.maskMap
}

func (self *MaskManager) SaveMasks() error {
	err := os.MkdirAll(configDir, 0755)
	valueMapJson, err := json.Marshal(self.maskMap)
	err = os.WriteFile(mapFilePath, valueMapJson, 0644)
	return err
}

func (self *MaskManager) UpdateMask(value string, mask string) {
	self.maskMap[value] = mask
}

func (self *MaskManager) GetMask(value string) (string, bool) {
	mask, isSuccesful := self.maskMap[value]
	return mask, isSuccesful
}
