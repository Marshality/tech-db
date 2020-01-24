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
		} else {
			query = fmt.Sprintf(SelectPostsByThreadFlat, "", "ASC", "$2")
		}
	case TREE_SORT:
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadTree, "WHERE", "DESC", "$2")
		} else {
			query = fmt.Sprintf(SelectPostsByThreadTree, "WHERE", "ASC", "$2")
		}
	case PARENT_TREE_SORT:
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadParentTree, "DESC",
				"WHERE roots.root <= $2 ORDER BY roots.path[1] DESC, " +
				"array_remove(roots.path, roots.path[1])")
		} else {
			query = fmt.Sprintf(SelectPostsByThreadParentTree, "ASC",
				"WHERE roots.root <= $2 ORDER BY roots.path[1] ASC, " +
				"array_remove(roots.path, roots.path[1])")
		}
	}

	return query
}

func (qm *QueryManager) SelectPostsByThreadSince(desc bool, sort string) string {
	var query string

	switch sort {
	case FLAT_SORT, "":
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadFlat, "AND id < $2", "DESC", "$3")
		} else {
			query = fmt.Sprintf(SelectPostsByThreadFlat, "AND id > $2", "ASC", "$3")
		}
	case TREE_SORT:
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadTree,"JOIN posts P ON P.id = $2 WHERE posts.path < p.path AND",
				"DESC",
				"$3")
		} else {
			query = fmt.Sprintf(SelectPostsByThreadTree,"JOIN posts P ON P.id = $2 WHERE posts.path > p.path AND",
				"ASC",
				"$3")
		}
	case PARENT_TREE_SORT:
		if desc {
			query = fmt.Sprintf(SelectPostsByThreadParentTree, "DESC",
				"JOIN roots R ON R.id = $2 WHERE roots.root <= $3 + R.root " +
				"AND (roots.root > R.root OR roots.root = R.root AND roots.path > R.path) ORDER BY roots.path[1] DESC, " +
				"array_remove(roots.path, roots.path[1])")
		} else {
			query = fmt.Sprintf(SelectPostsByThreadParentTree, "ASC",
				"JOIN roots R ON R.id = $2 WHERE roots.root <= $3 + R.root " +
				"AND (roots.root > R.root OR roots.root = R.root AND roots.path > R.path) ORDER BY roots.path[1] ASC, " +
				"array_remove(roots.path, roots.path[1])")
		}
	}

	return query
}
