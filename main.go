package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"syscall"

	"github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/term"
)

type (
	ConfigEntry struct {
		Uuid    string               `json:"uuid"`
		Exports []ConfigEntryExports `json:"exports"`
	}
	ConfigEntryExports struct {
		Field    string `json:"field"`
		Variable string `json:"variable"`
	}
	Entry struct {
		isDeleted bool
		entity    gokeepasslib.Entry
	}
	Group struct {
		isDeleted bool
		level     int
		entity    gokeepasslib.Group
	}
)

func main() {
	if len(os.Args) < 3 || slices.Contains(os.Args[1:3], "") {
		fmt.Fprint(os.Stderr, "Missing arguments.\nUsage: keexp <database_path> <config_path>\n")
		os.Exit(1)
	}

	if len(os.Args) > 3 {
		fmt.Fprint(os.Stderr, "Too many arguments.\nUsage: keexp <database_path> <config_path>\n")
		os.Exit(1)
	}

	databasePath := os.Args[1]
	configPath := os.Args[2]

	config, err := readConfig(configPath)
	checkError(err)

	fmt.Fprint(os.Stderr, "Enter password: ")
	passwordBytes, err := term.ReadPassword(syscall.Stdin)
	fmt.Fprintln(os.Stderr)
	checkError(err)

	databaseFile, err := os.Open(databasePath)
	checkError(err)

	database := gokeepasslib.NewDatabase()
	database.Credentials = gokeepasslib.NewPasswordCredentials(string(passwordBytes))
	err = gokeepasslib.NewDecoder(databaseFile).Decode(database)
	checkError(err)

	err = database.UnlockProtectedEntries()
	checkError(err)

	entriesByUuid := getEntriesByUuid(database)

	for _, configEntry := range config {
		uuidBytes, err := hex.DecodeString(configEntry.Uuid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: UUID \"%s\" could not be decoded: %s\n", configEntry.Uuid, err)
			continue
		}

		uuidLength := len(uuidBytes)
		if uuidLength != 16 {
			fmt.Fprintf(os.Stderr, "Warning: UUID \"%s\" has invalid length: %d\n", configEntry.Uuid, uuidLength)
			continue
		}

		uuid := gokeepasslib.UUID(uuidBytes)
		entry, entryExists := entriesByUuid[uuid]
		if !entryExists {
			fmt.Fprintf(os.Stderr, "Warning: Entry with UUID \"%s\" could not be found\n", configEntry.Uuid)
			continue
		}
		if entry.isDeleted {
			fmt.Fprintf(os.Stderr, "Warning: Entry with UUID \"%s\" is in recycle bin\n", configEntry.Uuid)
		}

		valuesByKey := getValuesByKey(entry)
		for _, export := range configEntry.Exports {
			value, valueExists := valuesByKey[export.Field]
			if !valueExists {
				fmt.Fprintf(os.Stderr, "Warning: Field \"%s\" for entry with UUID \"%s\" could not be found\n", export.Field, configEntry.Uuid)
				continue
			}
			if value == "" {
				fmt.Fprintf(os.Stderr, "Warning: Field \"%s\" for entry with UUID \"%s\" is empty\n", export.Field, configEntry.Uuid)
				continue
			}
			escapedValue := strings.ReplaceAll(value, "'", "'\\''")
			fmt.Fprintln(os.Stdout, "export "+export.Variable+"='"+escapedValue+"';")
		}
	}
}

func readConfig(configPath string) ([]ConfigEntry, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	configJson, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config []ConfigEntry
	err = json.Unmarshal(configJson, &config)
	return config, err
}

func getEntriesByUuid(database *gokeepasslib.Database) map[gokeepasslib.UUID]*Entry {
	entriesByUuid := make(map[gokeepasslib.UUID]*Entry)
	groups := make([]Group, len(database.Content.Root.Groups))
	for i, group := range database.Content.Root.Groups {
		groups[i] = Group{
			isDeleted: false,
			level:     0,
			entity:    group,
		}
	}

	for len(groups) > 0 {
		currentGroup := groups[0]
		isRecycleBin := currentGroup.level == 1 && currentGroup.entity.Name == "Recycle Bin"
		isDeleted := currentGroup.isDeleted || isRecycleBin

		for _, entry := range currentGroup.entity.Entries {
			if entriesByUuid[entry.UUID] != nil {
				fmt.Fprintf(os.Stderr, "Warning: Found multiple entries with UUID \"%s\"\n", entry.UUID)
				continue
			}
			entriesByUuid[entry.UUID] = &Entry{
				isDeleted: isDeleted,
				entity:    entry,
			}
		}

		subGroupCount := len(currentGroup.entity.Groups)
		updatedGroups := make([]Group, len(groups)+subGroupCount-1)
		for i, subGroup := range currentGroup.entity.Groups {
			updatedGroups[i] = Group{
				isDeleted: isDeleted,
				level:     currentGroup.level + 1,
				entity:    subGroup,
			}
		}
		for i := 1; i < len(groups); i++ {
			updatedGroups[subGroupCount+i-1] = groups[i]
		}
		groups = updatedGroups
	}

	return entriesByUuid
}

func getValuesByKey(entry *Entry) map[string]string {
	valuesByKey := make(map[string]string, len(entry.entity.Values))
	for _, value := range entry.entity.Values {
		if valuesByKey[value.Key] != "" {
			fmt.Fprintf(os.Stderr, "Warning: Found values with duplicate key \"%s\" for entry with UUID \"%s\"\n", value.Key, entry.entity.UUID)
			continue
		}
		valuesByKey[value.Key] = value.Value.Content
	}
	return valuesByKey
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
