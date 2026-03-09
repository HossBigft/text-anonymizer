package maskmanager

import (
	patternmanager "anonymizer/patternManager"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/lucasjones/reggen"
)

type (
	MaskManager struct {
		valueToMaskMap map[string]string
	}
)

var configDir = filepath.Join(os.Getenv("HOME"), ".config", "anonymizer")
var mapFileName = "map.json"
var mapFilePath = filepath.Join(configDir, mapFileName)

func NewMaskManager() *MaskManager {
	newManager := MaskManager{}
	newManager.valueToMaskMap, _ = readMasks(mapFilePath)
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
	return self.valueToMaskMap
}

func (self *MaskManager) SaveMasks() error {
	err := os.MkdirAll(configDir, 0755)
	valueMapJson, err := json.Marshal(self.valueToMaskMap)
	err = os.WriteFile(mapFilePath, valueMapJson, 0644)
	return err
}

func (self *MaskManager) UpdateMask(value string, mask string) {
	self.valueToMaskMap[value] = mask
}

func (self *MaskManager) MapValuesToMasks(match patternmanager.PatternMatch) map[string]string {
	isMasksUpdated := false
	for _, value := range match.Matches {
		mask, present := self.valueToMaskMap[value]
		if present == false {
			mask = self.GetRandomStringByRegex(match.MaskPattern.Regex)
			self.valueToMaskMap[value] = mask
			isMasksUpdated = true
		}
	}
	if isMasksUpdated {
		self.SaveMasks()
	}
	return self.valueToMaskMap
}

func (self *MaskManager) GetRandomStringByRegex(regex string, maxLength_optional ...int) string {
	var maxLength int
	if len(maxLength_optional) == 0 {
		maxLength = 7
	} else {
		maxLength = maxLength_optional[0]
	}
	randomString, _ := reggen.Generate(regex, maxLength)
	return randomString
}

func (self *MaskManager) GetMasksToValuesMap() map[string]string {
	reverseMap := make(map[string]string)
	for value, mask := range self.valueToMaskMap {
		reverseMap[mask] = value
	}
	return reverseMap
}
