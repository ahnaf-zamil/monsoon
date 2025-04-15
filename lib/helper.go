package lib

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
)

/* Helper functions */
func GenerateSocketID() string {
	// Generates 8 byte IDs for sockets
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func PrettyPrintSyncMap(data *sync.Map) {
	// Handy function for pretty printing sync.Map objects
	m := map[string]any{}
	data.Range(func(key, value any) bool {
		m[fmt.Sprint(key)] = value
		return true
	})

	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
