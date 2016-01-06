// This is the primary file of LibreJS-Gopher

package librejsgopher

import (
	"errors"
	"io/ioutil"
	"strings"
)

var LicensesCapitalizedStrings []string // An array of strings where each string is a license name or sub-string that needs to be capitalized

var LicenseMap map[string]string // LicenseMap is a map of license names to magnet URLs

func init() {
	LicensesCapitalizedStrings = []string{"BSD", "CC", "GPL", "ISC", "MPL"}

	LicenseMap = map[string]string{
		"AGPL-3.0":      "magnet:?xt=urn:btih:0b31508aeb0634b347b8270c7bee4d411b5d4109&dn=agpl-3.0.txt",
		"Apache-2.0":    "magnet:?xt=urn:btih:8e4f440f4c65981c5bf93c76d35135ba5064d8b7&dn=apache-2.0.txt",
		"Artistic-2.0":  "magnet:?xt=urn:btih:54fd2283f9dbdf29466d2df1a98bf8f65cafe314&dn=artistic-2.0.txt",
		"BSD-3.0":       "magnet:?xt=urn:btih:c80d50af7d3db9be66a4d0a86db0286e4fd33292&dn=bsd-3-clause.txt",
		"CC0":           "magnet:?xt=urn:btih:90dc5c0be029de84e523b9b3922520e79e0e6f08&dn=cc0.txt",
		"Expat":         "magnet:?xt=urn:btih:d3d9a9a6595521f9666a5e94cc830dab83b65699&dn=expat.txt",
		"FreeBSD":       "magnet:?xt=urn:btih:87f119ba0b429ba17a44b4bffcab33165ebdacc0&dn=freebsd.txt",
		"GPL-2.0":       "magnet:?xt=urn:btih:cf05388f2679ee054f2beb29a391d25f4e673ac3&dn=gpl-2.0.txt",
		"GPL-3.0":       "magnet:?xt=urn:btih:1f739d935676111cfff4b4693e3816e664797050&dn=gpl-3.0.txt",
		"ISC":           "magnet:?xt=urn:btih:b8999bbaf509c08d127678643c515b9ab0836bae&dn=ISC.txt",
		"LGPL-2.1":      "magnet:?xt=urn:btih:5de60da917303dbfad4f93fb1b985ced5a89eac2&dn=lgpl-2.1.txt",
		"LGPL-3.0":      "magnet:?xt=urn:btih:0ef1b8170b3b615170ff270def6427c317705f85&dn=lgpl-3.0.txt",
		"MPL-2.0":       "magnet:?xt=urn:btih:3877d6d54b3accd4bc32f8a48bf32ebc0901502a&dn=mpl-2.0.txt",
		"Public-Domain": "magnet:?xt=urn:btih:e95b018ef3580986a04669f1b5879592219e2a7a&dn=public-domain.txt",
		"X11":           "magnet:?xt=urn:btih:5305d91886084f776adcf57509a648432709a7c7&dn=x11.txt",
		"XFree86":       "magnet:?xt=urn:btih:12f2ec9e8de2a3b0002a33d518d6010cc8ab2ae9&dn=xfree86.txt",
	}
}

// GetFileLicense
// This function will get the license of the file, assuming it uses a valid LibreJS short-form header.
func GetFileLicense(file string) (LibreJSMetaInfo, error) {
	var getError error
	var metaInfo LibreJSMetaInfo

	fileContentBytes, fileReadError := ioutil.ReadFile(file) // Get the fileContent or if the file does not exist (or we do not have the permission) assign to fileReadError

	if fileReadError == nil { // If there was no read error
		fileContent := string(fileContentBytes[:])           // Convert to string
		fileContentLines := strings.Split(fileContent, "\n") // Split each new line into an []string
		fileContentLinesCount := len(fileContentLines)

		if fileContentLinesCount > 1 { // If this file is not a single line or empty
			fileLineParserChannel := make(chan LibreJSMetaInfo) // Make a channel that takes LibreJSMetaInfo
			linesParsed := 0                                    // Define linesParsed as the number of lines parsed by FileLicenseLineParser

			for _, lineContent := range fileContentLines { // For each license
				go FileLicenseLineParser(fileLineParserChannel, lineContent) // Asynchronously call FileLicenseLineParser
			}

		LineParserLoop:
			for libreJsMetaInfo := range fileLineParserChannel { // Constantly listen for channel input
				var endChannelListening bool

				linesParsed++ // Add one two linesParsed

				if libreJsMetaInfo.License != "" { // If the provided LibreJSMetaInfo has a valid License
					metaInfo = libreJsMetaInfo // Assign metaInfo as provided libreJsMetaInfo
					endChannelListening = true
				}

				if (fileContentLinesCount == linesParsed) || (endChannelListening) { // If we have parsed all lines or found the header info
					close(fileLineParserChannel) // Close the channel
					break LineParserLoop         // Break the loop
				}
			}

			if metaInfo.License == "" { // If there is no License defined by the end of the file
				getError = errors.New("LibreJS short-form header does not exist in this file.")
			}
		} else { // If the length of the file is 1 line or none
			getError = errors.New("File is either empty or does not contain the necessary individual lines required by LibreJS short-form blocks.")
		}
	} else { // If the file does not exist
		getError = errors.New(file + " does not exist.")
	}

	return metaInfo, getError
}

// FileLicenseLineParser
// This function handles individual line parsing
func FileLicenseLineParser(returnContentChannel chan LibreJSMetaInfo, lineContent string) {
	metaInfo := LibreJSMetaInfo{}

	lineContent = strings.Replace(lineContent, "//", "", -1) // Replace any // with nothing
	lineContent = strings.Replace(lineContent, "*", "", -1)  // Replace any * (block quotes) with nothing
	lineContent = strings.TrimPrefix(lineContent, " ")       // Trim any prefixed whitespace

	if strings.HasPrefix(lineContent, "@license") { // If the line starts with @license
		licenseHeaderFragments := strings.SplitN(lineContent, " ", 3)  // Split the license header info into three segments, separated by whitespace
		metaInfo.License = ParseLicenseName(licenseHeaderFragments[2]) // Define License as the parsed license name of the last item in fragments index
		metaInfo.Magnet = licenseHeaderFragments[1]                    // Define Magnet as the second item in the fragments index
	}

	returnContentChannel <- metaInfo
}

// GetMagnetLink
// This function will get a magnet link of the associated license exists
// Returns string for magnet link, error if item does not exist
func GetMagnetLink(license string) (string, error) {
	var magnetLinkFetchError error

	license = ParseLicenseName(license) // Parse the license name first
	magnetURL, licenseExists := LicenseMap[license]

	if !licenseExists { // If the license does not exist
		magnetLinkFetchError = errors.New(license + " does not exist.")
	}

	return magnetURL, magnetLinkFetchError
}

// ParseLicenseName
// This function will attempt to parse the provided license into a more logic naming scheme used in LicenseMap
func ParseLicenseName(license string) string {
	license = strings.ToLower(license) // Lowercase the entire string to make selective capitalization easier

	for _, licenseCapitalizedString := range LicensesCapitalizedStrings { // For each capitalized string of a license in LicensesCapitalizedStrings
		license = strings.Replace(license, strings.ToLower(licenseCapitalizedString), licenseCapitalizedString, -1) // Replace any lowercase instance with capitalized instance
	}

	license = strings.ToTitle(license)               // Title the license (example: apache -> Apache)
	license = strings.Replace(license, " ", "-", -1) // Replace whitespacing with hyphens

	return license
}
