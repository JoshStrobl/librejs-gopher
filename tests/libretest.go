package main

import (
	"fmt"
	"github.com/JoshStrobl/librejs-gopher"
	"io/ioutil"
	"os"
)

func main() {
	// Bootstrap test
	fmt.Println("Bootstrapping test.")
	exampleJSContent := "console.log('hello world!');"              // Create an example JS file
	ioutil.WriteFile("test.min.js", []byte(exampleJSContent), 0744) // Write the file

	fmt.Println("") // New Line

	// Test Parsing License Name

	fmt.Println("Testing ParseLicenseName.")
	parsedApacheLicenseName := librejsgopher.ParseLicenseName("apache 2.0")
	fmt.Println("apache 2.0 parsed as: " + parsedApacheLicenseName)

	fmt.Println("") // New Line

	// Test GetMagnetLink on Apache 2.0

	fmt.Println("Testing GetMagnetLink.")
	magnetLink, _ := librejsgopher.GetMagnetLink(parsedApacheLicenseName)
	fmt.Println("Magnet Link for " + parsedApacheLicenseName + " is: " + magnetLink)

	fmt.Println("") // New Line

	// Test GetMagnetLink on invalid license

	fmt.Println("Testing GetMagnetLink on invalid license.")
	_, magnetError := librejsgopher.GetMagnetLink("Invalid License")
	fmt.Println(magnetError) // Print the error

	fmt.Println("") // New Line

	// Test Adding License Info

	fmt.Println("Testing AddLicenseInfo on example JS file.")
	newFileContent, _ := librejsgopher.AddLicenseInfo("Apache-2.0", "test.min.js", true) // Automatically write the file
	fmt.Println(newFileContent)

	fmt.Println("") // New Line

	// Test GetFileLicense on existant file

	fmt.Println("Testing GetFileLicense on example JS file.")
	libreJsMetaInfo, _ := librejsgopher.GetFileLicense("test.min.js")
	fmt.Printf("LibreJsMetaInfo struct: %v\n", libreJsMetaInfo)
	fmt.Println("License: " + libreJsMetaInfo.License)
	fmt.Println("Magnet Link: " + libreJsMetaInfo.Magnet)

	fmt.Println("") // New Line

	// Test GetFileLicense on non-existant file

	fmt.Println("Testing GetFileLicense on non-existant file.")
	_, getFileLicenseError := librejsgopher.GetFileLicense("nonexistantfile.min.js")
	fmt.Println(getFileLicenseError) // Print the error

	fmt.Println("") // New Line

	// Destroy test

	fmt.Println("Breaking down test.")
	os.Remove("test.min.js")
}
