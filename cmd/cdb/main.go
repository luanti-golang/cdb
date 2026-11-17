package main

import (
	"cdb"
	"flag"
	"fmt"
	"os"
)

func main() {
	helpFlag := flag.Bool("help", false, "shows the help")
	installFlag := flag.Bool("install", false, "install a package")
	packageFlag := flag.String("package", "", "package name")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	var err error
	if *installFlag {
		err = cdb.Install(*packageFlag)
	} else {
		flag.Usage()
		return
	}
	if err != nil {
		fmt.Printf("Install failed: %v\n", err)
		os.Exit(1)
	}
}
