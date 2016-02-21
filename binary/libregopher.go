// LibreGopher is a CLI companion to librejs-gopher (which it implements).

package main

import (
	"errors"
	"fmt"
	"github.com/JoshStrobl/librejs-gopher"
	"github.com/JoshStrobl/nflag"
	"sort"
	"strings"
)

// LibreGopher initialization
func init() {
	nflag.Configure(nflag.ConfigOptions{
		ProgramDescription: "LibreGopher is a CLI companion to the Golang package librejs-gopher.",
		ShowHelpIfNoArgs:   true,
	})

	// Add LibreJS header info
	nflag.Set("add", nflag.Flag{
		Descriptor:   "Setting add to a license will cause it to be added to the file.",
		Type:         "string",
		Required:     false,
		AllowNothing: false,
	})

	// File flag
	nflag.Set("file", nflag.Flag{
		Descriptor:   "File we are adding LibreJS header to or getting the licensing information from.",
		Type:         "string",
		Required:     false,
		AllowNothing: false, // Must pass something when using this flag
	})

	// List all valid LibreJS licenses and cooresponding magnet links
	nflag.Set("list-licenses", nflag.Flag{
		Descriptor:   "Gets all the valid LibreJS licenses and cooresponding magnet links.",
		Type:         "bool",
		DefaultValue: false,
		AllowNothing: true,
	})

	// Get File Info flag
	nflag.Set("get-info", nflag.Flag{
		Descriptor:   "Whether we should get the LibreJS header information of this file.",
		Type:         "bool",
		DefaultValue: false,
		Required:     false,
		AllowNothing: true, // Not required to pass value, assumed to be true
	})
}

// SortedLicenses will get the keys of the librejsgopher LicenseMap, ensure they are sorted, then proceed to return them
func SortedLicenses() []string {
	keys := []string{} // Make the string slice

	for licenseName := range librejsgopher.LicenseMap { // For each licenseName
		keys = append(keys, []string{licenseName}...) // Append the licenseName to keys
	}

	sort.Strings(keys) // Sort the keys
	return keys
}

func main() {
	nflag.Parse()

	addLicenseString, _ := nflag.GetAsString("add")     // Get license name provided
	fileProvided, _ := nflag.GetAsString("file")        // Get file name provided
	listLicenses, _ := nflag.GetAsBool("list-licenses") // Get a boolean value as to whether or not we should list licenses
	getInfo, _ := nflag.GetAsBool("get-info")           //

	var derp error // Define derp as the error struct we'll add errors from commands or constructed errors to

	if listLicenses { // If we are listing the licenses available
		fmt.Println("Below are the flags that LibreJS Gopher expose:")

		sortedLicenseNames := SortedLicenses()

		for _, licenseName := range sortedLicenseNames { // For each licenseName and magnetURL in LicenseMap
			magnetURL, _ := librejsgopher.GetMagnetLink(licenseName)             // Get the Magnet URL of the license
			licenseName = licenseName + strings.Repeat(" ", 15-len(licenseName)) // Change licenseName to ensure the width of each license name is 10char long
			fmt.Println(licenseName + magnetURL)
		}
	} else if fileProvided != "" { // If a file was passed
		if addLicenseString != "" { // If we are adding license information to a file
			_, derp = librejsgopher.AddLicense(addLicenseString, fileProvided, true) // Automatically write the license details to the file and return any error
		} else if getInfo { // If we are getting the information of this file
			var metaInfo librejsgopher.LibreJSMetaInfo                  // Define metaInfo as a LibreJsMetaInfo structure
			metaInfo, derp = librejsgopher.GetFileLicense(fileProvided) // Get any license information provided by GetFileLicense

			if derp == nil { // If there was no error
				fmt.Println("License: " + metaInfo.License)
				fmt.Println("Magnet Link: " + metaInfo.Magnet)
			}
		} else { // If we're not really sure what we're supposed to be doing, since no add flag or get-info flag was provided
			derp = errors.New("You need to pass the " + nflag.Config.OSSpecificFlagString + "add or " + nflag.Config.OSSpecificFlagString + "get-info flag when using the file flag.")
		}
	} else {
		nflag.PrintFlags() // Print the flags
	}

	if derp != nil { // If there was an error
		fmt.Printf("%v", derp) // Print the error
	}
}
