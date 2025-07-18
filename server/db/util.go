package db

import "fmt"

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
