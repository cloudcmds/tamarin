package tests

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/cloudcmds/tamarin/evaluator"
	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/internal/parser"
	"github.com/cloudcmds/tamarin/object"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name              string
	Text              string
	ExpectedValue     string
	ExpectedType      string
	ExpectedErr       string
	ExpectedErrLine   int
	ExpectedErrColumn int
}

func readFile(name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func parseExpectedValue(filename, text string) (TestCase, error) {
	result := TestCase{}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "// ") {
			continue
		}
		line = strings.TrimPrefix(line, "// ")
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "expected value":
			result.ExpectedValue = val
		case "expected type":
			result.ExpectedType = val
		case "expected error":
			result.ExpectedErr = val
		case "expected error line":
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return result, err
			}
			result.ExpectedErrLine = intVal
		case "expected error column":
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return result, err
			}
			result.ExpectedErrColumn = intVal
		}
	}
	return result, nil
}

func getTestCase(name string) (TestCase, error) {
	input := readFile(name)
	testCase, err := parseExpectedValue(name, input)
	testCase.Name = name
	testCase.Text = input
	return testCase, err
}

func execute(ctx context.Context, input string) (object.Object, error) {
	result, err := exec.Execute(ctx, exec.Opts{
		Input:    string(input),
		Importer: &evaluator.SimpleImporter{},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func listTestFiles() []string {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	var testFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".tm") {
			testFiles = append(testFiles, f.Name())
		}
	}
	return testFiles
}

func TestFiles(t *testing.T) {
	only := "" // "11-19-22"
	for _, name := range listTestFiles() {
		if !strings.HasSuffix(name, ".tm") {
			continue
		}
		if only != "" && !strings.Contains(name, only) {
			continue
		}
		t.Run(name, func(t *testing.T) {
			tc, err := getTestCase(name)
			require.Nil(t, err)
			ctx := context.Background()
			result, err := execute(ctx, tc.Text)
			expectedType := object.Type(tc.ExpectedType)

			if tc.ExpectedValue != "" {
				require.Equal(t, tc.ExpectedValue, result.Inspect())
			}
			if tc.ExpectedType != "" {
				require.Equal(t, expectedType, result.Type())
			}
			if tc.ExpectedErr != "" {
				require.NotNil(t, err)
				require.Equal(t, tc.ExpectedErr, err.Error())
			}
			if tc.ExpectedErrColumn != 0 {
				require.NotNil(t, err)
				parserErr, ok := err.(parser.ParserError)
				require.True(t, ok)
				require.Equal(t, tc.ExpectedErrColumn, parserErr.StartPosition().ColumnNumber())
			}
			if tc.ExpectedErrLine != 0 {
				require.NotNil(t, err)
				parserErr, ok := err.(parser.ParserError)
				require.True(t, ok)
				require.Equal(t, tc.ExpectedErrLine, parserErr.StartPosition().LineNumber())
			}
		})
	}
}
