package lib

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

/*
This is the place where I dump all my helper functions

May refactor/reorganize later
*/

/* General helper functions */
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

/* API-sided */
func WriteAPIError(c *gin.Context, message string, rs *APIResponse, r_code int) {
	rs.Err = true
	rs.Message = message
	c.JSON(r_code, rs)
}

func HandleServerError(c *gin.Context, rs *APIResponse, err error) {
	// Instead of rewriting the same two lines every single time, just run it for error cases
	WriteAPIError(c, "Internal server error", rs, http.StatusInternalServerError)
	log.Println(err)
}

/* DB-sided */
func GenerateDBQueryFields(cols []UserColumn) string {
	/* Generates SELECT query fields using provided columns */

	selected_columns := ""

	for i, c := range cols {
		if i == len(cols)-1 {
			// For last column, do not include comma and space
			selected_columns = selected_columns + fmt.Sprintf("%s.%s", c.Table, c.Column)
		} else {
			selected_columns = selected_columns + fmt.Sprintf("%s.%s, ", c.Table, c.Column)
		}
	}
	return selected_columns
}
