package cdb

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed offline-packages.json
var offlinePackagesJSON []byte

type PackageRepository struct {
	// map of package-name to possible package matches
	packages map[string][]*VCSPackage
}

// CreateRepository creates an online repository and falls back to the
// offline repository if that fails.
func CreateRepository() (*PackageRepository, error) {
	repo, err := CreateOnlineRepository()
	if err == nil {
		return repo, nil
	}

	repo, offlineErr := CreateOfflineRepository()
	if offlineErr != nil {
		return nil, fmt.Errorf("online repository failed (%v), offline repository failed: %w", err, offlineErr)
	}
	return repo, nil
}

// CreateOnlineRepository fetches all packages from the api.
func CreateOnlineRepository() (*PackageRepository, error) {
	packages, err := GetVCSPackages()
	if err != nil {
		return nil, err
	}
	return newPackageRepository(packages), nil
}

// CreateOfflineRepository creates a repository from the embedded offline-packages.json.
func CreateOfflineRepository() (*PackageRepository, error) {
	var packages []*VCSPackage
	if err := json.Unmarshal(offlinePackagesJSON, &packages); err != nil {
		return nil, fmt.Errorf("decode offline packages: %w", err)
	}
	return newPackageRepository(packages), nil
}

// newPackageRepository populates the name -> packages map. Packages are indexed
// by their name and by the name part of each alias ("author/name"), so lookups
// of renamed packages still match.
func newPackageRepository(packages []*VCSPackage) *PackageRepository {
	repo := &PackageRepository{packages: make(map[string][]*VCSPackage, len(packages))}
	for _, pkg := range packages {
		if pkg == nil {
			continue
		}
		repo.add(pkg.Name, pkg)
		for _, alias := range pkg.Aliases {
			name := alias[strings.LastIndex(alias, "/")+1:]
			repo.add(name, pkg)
		}
	}
	return repo
}

func (r *PackageRepository) add(name string, pkg *VCSPackage) {
	if name == "" {
		return
	}
	for _, existing := range r.packages[name] {
		if existing == pkg {
			return
		}
	}
	r.packages[name] = append(r.packages[name], pkg)
}

// Find returns all packages matching the given name.
func (r *PackageRepository) Find(name string) []*VCSPackage {
	return r.packages[name]
}
