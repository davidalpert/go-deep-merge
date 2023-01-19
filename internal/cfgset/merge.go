package cfgset

import (
	"fmt"
	"github.com/davidalpert/go-deep-merge/internal/app"
	"github.com/davidalpert/go-deep-merge/v1"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"path"
	"sort"
	"strings"
)

type MergeOptions struct {
	SourceFolder string
	Debug        bool
}

type MergeResult struct {
	AppDir      string
	MergeBySlug map[string]map[string]interface{}
}

// Merge merges a source folder of config files grouped by app
// assuming that each app folder contains a default.yaml and one or more
// slug.yaml (e.g. dev.yaml, prd.yaml, etc)
func Merge(o MergeOptions) ([]MergeResult, error) {
	fis, err := afero.ReadDir(app.Fs, o.SourceFolder)
	if err != nil {
		return nil, fmt.Errorf("reading source folder %#v: %#v", o.SourceFolder, err)
	}

	appDirs := make([]string, 0)
	for _, fi := range fis {
		if fi.IsDir() {
			appDirs = append(appDirs, path.Join(o.SourceFolder, fi.Name()))
		}
	}

	result := make([]MergeResult, 0)
	for _, appDir := range appDirs {
		var defaultFile string
		var overrideFiles = make([]string, 0)

		fis, err = afero.ReadDir(app.Fs, appDir)
		for _, fi := range fis {
			name := fi.Name()
			if strings.HasSuffix(name, "default.yaml") {
				defaultFile = path.Join(appDir, name)
			} else if strings.HasSuffix(name, ".yaml") {
				overrideFiles = append(overrideFiles, path.Join(appDir, name))
			}
		}

		sort.Slice(overrideFiles, func(i, j int) bool {
			return len(overrideFiles[i]) < len(overrideFiles[j])
		})

		destFile, err := afero.ReadFile(app.Fs, defaultFile)
		if err != nil {
			return nil, fmt.Errorf("read dest file %#v: %#v", defaultFile, err)
		}

		mergeResultBySlug := make(map[string]map[string]interface{})
		for _, override := range overrideFiles {
			var dest map[string]interface{}
			if err := yaml.Unmarshal(destFile, &dest); err != nil {
				return nil, fmt.Errorf("unmarshalling dest: %#v", err)
			}

			slug := strings.TrimSuffix(path.Base(override), path.Ext(override))
			if strings.ContainsAny(slug, ".") {
				baseSlug := strings.Split(slug, ".")[0]
				// merge on top of another

				r, err := v1.MergeWithOptions(mergeResultBySlug[baseSlug], dest, v1.NewConfigDeeperMergeBang().WithMergeHashArrays(true).WithDebug(o.Debug))
				if err != nil {
					return nil, fmt.Errorf("merging files %#v -> %#v: %#v", override, defaultFile, err)
				}

				dest = r
			}

			sourceFile, err := afero.ReadFile(app.Fs, override)
			if err != nil {
				return nil, fmt.Errorf("read source file %#v: %#v", override, err)
			}

			var src map[string]interface{}
			if err := yaml.Unmarshal(sourceFile, &src); err != nil {
				return nil, fmt.Errorf("unmarshalling src: %#v", err)
			}

			r, err := v1.MergeWithOptions(src, dest, v1.NewConfigDeeperMergeBang().WithMergeHashArrays(true).WithDebug(o.Debug))
			if err != nil {
				return nil, fmt.Errorf("merging files %#v -> %#v: %#v", override, defaultFile, err)
			}

			mergeResultBySlug[slug] = r
		}

		result = append(result, MergeResult{
			path.Base(appDir),
			mergeResultBySlug,
		})
	}
	return result, nil
}
