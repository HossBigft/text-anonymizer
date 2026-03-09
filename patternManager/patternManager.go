package patternmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type (
	Pattern struct {
		Name  string `json:"name"`
		Regex string `json:"regex"`
	}

	PatternManager struct {
		maskPatterns   map[string]Pattern
		filterPatterns map[string]Pattern
	}

	PatternMatch struct {
		Matches     []string
		MaskPattern Pattern
	}
)

func NewPatternManager() *PatternManager {
	newManager := PatternManager{}
	newManager.maskPatterns = make(map[string]Pattern)
	newManager.loadPatterns()

	newManager.filterPatterns = make(map[string]Pattern)
	newManager.loadFilters()
	return &newManager
}

var patternsConfigFileName = "maskPatterns.json"
var filtersConfigFileName = "filterPatterns.json"
var configDir = filepath.Join(os.Getenv("HOME"), "/.config/ae/")
var patternsFilePath = filepath.Join(configDir, patternsConfigFileName)
var filtersFilePath = filepath.Join(configDir, filtersConfigFileName)

func readConfig(path string) (map[string]Pattern, error) {
	patterns := make(map[string]Pattern)
	patternsFileHandle, err := os.ReadFile(path)
	if err != nil {

		return patterns, err
	}
	err = json.Unmarshal(patternsFileHandle, &patterns)
	return patterns, err
}

func (self *PatternManager) loadPatterns() error {
	IPV4_REGEX := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}`
	FQDN_REGEX := `[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*\.[a-z]{2,}`
	patterns, err := readConfig(patternsFilePath)
	if err != nil || len(patterns) == 0 {
		patterns = make(map[string]Pattern)
		patterns["ipv4"] = Pattern{Name: "ipv4", Regex: IPV4_REGEX}
		patterns["fqdn"] = Pattern{Name: "fqdn", Regex: FQDN_REGEX}
		for _, pattern := range patterns {
			self.maskPatterns[pattern.Name] = pattern
		}
		fmt.Fprintf(os.Stderr, "Config file not found. Created new one in %q \n", patternsFilePath)
		self.SavePatterns()
	}
	for _, pattern := range patterns {
		self.maskPatterns[pattern.Name] = pattern
	}
	return err
}

func (self *PatternManager) loadFilters() error {
	PHP_REGEX := `\.php$`
	filters, err := readConfig(filtersFilePath)
	if err != nil || len(filters) == 0 {
		filters = make(map[string]Pattern)
		filters["php"] = Pattern{Name: "php", Regex: PHP_REGEX}
		for _, pattern := range filters {
			self.filterPatterns[pattern.Name] = pattern
		}
		fmt.Fprintf(os.Stderr, "Filter pattern file not found. Created new one in %q \n", filtersFilePath)
		self.SaveFilters()
	}
	for _, pattern := range filters {
		self.filterPatterns[pattern.Name] = pattern
	}
	return err
}
func (self *PatternManager) GetPatterns() []Pattern {
	maskPatterns := []Pattern{}
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
	err = os.WriteFile(patternsFilePath, valueMapJson, 0644)
	return err
}

func (self *PatternManager) SaveFilters() error {
	err := os.MkdirAll(configDir, 0755)
	valueMapJson, err := json.Marshal(self.filterPatterns)
	err = os.WriteFile(filtersFilePath, valueMapJson, 0644)
	return err
}

func (self *PatternManager) MapValuesToPatterns(rawLine string) ([]PatternMatch, error) {
	var matchList []PatternMatch
	var err error
	for _, pattern := range self.GetPatterns() {
		var regex *regexp.Regexp
		regex, err = regexp.Compile(pattern.Regex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to compile regex %q\n", pattern.Regex)
			self.RemovePatternByName(pattern.Name)
			continue
		}
		sensitive_values := regex.FindAllString(rawLine, -1)
		matchList = append(matchList, PatternMatch{MaskPattern: pattern, Matches: sensitive_values})

	}

	for _, filter := range self.filterPatterns {
		var filteredMatches []string
		filterRegex, err := regexp.Compile(filter.Regex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to compile regex %q\n", filter.Regex)
			self.RemoveFilterByName(filter.Name)
			continue
		}
		for i, patternMatch := range matchList {
			for _, match := range patternMatch.Matches {
				if filterRegex.MatchString(match) {
					continue
				}
				filteredMatches = append(filteredMatches, match)
			}
			matchList[i].Matches = filteredMatches
		}

	}
	return matchList, err
}

func (self *PatternManager) AddPattern(pattern Pattern) {
	self.maskPatterns[pattern.Name] = pattern
}
func (self *PatternManager) RemovePatternByName(name string) (Pattern, error) {
	var err error
	pattern, present := self.maskPatterns[name]
	if !present {
		err = errors.New("Pattern not found")
	}
	delete(self.maskPatterns, name)
	return pattern, err
}
func (self *PatternManager) RemoveFilterByName(name string) (Pattern, error) {
	var err error
	pattern, present := self.filterPatterns[name]
	if !present {
		err = errors.New("Filter pattern not found")
	}
	delete(self.filterPatterns, name)
	return pattern, err
}
