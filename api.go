package cdb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const vcsPackagesURL = "https://content.luanti.org/api/packages/?fmt=vcs"

// GetVCSPackages fetches the list of packages with their VCS information.
func GetVCSPackages() ([]*VCSPackage, error) {
	resp, err := http.Get(vcsPackagesURL)
	if err != nil {
		return nil, fmt.Errorf("request packages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request packages: unexpected status %s", resp.Status)
	}

	var packages []*VCSPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, fmt.Errorf("decode packages: %w", err)
	}
	return packages, nil
}
