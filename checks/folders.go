package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/math2001/gocmt/cmt"
)

func Folders(
	c *cmt.Check,
	args map[string]interface{},
) {
	path := args["path"].(string)

	var name string = path
	if v, ok := args["name"]; ok {
		name = v.(string)
	}

	recursive := args["recursive"].(bool) // returns false by deafult

	targets := args["target"].(map[interface{}]interface{})

	c.AddItem(&cmt.CheckItem{
		Name:  "folder_path",
		Value: path,
	})
	c.AddItem(&cmt.CheckItem{
		Name:  "folder_name",
		Value: name,
	})

	root, err := os.Stat(path)
	if err != nil {
		c.AddItem(&cmt.CheckItem{
			Name:         "folder_status",
			Value:        "nok",
			IsAlert:      true,
			AlertMessage: fmt.Sprintf("check_folder - %s missing", path),
		})
	}

	var filesCount int64
	var dirsCount int64
	var totalSize int64
	var minTime time.Time
	var maxTime time.Time
	var filenames []string

	// walk doesn't follow symlinks ($ go doc filepath.Walk)
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if os.SameFile(info, root) {
			return nil
		}

		// fmt.Fprintf(c.DebugBuffer(), "filepath: %s\n", path)
		// we actuall ignore the errors

		totalSize += info.Size()

		if maxTime.IsZero() || info.ModTime().After(maxTime) {
			maxTime = info.ModTime()
		}

		if minTime.IsZero() || info.ModTime().Before(minTime) {
			minTime = info.ModTime()
		}

		if info.IsDir() {
			dirsCount++
			if recursive {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		filenames = append(filenames, filepath.Base(path))

		filesCount++
		if !info.Mode().IsRegular() {
			return fmt.Errorf("found non regular file: %s", path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	c.AddItem(&cmt.CheckItem{
		Name:        "folder_files",
		Value:       filesCount,
		Description: "number of files",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "folder_dirs",
		Value:       dirsCount,
		Description: "number of directories",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "folder_size",
		Value:       totalSize,
		Description: "Total size (in bytes)",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "folder_age_min",
		Value:       time.Since(maxTime),
		Description: "Min age (seconds)",
		Unit:        "sec",
	})
	c.AddItem(&cmt.CheckItem{
		Name:        "folder_age_max",
		Value:       time.Since(minTime),
		Description: "Max age (seconds)",
		Unit:        "sec",
	})

	ci := &cmt.CheckItem{
		Name:    "folder_status",
		Value:   "nok",
		IsAlert: true,
	}
	defer c.AddItem(ci)

	if min, ok := targets["files_min"]; ok && filesCount < int64(min.(int)) {
		ci.AlertMessage = fmt.Sprintf("check_folder %q: too few files (%d)", path, filesCount)
		return
	}

	if max, ok := targets["files_max"]; ok && filesCount > int64(max.(int)) {
		ci.AlertMessage = fmt.Sprintf("check_folder %q: too many files (%d)", path, filesCount)
		return
	}

	if min, ok := targets["size_min"]; ok && totalSize < int64(min.(int)) {
		ci.AlertMessage = fmt.Sprintf("check_folder %q: too small (%d)", path, totalSize)
		return
	}

	if max, ok := targets["size_max"]; ok && totalSize > int64(max.(int)) {
		ci.AlertMessage = fmt.Sprintf("check_folder %q: too big (%d)", path, totalSize)
		return
	}

	if max, ok := targets["age_max"]; ok && !maxTime.IsZero() && int64(time.Since(maxTime).Seconds()) > int64(max.(int)) {
		ci.AlertMessage = fmt.Sprintf("check_folder %q: some files too old (%.0f sec)", path, time.Since(maxTime).Seconds())
		return
	}

	if min, ok := targets["age_min"]; ok && !minTime.IsZero() && int64(time.Since(minTime).Seconds()) < int64(min.(int)) {
		ci.AlertMessage = fmt.Sprintf("check_folder %q: some files too old (%.0f sec)", path, time.Since(minTime).Seconds())
		return
	}

	if v, ok := targets["has_files"]; ok {
		for _, v := range v.([]interface{}) {
			requiredFile := v.(string)
			var found bool
			for _, actualFile := range filenames {
				if requiredFile == actualFile {
					found = true
				}
			}
			if !found {
				ci.AlertMessage = fmt.Sprintf("check_folder - %s : expected file not found (%s)", path, requiredFile)
				return
			}
		}
	}

	ci.Value = "ok"
	ci.IsAlert = false
}
