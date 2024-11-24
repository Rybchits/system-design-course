package commands

import (
	"os"
	"path/filepath"
	"shell/internal/command_meta"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeDir(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test_cd")
	require.NoError(t, err, "Could not create temp dir")
	defer os.RemoveAll(baseDir)

	innerDir, err := os.MkdirTemp(baseDir, "inner")
	require.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	cases := []struct {
		name         string
		argPath      string
		expectedPath string
		expectErr    bool
	}{
		{
			name:         "return home",
			argPath:      "",
			expectedPath: homeDir,
		},
		{
			name:         "absolute path",
			argPath:      innerDir,
			expectedPath: innerDir,
		},
		{
			name:         "relative path",
			argPath:      innerDir[len(baseDir)+1:],
			expectedPath: innerDir,
		},
		{
			name:      "non-existent path",
			argPath:   filepath.Join(baseDir, "foo", "bar"),
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Перейти в базовую директорию в начале теста
			require.NoError(t, os.Chdir(baseDir))

			args := make([]string, 0, 1)
			if tc.argPath != "" {
				args = append(args, tc.argPath)
			}
			meta := command_meta.CommandMeta{Name: "cd", Args: args}

			oldPwd, err := os.Getwd()
			require.NoError(t, err)

			cmd := ChangeDirCommand{meta}
			err = cmd.Execute()

			if tc.expectErr {
				require.Error(t, err)
			} else {
				newPwd, err := os.Getwd()
				require.NoError(t, err)
				require.NotEqual(t, oldPwd, newPwd, "pwd did not change after execution. Test is inconclusive")

				// os.MkdirTemp создает папку в /var,
				// что на macos на самом деле симлинк на /private/var
				// поэтому, чтобы тест честно проходил, надо здесь резолвить симлинк
				expectedDir, err := filepath.EvalSymlinks(tc.expectedPath)
				require.NoError(t, err)
				require.Equal(t, expectedDir, newPwd)
			}

		})
	}

}
