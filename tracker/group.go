package tracker

import (
	"strings"
)

// Group containing users
// used to divide the users in the tracker message
type Group struct {
	Name     string  `json:"name"`
	Users    []*User `json:"users"`
	Notes    string  `json:"notes"`
	Location string  `json:"location"`
}

// Creates a new group with empty user list
func NewGroup(name string) *Group {
	return &Group{
		Name:  strings.ToLower(name),
		Users: []*User{},
	}
}

func (g *Group) AddUser(username string) {
	user := &User{
		ID:        "",
		Usernames: []string{username},
		Group:     g.Name,
		Status:    StatusUntracked,
	}

	g.Users = append(g.Users, user)
}

// Removes user by username, searching is case insensitive and can be partial
// if you have the user's ID it's highly recommended to use RemoveUserByID.
// This function should be used when a user inputs a username to remove.
func (g *Group) RemoveUserByUsername(username string) bool {
	if len(g.Users) == 0 {
		return false
	}

	i, user := SearchUsersWithUserCreatedName(g.Users, func(u *User) string { return u.GetUsername() }, username, false, false)
	if i > -1 && user != nil {
		g.Users = append(g.Users[:i], g.Users[i+1:]...)
		return true
	}

	return false
}

// Removes user by ID
func (g *Group) RemoveUserByID(id string) bool {
	if g.Users == nil || len(g.Users) == 0 {
		return false
	}
	for i, user := range g.Users {
		if user.ID == id {
			g.Users = append(g.Users[:i], g.Users[i+1:]...)
			return true
		}
	}

	return false
}
