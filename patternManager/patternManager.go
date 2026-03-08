package patternmanager

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

type (
	MaskPattern struct {
		Name  string `json:"name"`
		Regex string `json:"regex"`
	}

	PatternManager struct {
		maskPatterns map[string]MaskPattern
	}

	PatternMatch struct {
		Matches     []string
		MaskPattern MaskPattern
	}
)

func NewPatternManager() *PatternManager {
	newManager := PatternManager{}
	newManager.maskPatterns = make(map[string]MaskPattern)
	newManager.loadConfig()
	return &newManager
}

var configDir = filepath.Join(os.Getenv("HOME") + "/.config/anonymizer/")
var configFilePath = os.Getenv("HOME") + "/.config/anonymizer/maskPatterns.json"

func readConfig(path string) ([]MaskPattern, error) {
	var patterns []MaskPattern
	patternsFileHandle, err := os.ReadFile(path)
	if err != nil {
		return patterns, err
	}
	err = json.Unmarshal(patternsFileHandle, &patterns)
	return patterns, err
}

func (self *PatternManager) loadConfig() error {
	IPV4_REGEX := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}`
	FQDN_REGEX := `(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)`
	patterns, err := readConfig(configFilePath)
	if err != nil {
		patterns = append(patterns, MaskPattern{Name: "ipv4", Regex: IPV4_REGEX})
		patterns = append(patterns, MaskPattern{Name: "fqdn", Regex: FQDN_REGEX})
		self.SavePatterns()
	}
	for _, pattern := range patterns {
		self.maskPatterns[pattern.Name] = pattern
	}
	return err
}

func (self *PatternManager) GetPatterns() []MaskPattern {
	maskPatterns := []MaskPattern{}
	for _, pattern := range self.maskPatterns {
		maskPatterns = append(maskPatterns, pattern)
	}
	return maskPatterns
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func (self *PatternManager) SavePatterns() error {
	err := os.MkdirAll(configDir, 0755)
	valueMapJson, err := json.Marshal(self.maskPatterns)
	err = os.WriteFile(configFilePath, valueMapJson, 0644)
	return err
}

func (self *PatternManager) MapValuesToPatterns(rawLine string) ([]PatternMatch, error) {
	var valuesToMaskMap []PatternMatch
	var err error
	for _, pattern := range self.GetPatterns() {
		var regex *regexp.Regexp
		regex, err = regexp.Compile(pattern.Regex)
		sensitive_values := regex.FindAllString(rawLine, -1)
		valuesToMaskMap = append(valuesToMaskMap, PatternMatch{MaskPattern: pattern, Matches: sensitive_values})

	}
	return valuesToMaskMap, err
}

func (self *PatternManager) AddPattern(pattern MaskPattern) {
	self.maskPatterns[pattern.Name] = pattern
}
func (self *PatternManager) RemovePatternByName(name string) (MaskPattern, error) {
	var err error
	pattern, present := self.maskPatterns[name]
	if !present {
		err = errors.New("Pattern not found")
	}
	delete(self.maskPatterns, name)
	return pattern, err
}
