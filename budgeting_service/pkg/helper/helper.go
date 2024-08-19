package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {

	var (
		i    = 1
		args = []interface{}{}
	)

	for k, v := range params {
		if k != "" {
			namedQuery = strings.ReplaceAll(namedQuery, "@"+k, "$"+strconv.Itoa(i))

			args = append(args, v)
			i++
		}
	}

	return namedQuery, args
}

func MakeWherePartOfQueryWithSearchFieldOfRequest(search string) (string, error) {

	var (
		paramsOfWhere = make(map[string]string)
		searchParams  = []string{}
		searchParam   = []string{}
		columnName    string
		input         string
		query         string
		err           error
	)

	// spliting search into small parts
	// each part has one column_name and its input like this name=Tashkent
	searchParams = strings.Split(search, ",")

	// iterating search parts to put them in map
	for _, field := range searchParams {
		// spliting column_name and its input
		searchParam = strings.Split(field, "=")

		// check invalid input
		if len(searchParam) != 2 || len(searchParam[1]) == 0 || len(searchParam[0]) == 0 {
			return "", fmt.Errorf("invalid search field %s", search)
		}

		columnName = searchParam[0]
		input = searchParam[1]

		paramsOfWhere[columnName] = input
	}

	if len(paramsOfWhere) == 0 {
		return "", nil
	}
	query = ` where `

	for columnName, input := range paramsOfWhere {
		// checking what is the input type
		// UUID
		if _, err = uuid.Parse(input); err == nil {
			query += fmt.Sprintf(" %s = '%s' and ", columnName, input)
			// bool
		} else if _, err = strconv.ParseBool(input); err == nil {
			query += fmt.Sprintf(" %s = %s and ", columnName, input)
			// float
		} else if _, err = strconv.ParseFloat(input, 64); err == nil {
			query += fmt.Sprintf(" %s = %s and ", columnName, input)
			// int
		} else if _, err = strconv.Atoi(input); err == nil {
			query += fmt.Sprintf(" %s = %s and ", columnName, input)
			// time with time zone
		} else if _, err = time.Parse("02-01-2006 15:04:05", input); err == nil {
			query += fmt.Sprintf(" %s = '%s'::timestamp without time zone and ", columnName, input)
			// time with time zone
		} else if _, err = time.Parse("02-01-2006T15:04:05Z", input); err == nil {
			query += fmt.Sprintf(" %s = '%s'::timestamp with time zone ann ", columnName, input)
		} else {
			// string
			query += fmt.Sprintf(" %s ilike '%%%s%%' and ", columnName, input)
		}
	}

	query = strings.TrimSuffix(query, "and ")

	return query, nil
}
