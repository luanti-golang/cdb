package cdb

import (
	"encoding/json"
	"testing"
)

func offlineRepo(t *testing.T) *PackageRepository {
	t.Helper()
	repo, err := CreateOfflineRepository()
	if err != nil {
		t.Fatalf("CreateOfflineRepository: %v", err)
	}
	return repo
}

func TestOfflineRepositoryNotEmpty(t *testing.T) {
	if len(offlineRepo(t).packages) == 0 {
		t.Fatal("offline repository is empty")
	}
}

func TestOfflineRepositoryFind(t *testing.T) {
	repo := offlineRepo(t)

	found := repo.Find("repixture")
	if len(found) != 1 {
		t.Fatalf("Find(repixture) returned %d packages, want 1", len(found))
	}
	pkg := found[0]
	if pkg.Author != "Wuzzy" || pkg.Type != "game" || pkg.Title != "Repixture" || !pkg.Featured {
		t.Errorf("unexpected package: %+v", pkg)
	}
	if pkg.Repo == nil || *pkg.Repo != "https://codeberg.org/Wuzzy/Repixture" {
		t.Errorf("unexpected repo: %v", pkg.Repo)
	}
}

func TestOfflineRepositoryNullRepo(t *testing.T) {
	found := offlineRepo(t).Find("realism_512")
	if len(found) != 1 {
		t.Fatalf("Find(realism_512) returned %d packages, want 1", len(found))
	}
	if found[0].Repo != nil {
		t.Errorf("Repo = %q, want nil", *found[0].Repo)
	}
}

func TestOfflineRepositoryUnknown(t *testing.T) {
	if got := offlineRepo(t).Find("does_not_exist_xyz"); len(got) != 0 {
		t.Errorf("Find(unknown) = %v, want empty", got)
	}
}

func TestOfflineRepositoryAliasLookup(t *testing.T) {
	repo := offlineRepo(t)

	// happy_land has the alias "a-zA0-02A/gamingoont_happy_land"
	byName := repo.Find("happy_land")
	byAlias := repo.Find("gamingoont_happy_land")
	if len(byName) != 1 || len(byAlias) != 1 {
		t.Fatalf("got %d by name, %d by alias, want 1 each", len(byName), len(byAlias))
	}
	if byName[0] != byAlias[0] {
		t.Error("name and alias lookups returned different packages")
	}
}

func TestOfflineRepositoryDuplicateNames(t *testing.T) {
	// several packages share the name "apartment"
	found := offlineRepo(t).Find("apartment")
	if len(found) < 2 {
		t.Fatalf("Find(apartment) returned %d packages, want >= 2", len(found))
	}
	seen := map[*VCSPackage]bool{}
	for _, p := range found {
		if seen[p] {
			t.Error("package listed twice")
		}
		seen[p] = true
	}
}

func TestOfflineRepositoryContainsAllPackages(t *testing.T) {
	var all []*VCSPackage
	if err := json.Unmarshal(offlinePackagesJSON, &all); err != nil {
		t.Fatal(err)
	}
	repo := offlineRepo(t)
	for _, want := range all {
		ok := false
		for _, got := range repo.Find(want.Name) {
			if got.Author == want.Author && got.Name == want.Name {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("package %s/%s missing from repository", want.Author, want.Name)
		}
	}
}

func TestNewPackageRepositorySkipsInvalid(t *testing.T) {
	repo := newPackageRepository([]*VCSPackage{nil, {Name: ""}, {Name: "a"}})
	if len(repo.packages) != 1 || len(repo.Find("a")) != 1 {
		t.Errorf("unexpected packages: %v", repo.packages)
	}
}
