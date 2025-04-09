// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleStruct struct {
	Name string `json:"name" yaml:"name"`
	Age  int    `json:"age" yaml:"age"`
}

func TestPathExists(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "testfile")
	defer os.Remove(tmpFile.Name())

	assert.True(t, PathExists(tmpFile.Name()))
	assert.False(t, PathExists("/some/nonexistent/path"))
}

func TestLoadObject(t *testing.T) {
	input := map[string]interface{}{"name": "John", "age": 30}
	var output sampleStruct

	LoadObject(input, &output)

	assert.Equal(t, "John", output.Name)
	assert.Equal(t, 30, output.Age)
}

func TestNormalizeFilename(t *testing.T) {
	fn := "my/test/file.json"
	fp := "/tmp/test"

	result, err := NormalizeFilename(fn, fp)

	assert.NoError(t, err)
	assert.Contains(t, result, "my_test_file.json")
	assert.Contains(t, result, "/tmp/test")
}

func TestWriteAndReadBytesToDisk(t *testing.T) {
	tmpDir := t.TempDir()
	dst := filepath.Join(tmpDir, "output.txt")

	err := WriteBytesToDisk([]byte("hello"), dst, false)
	assert.NoError(t, err)

	data, err := os.ReadFile(dst)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestWriteBytesToDisk_Overwrite(t *testing.T) {
	tmpDir := t.TempDir()
	dst := filepath.Join(tmpDir, "file.txt")

	// Write once
	err := WriteBytesToDisk([]byte("original"), dst, false)
	assert.NoError(t, err)

	// Overwrite
	err = WriteBytesToDisk([]byte("updated"), dst, true)
	assert.NoError(t, err)

	data, _ := os.ReadFile(dst)
	assert.Equal(t, "updated", string(data))
}

func TestWriteJsonToDisk(t *testing.T) {
	tmpDir := t.TempDir()
	obj := sampleStruct{"Alice", 25}

	err := WriteJsonToDisk(obj, "data.json", tmpDir)
	assert.NoError(t, err)

	content, _ := os.ReadFile(filepath.Join(tmpDir, "data.json"))
	assert.Contains(t, string(content), "Alice")
	assert.Contains(t, string(content), "25")
}

func TestWriteYamlToDisk(t *testing.T) {
	tmpDir := t.TempDir()
	obj := sampleStruct{"Bob", 40}

	err := WriteYamlToDisk(obj, "data.yaml", tmpDir)
	assert.NoError(t, err)

	content, _ := os.ReadFile(filepath.Join(tmpDir, "data.yaml"))
	assert.Contains(t, string(content), "Bob")
	assert.Contains(t, string(content), "40")
}

func TestWrite(t *testing.T) {
	tmpDir := t.TempDir()
	obj := sampleStruct{"Test", 10}

	err := Write(obj, "file.json", tmpDir, "json")
	assert.NoError(t, err)

	err = Write(obj, "file.yaml", tmpDir, "yaml")
	assert.NoError(t, err)

	jsonData, _ := os.ReadFile(filepath.Join(tmpDir, "file.json"))
	yamlData, _ := os.ReadFile(filepath.Join(tmpDir, "file.yaml"))

	assert.Contains(t, string(jsonData), "Test")
	assert.Contains(t, string(yamlData), "Test")
}

func TestReadStringFromFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(tmpFile, []byte("hello world"), 0644)

	data, err := ReadStringFromFile(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", data)
}

func TestEnsurePathExists(t *testing.T) {
	tmpDir := filepath.Join(t.TempDir(), "newdir")

	err := EnsurePathExists(tmpDir)
	assert.NoError(t, err)
	assert.True(t, PathExists(tmpDir))
}

func TestReadObjectFromDisk_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	obj := sampleStruct{"Test", 99}
	_ = WriteJsonToDisk(obj, "readme.json", tmpDir)

	var output sampleStruct
	err := ReadObjectFromDisk(filepath.Join(tmpDir, "readme.json"), &output)

	assert.NoError(t, err)
	assert.Equal(t, "Test", output.Name)
	assert.Equal(t, 99, output.Age)
}
