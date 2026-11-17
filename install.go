package cdb

import "fmt"

func Install(name string) error {
	repo, err := CreateOfflineRepository()
	if err != nil {
		return fmt.Errorf("repo creation failed: %v", err)
	}

	pkgs := repo.FindByName(name)
	for _, pkg := range pkgs {
		fmt.Printf("Package: %v\n", pkg)
	}

	return nil
}
