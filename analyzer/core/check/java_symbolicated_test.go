package check

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"analyzer/core/mast"
	"analyzer/core/symbolication"
	"analyzer/core/translation"
	ts "analyzer/core/treesitter"
)

func TestJavaSymbolicatedChecker(t *testing.T) {
	testcases := []struct {
		file     string
		expected string
	}{
		{
			file:     "import_slf4j.java",
			expected: "Logger",
		},
		{
			file:     "import_log4j.java",
			expected: "Logger",
		},
		{
			file:     "import_java_util.java",
			expected: "Logger",
		},
		{
			file:     "import_unsupported.java",
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run("test if java symbolicated checker correctly sets logger imported flag", func(t *testing.T) {
			path := _metaTestCheckDir + "java/" + tc.file
			suffix := filepath.Ext(path)
			// Build tree-sitter node.
			tsNode, err := ts.ParseFile(path)
			require.NoError(t, err)

			// Translate to MAST node.
			mastNode, err := translation.Run(tsNode, suffix)
			require.NoError(t, err)
			require.NotNil(t, mastNode)

			// Run symbolication.
			symbols, err := symbolication.Run([]mast.Node{mastNode}, suffix)
			require.NoError(t, err)
			require.NotNil(t, symbols)

			// Create a symbolicated checker and check the mast node against itself. We are only
			// interested in finding out if the importedLoggerClass is properly set.
			symbolicatedChecker := SymbolicatedChecker{
				GenericChecker: GenericChecker{},
				baseSymbols:    symbols,
				lastSymbols:    symbols,
			}
			langChecker := NewJavaSymbolicatedChecker(&symbolicatedChecker, true /* loggingOn */)
			require.True(t, langChecker.LoggingOn)
			r, err := symbolicatedChecker.Equal([]mast.Node{mastNode}, []mast.Node{mastNode}, langChecker)
			require.NoError(t, err)
			require.True(t, r)

			require.Equal(t, tc.expected, langChecker.importedLoggerClass)
		})
	}

}
