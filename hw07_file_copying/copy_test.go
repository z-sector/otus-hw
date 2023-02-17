package main

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyError(t *testing.T) {
	tempFile, err := os.CreateTemp(os.TempDir(), "test_copy_error_*")
	require.NoError(t, err)
	defer func() {
		errR := os.Remove(tempFile.Name())
		require.NoError(t, errR)
	}()

	tests := []struct {
		name          string
		from          string
		to            string
		offset        int64
		limit         int64
		expectedError error
	}{
		{
			name:          "from is empty",
			from:          "",
			to:            tempFile.Name(),
			expectedError: ErrFromPathIsEmpty,
		},
		{
			name:          "to is empty",
			from:          "testdata/input100.txt",
			to:            "",
			expectedError: ErrToPathIsEmpty,
		},
		{
			name:          "offset is not positive",
			from:          "testdata/input100.txt",
			to:            tempFile.Name(),
			offset:        -1,
			expectedError: ErrNegativeOffset,
		},
		{
			name:          "limit is not positive",
			from:          "testdata/input100.txt",
			to:            tempFile.Name(),
			limit:         -1,
			expectedError: ErrNegativeLimit,
		},
		{
			name:          "unsupported file",
			from:          "/dev/urandom",
			to:            tempFile.Name(),
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "offset exceeds file size",
			from:          "testdata/input100.txt",
			to:            tempFile.Name(),
			offset:        101,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			name:          "file is a directory",
			from:          "testdata",
			to:            tempFile.Name(),
			expectedError: ErrFileIsDirectory,
		},
		{
			name:          "file does not exist",
			from:          "testdata/not_exists.txt",
			to:            tempFile.Name(),
			expectedError: fs.ErrNotExist,
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}

	err = tempFile.Close()
	require.NoError(t, err)
}

func TestCopySuccess(t *testing.T) {
	tests := []struct {
		offset           int64
		limit            int64
		expectedFileName string
	}{
		{
			offset:           0,
			limit:            0,
			expectedFileName: "out_offset0_limit0",
		},
		{
			offset:           0,
			limit:            10,
			expectedFileName: "out_offset0_limit10",
		},
		{
			offset:           0,
			limit:            1000,
			expectedFileName: "out_offset0_limit1000",
		},
		{
			offset:           0,
			limit:            10000,
			expectedFileName: "out_offset0_limit10000",
		},
		{
			offset:           100,
			limit:            1000,
			expectedFileName: "out_offset100_limit1000",
		},
		{
			offset:           6000,
			limit:            1000,
			expectedFileName: "out_offset6000_limit1000",
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.expectedFileName, func(t *testing.T) {
			temp, err := os.CreateTemp(os.TempDir(), "test_copy_*")
			require.NoError(t, err)
			defer func() {
				err = os.Remove(temp.Name())
				require.NoError(t, err)
			}()

			err = temp.Close()
			require.NoError(t, err)

			err = Copy("testdata/input.txt", temp.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			actualBytes, err := os.ReadFile(temp.Name())
			require.NoError(t, err)

			expectedBytes, err := os.ReadFile("testdata/" + tc.expectedFileName + ".txt")
			require.NoError(t, err)

			require.Equal(t, expectedBytes, actualBytes)
		})
	}
}
