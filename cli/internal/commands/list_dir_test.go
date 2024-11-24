package commands

import (
	"io"
	"os"
	"path/filepath"
	"shell/internal/command_meta"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func make_inner_files(t *testing.T, baseDir string, structure map[string]any) error {
	for entryName, v := range structure {
		entryFullPath := filepath.Join(baseDir, entryName)
		switch value := v.(type) {
		case bool:
			// file
			_, err := os.Create(entryFullPath)
			if err != nil {
				return err
			}
		case map[string]any:
			// directory
			err := os.Mkdir(entryFullPath, 0777)
			if err != nil {
				return err
			}

			err = make_inner_files(t, entryFullPath, value)
			if err != nil {
				return err
			}
		default:
			t.Fatal("Unexpected type in setup")
		}
	}

	return nil
}

func TestListDir(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test_ls")
	require.NoError(t, err, "Could not create temp dir")
	defer os.RemoveAll(baseDir)

	treeStructure := map[string]any{
		"file1": true,
		"file2": true,
		"dir": map[string]any{
			"ff1": true,
			"ff2": true,
			"dd":  map[string]any{},
		},
	}
	err = make_inner_files(t, baseDir, treeStructure)
	require.NoError(t, err)
	require.NoError(t, os.Chdir(baseDir))

	cases := []struct {
		name           string
		argPath        string
		expectedResult []string
		expectErr      bool
	}{
		{
			name:           "current folder",
			argPath:        "",
			expectedResult: []string{"file1", "file2", "dir"},
		},
		{
			name:           "local folder",
			argPath:        "dir",
			expectedResult: []string{"ff1", "ff2", "dd"},
		},
		{
			name:           "local path to empty",
			argPath:        "dir/dd",
			expectedResult: nil,
		},
		{
			name:           "absolute path to empty",
			argPath:        filepath.Join(baseDir, "dir", "dd"),
			expectedResult: nil,
		},
		{
			name:      "Non existing dir",
			argPath:   filepath.Join(baseDir, "foo", "bar"),
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := make([]string, 0, 1)
			if tc.argPath != "" {
				args = append(args, tc.argPath)
			}
			meta := command_meta.CommandMeta{Name: "ls", Args: args}

			rp, wp, err := os.Pipe()
			if err != nil {
				t.Fatal("Cant create pipe", err)
			}
			defer rp.Close()

			cmd := ListDirCommand{wp, meta}
			err = cmd.Execute()
			wp.Close()

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				actualBytes, err := io.ReadAll(rp)
				require.NoError(t, err)
				actual := strings.TrimSpace(string(actualBytes))

				var fileWithPermSlice []string
				if actual != "" {
					fileWithPermSlice = strings.Split(actual, "\n")
				}

				// Здесь проверим только сам состав файлов. Правильность прав проверим отдельным тестом
				fileSlice := make([]string, len(fileWithPermSlice))
				for i, val := range fileWithPermSlice {
					fileSlice[i] = strings.SplitN(val, " ", 2)[1]
				}

				require.ElementsMatch(t, tc.expectedResult, fileSlice)
			}

		})
	}
}

func TestListDirPermissionString(t *testing.T) {
	cases := []struct {
		fileMode os.FileMode
		expected string
	}{
		{0, "----------"},
		{os.ModeDir | 0755, "drwxr-xr-x"},
		{os.ModeDir | 0777, "drwxrwxrwx"},
		{0644, "-rw-r--r--"},
		{0600, "-rw-------"},
	}

	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			actual := permissionString(tc.fileMode)
			assert.Equal(t, tc.expected, actual, "Mode was %#o", tc.fileMode)
		})
	}

}
