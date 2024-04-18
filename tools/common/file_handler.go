package common

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strings"
)

var (
	errDirNotExists = errors.New("provided directory does not exist")
	errNotDirectory = errors.New("provided path is not a directory")
)

func ValidateDirectory(pathToDirectory string) error {
	fileInfo, err := os.Stat(pathToDirectory)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %s", errDirNotExists, pathToDirectory)
		}
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%w: %s", errNotDirectory, pathToDirectory)
	}

	return nil
}

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

func GenerateYamlFileFromObject(obj interface{}, filePath string) error {
	yamlVal := GenerateYaml(obj)

	err := os.WriteFile(filePath, []byte(yamlVal), 0600)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func GenerateYaml(obj interface{}) string {
	reflectValue := reflect.ValueOf(obj)
	var yamlBuilder strings.Builder
	generateYamlWithComments(&yamlBuilder, reflectValue, 0, "")
	return yamlBuilder.String()
}

// generateYamlWithComments uses a "comment" tag in the struct definition to generate YAML with comments on corresponding lines.
// Note: Map support is missing!
func generateYamlWithComments(yamlBuilder *strings.Builder, obj reflect.Value, indentLevel int, commentPrefix string) {
	t := obj.Type()

	indentPrefix := strings.Repeat("  ", indentLevel)
	originalCommentPrefix := commentPrefix
	for i := 0; i < t.NumField(); i++ {
		commentPrefix = originalCommentPrefix
		field := t.Field(i)
		value := obj.Field(i)
		yamlTag := field.Tag.Get("yaml")
		commentTag := field.Tag.Get("comment")

		// comment-out non-required empty attributes
		if value.IsZero() && !strings.Contains(commentTag, "required") {
			commentPrefix = "# "
		}

		if value.Kind() == reflect.Struct {
			if commentTag == "" {
				yamlBuilder.WriteString(fmt.Sprintf("%s%s%s:\n", commentPrefix, indentPrefix, yamlTag))
			} else {
				yamlBuilder.WriteString(fmt.Sprintf("%s%s%s: # %s\n", commentPrefix, indentPrefix, yamlTag, commentTag))
			}
			generateYamlWithComments(yamlBuilder, value, indentLevel+1, commentPrefix)
			continue
		}

		if value.Kind() == reflect.Slice {
			if commentTag == "" {
				yamlBuilder.WriteString(fmt.Sprintf("%s%s%s:\n", commentPrefix, indentPrefix, yamlTag))
			} else {
				yamlBuilder.WriteString(fmt.Sprintf("%s%s%s: # %s\n", commentPrefix, indentPrefix, yamlTag, commentTag))
			}

			if value.Len() == 0 {
				yamlBuilder.WriteString(fmt.Sprintf("%s%s  -\n", commentPrefix, indentPrefix))
			}
			for j := 0; j < value.Len(); j++ {
				valueStr := getValueStr(value.Index(j))
				yamlBuilder.WriteString(fmt.Sprintf("%s%s  - %s\n", "", indentPrefix, valueStr))
			}
			continue
		}

		valueStr := getValueStr(value)
		if commentTag == "" {
			yamlBuilder.WriteString(fmt.Sprintf("%s%s%s: %s\n", commentPrefix, indentPrefix,
				yamlTag, valueStr))
		} else {
			yamlBuilder.WriteString(fmt.Sprintf("%s%s%s: %s # %s\n", commentPrefix, indentPrefix,
				yamlTag, valueStr, commentTag))
		}
	}
}

func getValueStr(value reflect.Value) string {
	valueStr := ""
	if value.Kind() == reflect.String {
		valueStr = fmt.Sprintf("\"%v\"", value.Interface())
	} else {
		valueStr = fmt.Sprintf("%v", value.Interface())
	}
	return valueStr
}
