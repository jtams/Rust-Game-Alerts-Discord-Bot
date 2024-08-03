package tracker

import "strings"

type Group struct {
	Name  string  `json:"name"`
	Users []*User `json:"users"`
}

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

func (g *Group) RemoveUserByUsername(username string) bool {
	for i, user := range g.Users {
		if user.GetUsername() == username {
			g.Users = append(g.Users[:i], g.Users[i+1:]...)
			return true
		}
	}

	return false
}

func (g *Group) RemoveUserByID(id string) {
	for i, user := range g.Users {
		if user.ID == id {
			g.Users = append(g.Users[:i], g.Users[i+1:]...)
			break
		}
	}
}
