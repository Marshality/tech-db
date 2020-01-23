package queries

import (
	"fmt"
	. "github.com/Marshality/tech-db/tools"
)

type QueryManager struct{}
var QM *QueryManager

func (qm *QueryManager) SelectThreadsWhereForum(desc bool) string {
	var query string

	switch desc {
	case true:
		query = fmt.Sprintf(SelectThreadsWhereForum, "DESC")
	case false:
		query = fmt.Sprintf(SelectThreadsWhereForum, "ASC")
	}

	return query
}

func (qm *QueryManager) SelectThreadsWhereForumAndCreated(desc bool) string {
	var query string

	switch desc {
	case true:
		query = fmt.Sprintf(SelectThreadsWhereForumAndCreated, "<=", "DESC")
	case false:
		query = fmt.Sprintf(SelectThreadsWhereForumAndCreated, ">=", "ASC")
	}

	return query
}

func (qm *QueryManager) SelectPostsByThread(desc bool, sort string) string {
	var query string

	switch sort {
	case FLAT_SORT, "":
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadFlat, "", "DESC", "$2")
		}

		query = fmt.Sprintf(SelectPostsByThreadFlat, "", "ASC", "$2")
	}

	return query
}

func (qm *QueryManager) SelectPostsByThreadSince(desc bool, sort string) string {
	var query string

	switch sort {
	case FLAT_SORT, "":
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadFlat, "AND id < $2", "DESC", "$3")
		}

		query = fmt.Sprintf(SelectPostsByThreadFlat, "AND id > $2", "ASC", "$3")
	}

	return query
}
