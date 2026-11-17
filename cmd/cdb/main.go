package main

import "flag"

func main() {
	helpFlag := flag.Bool("help", false, "shows the help")
	installFlag := flag.Bool("install", false, "install a package")
	packageFlag := flag.String("package", "", "package name")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *installFlag {
	}
}
