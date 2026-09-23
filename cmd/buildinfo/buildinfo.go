package buildinfo

import (
	"encoding/json"
)

const (
	versionDefaultBase = "0.0.0"
)

// GetVersion returns the version found in a JSON definition.
func GetVersion(definition []byte) (version string, err error) {
	partInfo := struct {
		Version string
	}{}
	if err := json.Unmarshal(definition, &partInfo); err != nil {
		return versionDefaultBase, err
	}
	if partInfo.Version == "" {
		return versionDefaultBase, err
	}
	return partInfo.Version, nil
}
