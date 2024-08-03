package tracker

import "strings"

// Takes a list of users, a function to get the username of the user and returns the user that the name matches
// TODO: This search algorithm is inefficient and it hurts
func SearchUsersWithUserCreatedName[T any](userList []T, getUsername func(T) string, username string, exact bool, casesensitive bool) (int, *T) {
	lowerUsername := strings.ToLower(username)

	// Exact match
	for i, user := range userList {
		if getUsername(user) == username {
			return i, &user
		}
	}

	if exact && !casesensitive {
		return -1, nil
	}

	// Case insensitive exact match
	for i, user := range userList {
		if strings.ToLower(getUsername(user)) == lowerUsername {
			return i, &user
		}
	}

	if exact {
		return -1, nil
	}

	// Match the start of the name
	for i, user := range userList {
		if strings.HasPrefix(getUsername(user), username) {
			return i, &user
		}
	}

	// Case insensitive match the start of the name
	if !casesensitive {
		for i, user := range userList {
			if strings.HasPrefix(strings.ToLower(getUsername(user)), lowerUsername) {
				return i, &user
			}
		}
	}

	// Substring match (more than 5 characters)
	for i, user := range userList {
		if len(username) > 5 && strings.Contains(getUsername(user), username) {
			return i, &user
		}
	}

	// Case insensitive substring match (more than 5 characters)
	if !casesensitive {
		for i, user := range userList {
			if len(username) > 5 && strings.Contains(strings.ToLower(getUsername(user)), lowerUsername) {
				return i, &user
			}
		}
	}

	return -1, nil
}
