package storage

import (
	"fmt"
	"strings"
)

var sortColumns = map[string]string{
	"date":    "date_sent",
	"subject": "subject",
	"from":    "from_address",
}

// buildWhere returns a SQL WHERE clause (possibly empty) and its bound
// arguments for the given search criteria. Column names are never derived
// from user input; only literal values are parameterized.
func buildWhere(s *Search) (string, []any) {
	if s == nil {
		return "", nil
	}

	var clauses []string
	var args []any

	if s.Query != "" {
		clauses = append(clauses, "(subject LIKE ? OR text_body LIKE ? OR html_body LIKE ?)")
		like := "%" + s.Query + "%"
		args = append(args, like, like, like)
	}
	if s.From != "" {
		clauses = append(clauses, "from_address LIKE ?")
		args = append(args, "%"+s.From+"%")
	}
	if s.To != "" {
		clauses = append(clauses, "to_addresses LIKE ?")
		args = append(args, "%"+s.To+"%")
	}
	if !s.Start.IsZero() {
		clauses = append(clauses, "date_sent >= ?")
		args = append(args, s.Start.UTC().Format(timeFormat))
	}
	if !s.End.IsZero() {
		clauses = append(clauses, "date_sent <= ?")
		args = append(args, s.End.UTC().Format(timeFormat))
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}

// buildOrderBy validates the requested sort field/direction against an
// allow-list before use, since these values come from client query params.
func buildOrderBy(s *Search) string {
	column := "date_sent"
	if s != nil {
		if c, ok := sortColumns[s.SortField]; ok {
			column = c
		}
	}

	dir := "DESC"
	if s != nil && strings.EqualFold(s.SortDir, "asc") {
		dir = "ASC"
	}

	return fmt.Sprintf("ORDER BY %s %s", column, dir)
}
