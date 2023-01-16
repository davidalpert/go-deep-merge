// +build ignore

// This program generates version.go. It can be invoked by running invoking go:generate
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"
)

func main() {
	versionStem, err := exec.Command("sbot", "get", "version").Output()
	if err != nil {
		fmt.Printf("can't read version from git tags: %s", err)
		fmt.Printf("defaulting to 0.0.1")
		versionStem = []byte("0.0.1")
	}

	if len(os.Args) < 2 {
		fmt.Printf("usage: go run .tools/version_gen.go <appName>\n\nrun from project root")
		os.Exit(1)
	}

	gitBranch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	//fmt.Printf("branch: >%s< (err: %#v)\n", string(gitBranch), err)
	if err != nil {
		panic(err)
	}

	gitSHA, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	//fmt.Printf("sha: >%s< (err: %#v)\n", string(gitSHA), err)
	if err != nil {
		panic(err)
	}

	var hasStaged bool
	if _, err := exec.Command("git", "diff-index", "--quiet", "--cached", "HEAD", "--").Output(); err != nil {
		hasStaged = true
	}

	var hasModified bool
	if _, err := exec.Command("git", "diff-files", "--quiet").Output(); err != nil {
		hasModified = true
	}

	var hasUntracked bool
	if _, err := exec.Command("git", "ls-files", "--exclude-standard", "--others").Output(); err != nil {
		hasUntracked = true
	}

	var isDirty = hasStaged || hasModified || hasUntracked

	_, err = exec.Command("mkdir", "-p", "./internal/version").Output()
	if err != nil {
		panic(err)
	}

	fmt.Printf("generating/updating: ./internal/version/detail.go\n")
	f, err := os.Create("internal/version/detail.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	packageTemplate.Execute(f, struct {
		AppName      string
		Version      string
		Branch       string
		SHA          string
		HasStaged    bool
		HasModified  bool
		HasUntracked bool
		IsDirty      bool
		Timestamp    time.Time
	}{
		AppName:      os.Args[1],
		Version:      string(bytes.Trim(versionStem, "\r\n")),
		Branch:       string(bytes.Trim(gitBranch, "\r\n")),
		SHA:          string(bytes.Trim(gitSHA, "\r\n")),
		HasStaged:    hasStaged,
		HasModified:  hasModified,
		HasUntracked: hasUntracked,
		IsDirty:      isDirty,
		Timestamp:    time.Now(),
	})
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at {{ .Timestamp }}
package version

import (
	"fmt"
	"os/user"
	"runtime"
	"sort"
	"strings"
)

// DetailStruct provides an easy way to grab all the govvv version details together
type DetailStruct struct {
	AppName              string ` + "`json:\"app_name\"`" + `
	BuildDate            string ` + "`json:\"build_date\"`" + `
	GitBranch            string ` + "`json:\"branch\"`" + `
	GitCommit            string ` + "`json:\"commit\"`" + `
	GitDirty             bool ` + "`json:\"dirty\"`" + `
	GitDirtyHasModified  bool ` + "`json:\"dirty_modified\"`" + `
	GitDirtyHasStaged    bool ` + "`json:\"dirty_staged\"`" + `
	GitDirtyHasUntracked bool ` + "`json:\"dirty_untracked\"`" + `
	GitWorkingState      string ` + "`json:\"working_state\"`" + `
	GitSummary           string ` + "`json:\"summary\"`" + `
	UserAgentString      string ` + "`json:\"user_agent\"`" + `
	Version              string ` + "`json:\"version\"`" + `
}

// NewVersionDetail builds a new version DetailStruct
func NewVersionDetail() DetailStruct {
	s := DetailStruct{
		AppName:              "{{ .AppName }}",
		BuildDate:            "{{ .Timestamp }}",
		GitBranch:            "{{ .Branch }}",
		GitCommit:            "{{ .SHA }}",
		GitDirty:             {{ .IsDirty }},
		GitDirtyHasModified:  {{ .HasModified }},
		GitDirtyHasStaged:    {{ .HasStaged }},
		GitDirtyHasUntracked: {{ .HasUntracked }},
		GitSummary:           "{{ .Timestamp }}",
		GitWorkingState:      "",
		Version:              "{{ .Version }}",
	}
	s.UserAgentString = s.ToUserAgentString()
	if s.GitDirty {
		s.GitWorkingState = "dirty"
	}
	return s
}

// Detail provides an easy global way to
var Detail = NewVersionDetail()

// ToUserAgentString formats a DetailStruct as a User-Agent string
func (s DetailStruct) ToUserAgentString() string {
	productName := s.AppName
	productVersion := s.Version

	productDetails := map[string]string{
		"sha": s.GitCommit,
	}

	if s.GitBranch != "main" {
		productDetails["branch"] = s.GitBranch
	}

	if s.GitDirty {
		productDetails["dirty"] = "true"
	}

	user, err := user.Current()
	if err == nil {
		username := user.Username
		if username == "" {
			username = "unknown"
		}

		productDetails["user"] = username // strings.Replace(user.Username, "a-", 1) // this is a northfield convention
	}

	detailParts := []string{}
	for k, v := range productDetails {
		detailParts = append(detailParts, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Slice(detailParts, func(i, j int) bool {
		return detailParts[i] < detailParts[j]
	})
	productDetail := strings.Join(detailParts, ", ")

	return fmt.Sprintf("%s/%s (%s) %s (%s)", productName, productVersion, productDetail, runtime.GOOS, runtime.GOARCH)
}
`))
