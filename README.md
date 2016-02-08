# LibreJS-Gopher #

LibreJS-Gopher is a Go package (as well as offering binary) that handles:

1. The management and organization of OSI-approved and LibreJS licenses and their known (if any) magnet links for usage with LibreJS.
2. The checking of licensing information of files.
3. The wrapping of JavaScript files with LibreJS compliant license header and footers.

**Current Stable Release:** `0.2`

## Contributing

This project leverages CodeUtils for development and adopts the CodeUtils Usage Spec. To learn how to contribute to this project and set up CodeUtils, read [here](https://github.com/StroblIndustries/CodeUtils/blob/master/CodeUtils-Usage-Spec.md).

## Binary ##

### Downloading ###

You can download the pre-compiled binary via the latest stable release [here](https://github.com/JoshStrobl/librejs-gopher/releases/tag/0.2), it is available in the "binary" folder of the tarball.

### Compiling ###

**Requirements:**

- Download and build [nflag](https://github.com/JoshStrobl/nflag).

**Instructions:**

Go to `binary` folder and run the following: `go build -o libregopher libregopher.go`

### Usage ###

You can see the usage of the binary by running:

`./libregopher` or `./libregopher --help`

## Package ##

### Import ###

You can import and use this software by using the line:

``` go
import "github.com/JoshStrobl/librejs-gopher"
```

Don't have it? Use:

```
go get github.com/JoshStrobl/librejs-gopher
```

### Usage ###

#### Variables ####

The following variables are exposed and *useful* when leveraging LibreJS-Gopher.

Name | Description | Type
----- | ----- | -----
LicenseMap | LicenseMap is a map of license names to magnet URLs | `map[string]string`

#### Structs ####

The following structs are used and exposed throughout the package

**LibreJSMetaInfo:** This structure contains the name and associated magnet link.


``` go
type LibreJSMetaInfo struct {
    License string
    Magnet  string
}
```

#### Funcs ####

The following functions are available for usage by LibreJS-Gopher.

##### AddLicense #####

This function will add a valid LibreJS short-form header and footer to the file. You can set to write the file automatically. We will always return new file content or an error.

``` go
func AddLicenseInfo(license string, file string, writeContentAutomatically bool) (string, error)
```

**Notes:**
1. `license` must correspond with a valid license exposed by LicenseMap.
2. `license` is automatically parsed using `ParseLicenseName()`

**Example:**

``` go
// --- Faux file content ---
// potato potato potato

newFileContent, addError := AddLicenseInfo("Apache-2.0", "potato.js", true)

// --- New File content ---

// @license magnet:?xt=urn:btih:8e4f440f4c65981c5bf93c76d35135ba5064d8b7&dn=apache-2.0.txt Apache-2.0
// potato potato potato
// @license-end
```

##### AddLicenseInfo #####

This function is a backwards-compatible function. Originally landing in `0.1`, this function was responsible for adding the license to the JavaScript file. That function has since been renamed to AddLicense and thus this function solely calls AddLicense. This function will be deprecated in `0.3`.

For documentation, refer to the AddLicense documentation.

##### GetFileLicense #####

This function will get the license of the file, assuming it uses a valid LibreJS [short-form header](http://www.gnu.org/software/librejs/free-your-javascript.html#magnet-link-license).

``` go
func GetFileLicense(file string) (LibreJSMetaInfo, error)
```

**Example:**

``` go
// --- Faux file content ---

// @license magnet:?xt=urn:btih:0b31508aeb0634b347b8270c7bee4d411b5d4109&dn=agpl-3.0.txt AGPL-3.0
// blah blah JS blah
// @license-end

// -- Go Call --

metaInfo, error := GetFileLicense("path/to/file.min.js")

// -- metaInfo --
// type LibreJsMetaInfo struct {
// License : AGPL-3.0
// Magnet : magnet:?xt=urn:btih:0b31508aeb0634b347b8270c7bee4d411b5d4109&dn=agpl-3.0.txt
// }
```

##### GetMagnetLink #####

This function will get a magnet link of the associated license exists. Returns string for magnet link, error if item does not exist.

``` go
func GetMagnetLink(license string) (string, error)
```

**Notes:**
1. `license` is automatically parsed using `ParseLicenseName()`

**Example:**

``` go
// Returns string: magnet:?xt=urn:btih:8e4f440f4c65981c5bf93c76d35135ba5064d8b7&dn=apache-2.0.txt
// Returns error: nil
magnetURL, magnetGetError := librejsgopher.GetMagnetLink("Apache-2.0")
```

##### ParseLicenseName #####

This function will attempt to parse the provided license into a more logic naming scheme used in LicenseMap.

``` go
func ParseLicenseName(license string) string
```

**Example:**

``` go
// Returns LGPL-2.1
parsedLicenseName := ParseLicenseName("lgpl 2.1")
```