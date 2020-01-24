package queries

import (
	"fmt"
	"github.com/Marshality/tech-db/models"
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

func (qm *QueryManager) SelectForumUsers(desc bool, limit uint64) string {
	var query string

	switch desc {
	case true:
		query = fmt.Sprintf(SelectUsersWhereForumSlug, "", "DESC")
	case false:
		query = fmt.Sprintf(SelectUsersWhereForumSlug, "", "ASC")
	}

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	return query
}

func (qm *QueryManager) SelectForumUsersSince(desc bool, limit uint64) string {
	var query string

	switch desc {
	case true:
		query = fmt.Sprintf(SelectUsersWhereForumSlug, "AND U.nickname < $2", "DESC")
	case false:
		query = fmt.Sprintf(SelectUsersWhereForumSlug, "AND U.nickname > $2", "ASC")
	}

	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	return query
}

// related[0] - user
// related[1] - forum
// related[2] - thread
func (qm *QueryManager) SelectPostWhereID(userFlag, threadFlag, forumFlag bool, related ...string) (func(string, *models.User, *models.Forum, *models.Thread, *models.Post)interface{}, string) {
	helper := func(columnName string, u *models.User, f *models.Forum, t *models.Thread, p *models.Post) interface{} {
		switch columnName {
		case "post_author":
			return &p.Author
		case "post_created":
			return &p.CreatedAt
		case "post_id":
			return &p.ID
		case "post_is_edited":
			return &p.IsEdited
		case "post_message":
			return &p.Message
		case "post_parent":
			return &p.Parent
		case "post_thread":
			return &p.Thread
		case "post_forum":
			return &p.Forum
		case "user_about":
			return &u.About
		case "user_email":
			return &u.Email
		case "user_fullname":
			return &u.Fullname
		case "user_nickname":
			return &u.Nickname
		case "forum_slug":
			return &f.Slug
		case "forum_posts":
			return &f.Posts
		case "forum_threads":
			return &f.Threads
		case "forum_title":
			return &f.Title
		case "forum_user":
			return &f.User
		case "thread_user":
			return &t.Author
		case "thread_created":
			return &t.CreatedAt
		case "thread_forum":
			return &t.Forum
		case "thread_id":
			return &t.ID
		case "thread_message":
			return &t.Message
		case "thread_slug":
			return &t.Slug
		case "thread_title":
			return &t.Title
		case "thread_votes":
			return &t.Votes
		default:
			return nil
		}
	}

	var query string
	extended1 := ""
	extended2 := ""

	var user, forum bool

	if userFlag {
		user = true
		extended1 += ", U.about AS user_about, U.email AS user_email, U.fullname AS user_fullname, U.nickname AS user_nickname "
	}

	if forumFlag {
		forum = true
		extended1 += ", F.slug AS forum_slug, F.posts AS forum_posts, F.threads AS forum_threads, F.title AS forum_title, F.usr AS forum_user "
	}

	if threadFlag {
		extended1 += ", T.author AS thread_user, T.created_at AS thread_created, T.forum AS thread_forum, " +
			"T.id AS thread_id, T.message AS thread_message, T.slug AS thread_slug, T.title AS thread_title, T.votes AS thread_votes "
	}

	if forum {
		extended2 += " JOIN forums F ON T.forum = F.slug "
	}

	if user {
		extended2 += " JOIN users U ON P.author = U.nickname "
	}

	query = fmt.Sprintf(SelectPostWhereID, extended1, extended2)

	return helper, query
}
