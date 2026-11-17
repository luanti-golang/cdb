package cdb

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
)

//go:embed offline-packages.json
var offlinePackagesJSON []byte

type PackageRepository struct {
	// map of "{author}/{name}" to package
	packages map[string]*VCSPackage
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

// newPackageRepository populates the "{author}/{name}" -> package map. Aliases
// ("author/name") are indexed as additional keys, so renamed packages still
// resolve. Real package keys always take precedence over aliases.
func newPackageRepository(packages []*VCSPackage) *PackageRepository {
	repo := &PackageRepository{packages: make(map[string]*VCSPackage, len(packages))}
	for _, pkg := range packages {
		if pkg != nil && pkg.Name != "" {
			repo.add(packageKey(pkg.Author, pkg.Name), pkg)
		}
	}
	for _, pkg := range packages {
		if pkg == nil {
			continue
		}
		for _, alias := range pkg.Aliases {
			repo.add(alias, pkg)
		}
	}
	return repo
}

func packageKey(author, name string) string {
	return author + "/" + name
}

// add registers pkg under key unless the key is already taken.
func (r *PackageRepository) add(key string, pkg *VCSPackage) {
	if _, exists := r.packages[key]; !exists {
		r.packages[key] = pkg
	}
}

// Find returns the package with the given "{author}/{name}" key (or alias).
func (r *PackageRepository) Find(key string) (*VCSPackage, bool) {
	pkg, ok := r.packages[key]
	return pkg, ok
}

// FindByName returns all packages with the given name, regardless of author,
// ordered by author.
func (r *PackageRepository) FindByName(name string) []*VCSPackage {
	var found []*VCSPackage
	seen := map[*VCSPackage]bool{}
	for _, pkg := range r.packages {
		if pkg.Name == name && !seen[pkg] {
			seen[pkg] = true
			found = append(found, pkg)
		}
	}
	sort.Slice(found, func(i, j int) bool { return found[i].Author < found[j].Author })
	return found
}
