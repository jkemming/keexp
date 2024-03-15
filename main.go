package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/term"
)

type (
	Config struct {
		Database string        `json:"database"`
		Entries  []ConfigEntry `json:"entries"`
	}
	ConfigEntry struct {
		Uuid    string               `json:"uuid"`
		Exports []ConfigEntryExports `json:"exports"`
	}
	ConfigEntryExports struct {
		Field    string `json:"field"`
		Variable string `json:"variable"`
	}
)

func main() {
	config, err := readConfig()
	checkError(err)

	fmt.Fprint(os.Stderr, "Enter password for ", config.Database, ": ")
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	fmt.Fprintln(os.Stderr)
	checkError(err)

	databaseFile, err := os.Open(config.Database)
	checkError(err)

	database := gokeepasslib.NewDatabase()
	database.Credentials = gokeepasslib.NewPasswordCredentials(string(passwordBytes))
	err = gokeepasslib.NewDecoder(databaseFile).Decode(database)
	checkError(err)

	err = database.UnlockProtectedEntries()
	checkError(err)

	entriesByUuid := getEntriesByUuid(database)

	for _, configEntry := range config.Entries {
		uuidBytes, err := hex.DecodeString(configEntry.Uuid)
		checkError(err)
		if len(uuidBytes) != 16 {
			fmt.Fprintln(os.Stderr, "Invalid UUID: "+configEntry.Uuid)
		}
		uuid := gokeepasslib.UUID(uuidBytes)
		valuesByKey := getValuesByKey(entriesByUuid[uuid])
		for _, export := range configEntry.Exports {
			value, valueExists := valuesByKey[export.Field]
			if !valueExists {
				// TODO Print out the entry's UUID
				fmt.Fprintln(os.Stderr, "Field could not be found: "+export.Field)
			}
			escapedValue := strings.ReplaceAll(value, "'", "'\\''")
			fmt.Fprintln(os.Stdout, "export "+export.Variable+"='"+escapedValue+"'")
		}
	}
}

func getEntriesByUuid(database *gokeepasslib.Database) map[gokeepasslib.UUID]*gokeepasslib.Entry {
	entriesByUuid := make(map[gokeepasslib.UUID]*gokeepasslib.Entry)
	groups := database.Content.Root.Groups

	for len(groups) > 0 {
		for _, entry := range groups[0].Entries {
			if entriesByUuid[entry.UUID] != nil {
				// TODO Print out the entry's UUID
				fmt.Fprintln(os.Stderr, "Warning: Found entries with duplicate UUID")
				continue
			}
			entriesByUuid[entry.UUID] = &entry
		}

		groups = append(groups, groups[0].Groups...)
		groups[0] = groups[len(groups)-1]
		groups = groups[:len(groups)-1]
	}

	return entriesByUuid
}

func getValuesByKey(entry *gokeepasslib.Entry) map[string]string {
	valuesByKey := make(map[string]string, len(entry.Values))
	for _, value := range entry.Values {
		if valuesByKey[value.Key] != "" {
			// TODO Print out the entry's UUID
			fmt.Fprintln(os.Stderr, "Warning: Found values with duplicate key")
			continue
		}
		valuesByKey[value.Key] = value.Value.Content
	}
	return valuesByKey
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
