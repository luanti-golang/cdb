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

	found := repo.FindByName("repixture")
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
	found := offlineRepo(t).FindByName("realism_512")
	if len(found) != 1 {
		t.Fatalf("Find(realism_512) returned %d packages, want 1", len(found))
	}
	if found[0].Repo != nil {
		t.Errorf("Repo = %q, want nil", *found[0].Repo)
	}
}

func TestOfflineRepositoryUnknown(t *testing.T) {
	if got := offlineRepo(t).FindByName("does_not_exist_xyz"); len(got) != 0 {
		t.Errorf("Find(unknown) = %v, want empty", got)
	}
}

func TestOfflineRepositoryFindByKey(t *testing.T) {
	repo := offlineRepo(t)

	pkg, ok := repo.Find("Wuzzy/repixture")
	if !ok || pkg.Name != "repixture" {
		t.Fatalf("Find(Wuzzy/repixture) = %v, %v", pkg, ok)
	}
	if _, ok := repo.Find("repixture"); ok {
		t.Error("bare name must not be a key")
	}
	if _, ok := repo.Find("Nobody/repixture"); ok {
		t.Error("wrong author must not match")
	}
}

func TestOfflineRepositoryAliasLookup(t *testing.T) {
	repo := offlineRepo(t)

	// happy_land has the alias "a-zA0-02A/gamingoont_happy_land"
	var byName *VCSPackage
	for _, p := range repo.FindByName("happy_land") {
		byName = p
	}
	byAlias, ok := repo.Find("a-zA0-02A/gamingoont_happy_land")
	if byName == nil || !ok || byName != byAlias {
		t.Fatalf("alias lookup mismatch: %v vs %v (ok=%v)", byName, byAlias, ok)
	}
}

func TestOfflineRepositoryDuplicateNames(t *testing.T) {
	// several authors publish a package named "apartment"
	found := offlineRepo(t).FindByName("apartment")
	if len(found) < 2 {
		t.Fatalf("Find(apartment) returned %d packages, want >= 2", len(found))
	}
	for i, p := range found {
		if i > 0 && found[i-1].Author > p.Author {
			t.Error("results not ordered by author")
		}
		if got, ok := offlineRepo(t).Find(p.Author + "/" + p.Name); !ok || got.Author != p.Author {
			t.Errorf("Find(%s/%s) failed", p.Author, p.Name)
		}
	}
}

func TestOfflineRepositoryContainsAllPackages(t *testing.T) {
	var all []*VCSPackage
	if err := json.Unmarshal(offlinePackagesJSON, &all); err != nil {
		t.Fatal(err)
	}
	repo := offlineRepo(t)
	for _, want := range all {
		got, ok := repo.Find(want.Author + "/" + want.Name)
		if !ok || got.Name != want.Name {
			t.Errorf("package %s/%s missing from repository", want.Author, want.Name)
		}
	}
}

func TestRealKeysTakePrecedenceOverAliases(t *testing.T) {
	real := &VCSPackage{Author: "a", Name: "x"}
	other := &VCSPackage{Author: "b", Name: "y", Aliases: []string{"a/x"}}
	repo := newPackageRepository([]*VCSPackage{other, real})
	if got, _ := repo.Find("a/x"); got != real {
		t.Error("alias overrode a real package key")
	}
}

func TestNewPackageRepositorySkipsInvalid(t *testing.T) {
	repo := newPackageRepository([]*VCSPackage{nil, {Author: "u"}, {Author: "u", Name: "a"}})
	if _, ok := repo.Find("u/a"); !ok || len(repo.packages) != 1 {
		t.Errorf("unexpected packages: %v", repo.packages)
	}
}
