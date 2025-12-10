package utils

import (
	"discountengine/internal/models"
	"encoding/json"
	"os"
	"sync"
)

// ConfigLoader handles loading of discount rules
type ConfigLoader struct {
	rules     []models.Rule
	rulesPath string
	mu        sync.RWMutex
}

// NewConfigLoaders creates a new ConfigLoaders
func NewConfigLoader(rulespath string) *ConfigLoader {
	return &ConfigLoader{
		rulesPath: rulespath,
	}
}

// LoadRules loads rules from JSON file
func (cl *ConfigLoader) LoadRules() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	file, err := os.ReadFile(cl.rulesPath)
	if err != nil {
		return err
	}

	var rules []models.Rule
	if err := json.Unmarshal(file, &rules); err != nil {
		return err
	}

	cl.rules = rules
	return nil
}

// GetRules return the current rules (thread-safe)
func (cl *ConfigLoader) GetRules() []models.Rule {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	// return a copy to prevent external modification
	rulesCopy := make([]models.Rule, len(cl.rules))
	copy(rulesCopy, cl.rules)
	return rulesCopy
}

// ReloadRules reloads rules from file
func (cl *ConfigLoader) ReloadRules() error {
	return cl.LoadRules()
}
