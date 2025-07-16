package lib

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

func GenerateDBUpdateFields(values map[UserColumn]string) string {
	/* Generates UPDATE fields using provided columns */
	update_fields := ""
	i := 1
	for col := range values {
		format := "%s = $%d, "
		if i == len(values) {
			format = "%s = $%d "
		}
		update_fields = update_fields + fmt.Sprintf(format, col.Column, i)
		i++
	}
	return update_fields
}

func GenerateDBOrFields(fields []UserColumn) string {
	/* Generates the OR conditions for given columns */
	or_fields := ""

	for i, col := range fields {
		if i == len(fields)-1 {
			// For last column, do not include OR and space
			or_fields = or_fields + fmt.Sprintf("%s.%s=$%d", col.Table, col.Column, i+1)
		} else {
			or_fields = or_fields + fmt.Sprintf("%s.%s=$%d OR ", col.Table, col.Column, i+1)
		}
	}

	return or_fields
}

func RandomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}
