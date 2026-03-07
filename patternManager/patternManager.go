package patternmanager

import (
	"encoding/json"
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
		maskPatterns []MaskPattern
	}
)

func NewPatternManager() *PatternManager {
	newManager := PatternManager{}
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

func (m *PatternManager) loadConfig() error {
	IPV4_REGEX := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}`
	FQDN_REGEX := `(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)`
	patterns, err := readConfig(configFilePath)
	if err != nil {
		patterns = append(patterns, MaskPattern{Name: "ipv4", Regex: IPV4_REGEX})
		patterns = append(patterns, MaskPattern{Name: "fqdn", Regex: FQDN_REGEX})
		m.maskPatterns = append(m.maskPatterns, patterns...)
		m.SavePatterns()
	}
	m.maskPatterns = append(m.maskPatterns, patterns...)
	return err
}

func (m *PatternManager) GetPatterns() []MaskPattern {
	return m.maskPatterns
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func (m *PatternManager) SavePatterns() error {
	err := os.MkdirAll(configDir, 0755)
	valueMapJson, err := json.Marshal(m.maskPatterns)
	err = os.WriteFile(configFilePath, valueMapJson, 0644)
	return err
}

func (self *PatternManager) MapSensitiveValuesToPatterns(rawLine string) (map[string]MaskPattern, error) {
	valuesToMaskMap := make(map[string]MaskPattern)
	var err error
	for _, pattern := range self.GetPatterns() {
		var regex *regexp.Regexp
		regex, err = regexp.Compile(pattern.Regex)
		sensitive_values := regex.FindAllString(rawLine, -1)
		for _, value := range sensitive_values {
			valuesToMaskMap[value] = pattern
		}
	}
	return valuesToMaskMap, err
}
