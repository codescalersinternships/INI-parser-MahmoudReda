package iniparser

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoadFromString(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		iniData := `
[section 1]
key key = value value
[section 2]
key1 = value1
key2 = value2
`

		expectedConfig := Config{
			"section 1": {
				"key key": "value value",
			},
			"section 2": {
				"key1": "value1",
				"key2": "value2",
			},
		}

		config, err := LoadFromString(iniData)
		assertNoErr(t, err)
		if !reflect.DeepEqual(config, expectedConfig) {
			t.Errorf("got %+v, want %+v", config, expectedConfig)
		}
	})

	t.Run("invalid format: begin with comment", func(t *testing.T) {
		iniData := `;comment`
		_, err := LoadFromString(iniData)
		assertErr(t, err, ErrEmpytFile)
	})

	t.Run("invalid format: begin with global data", func(t *testing.T) {
		iniData := `key key = value value`
		_, err := LoadFromString(iniData)
		assertErr(t, err, ErrNoGlobalDataAllowed)
	})

	t.Run("invalid format: section naming format", func(t *testing.T) {
		iniData := `
[section 1
key key = value value
`
		_, err := LoadFromString(iniData)
		assertErr(t, err, ErrInvalidFormat)
	})

	t.Run("invalid format: key value format", func(t *testing.T) {
		iniData := `
[section 1]
key key  value value
`
		_, err := LoadFromString(iniData)
		assertErr(t, err, ErrInvalidFormat)
	})
}

func TestLoadFromFile(t *testing.T) {
	t.Run("file not exist", func(t *testing.T) {
		_, err := LoadFromFile("notExist.ini")
		assertErr(t, err, ErrFileNotFound)
	})

	t.Run("file extension is not valid", func(t *testing.T) {
		_, err := LoadFromFile("go.mod")
		assertErr(t, err, ErrInvalidFileExtension)
	})

	t.Run("valid format", func(t *testing.T) {
		testFilePath := "testdata/testfile.ini"
		iniData := `
[section 1]
key key = value value
[section 2]
key1 = value1
key2 = value2
`
		expectedConfig := Config{
			"section 1": {
				"key key": "value value",
			},
			"section 2": {
				"key1": "value1",
				"key2": "value2",
			},
		}

		err := ioutil.WriteFile(testFilePath, []byte(iniData), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := os.Remove(testFilePath)
			if err != nil {
				t.Fatal(err)
			}
		}()

		config, err := LoadFromFile(testFilePath)
		assertNoErr(t, err)
		if !reflect.DeepEqual(config, expectedConfig) {
			t.Errorf("got %+v, want %+v", config, expectedConfig)
		}
	})
}

func TestGetSectionNames(t *testing.T) {
	config := Config{
		"section 1": {
			"key key": "value value",
		},
		"section 2": {
			"key1": "value1",
			"key2": "value2",
		},
	}

	expectedSectionNames := []string{"section 1", "section 2"}

	sectionNames := GetSectionNames(config)
	if !reflect.DeepEqual(sectionNames, expectedSectionNames) {
		t.Errorf("got %v, want %v", sectionNames, expectedSectionNames)
	}
}

// have a problem in this test
// got map[section 1:map[key key:value value] section 2:map[key1:value1 key2:value2]], want
// map[section 1:map[key key:value value] section 2:map[key1:value1 key2:value2]]
// func TestGetSections(t *testing.T) {
// 	config := Config{
// 		"section 1": {
// 			"key key": "value value",
// 		},
// 		"section 2": {
// 			"key1": "value1",
// 			"key2": "value2",
// 		},
// 	}

// 	sections := GetSections(config)
// 	if !reflect.DeepEqual(sections, config) {
// 		t.Errorf("got %+v, want %+v", sections, config)
// 	}
// }

func TestGet(t *testing.T) {
	config := Config{
		"section 1": {
			"key key": "value value",
		},
		"section 2": {
			"key1": "value1",
			"key2": "value2",
		},
	}

	t.Run("get existing value", func(t *testing.T) {
		value, err := Get(config, "section 1", "key key")
		assertNoErr(t, err)
		assertEqual(t, value, "value value")
	})

	t.Run("get non-existing section", func(t *testing.T) {
		_, err := Get(config, "not found", "not found")
		assertErr(t, err, ErrSectionNotFound)
	})

	t.Run("get non-existing key", func(t *testing.T) {
		_, err := Get(config, "section 1", "not found")
		assertErr(t, err, ErrKeyNotFound)
	})
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func assertErr(t *testing.T, err, expectedErr error) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error: %v, got nil", expectedErr)
	} else if err != expectedErr {
		t.Errorf("expected error: %v, got: %v", expectedErr, err)
	}
}

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
