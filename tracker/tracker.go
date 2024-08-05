package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var logger *slog.Logger = slog.Default()

// Player Tracker manages the tracking data and organization.
// The heart of the bot.
type PlayerTracker struct {
	BattleMetricsID string `json:"battleMetricsID"`
	ServerName      string `json:"serverName"`

	Groups []*Group `json:"groups"`

	// Update interval in seconds
	Interval int `json:"interval"`

	Running bool `json:"running"`

	// Signal channel after each update
	Channel chan time.Time `json:"-"`

	// Server population [online, capacity]
	Online [2]int `json:"online"`
}

// Creates new default player tracker.
// Includes default group names: Squad, Allies, Neighbors, Enemies, Others.
func NewPlayerTracker() *PlayerTracker {
	tracker := &PlayerTracker{
		BattleMetricsID: "",
		Groups:          []*Group{},
		Interval:        10,
		Running:         false,
		Channel:         make(chan time.Time),
		Online:          [2]int{0, 0},
	}

	// Default groups
	tracker.AddNewGroup("squad")
	tracker.AddNewGroup("allies")
	tracker.AddNewGroup("neighbors")
	tracker.AddNewGroup("enemies")
	tracker.AddNewGroup("others")

	return tracker
}

// Adds a group to the tracker.
func (tracker *PlayerTracker) AddGroup(group *Group) error {
	group.Name = strings.ToLower(group.Name)
	if tracker.GetGroupByName(group.Name) != nil {
		return errors.New("group with that name already exists")
	}

	tracker.Groups = append(tracker.Groups, group)
	return nil
}

// Creates a new group and adds it to the tracker.
func (tracker *PlayerTracker) AddNewGroup(name string) error {
	newGroup := NewGroup(name)
	return tracker.AddGroup(newGroup)
}

// Removes a group from the tracker by name.
func (tracker *PlayerTracker) RemoveGroup(groupName string) bool {
	groupName = strings.ToLower(groupName)
	for i, group := range tracker.Groups {
		if group.Name == groupName {
			tracker.Groups = append(tracker.Groups[:i], tracker.Groups[i+1:]...)
			return true
		}
	}
	return false
}

// Adds a user to a group by username.
func (tracker *PlayerTracker) AddUserToGroup(username string, groupName string) error {
	groupName = strings.ToLower(groupName)
	group := tracker.GetGroupByName(groupName)
	if group == nil {
		return errors.New("group not found")
	}

	group.AddUser(username)
	return nil
}

// Adds a user to a group by username.
func (tracker *PlayerTracker) AddUserToGroupByID(id string, groupName string) error {
	groupName = strings.ToLower(groupName)
	group := tracker.GetGroupByName(groupName)
	if group == nil {
		return errors.New("group not found")
	}

	group.AddUserByID(id)
	return nil
}

// Removes a user from a group by username.
func (tracker *PlayerTracker) RemoveUserByUsername(username string) bool {
	deleted := 0
	for _, group := range tracker.Groups {
		if group.RemoveUserByUsername(username) {
			deleted++
		}
	}

	return deleted > 0
}

// Moves a user from one group to another.
func (tracker *PlayerTracker) MoveUserToGroup(username string, groupName string) bool {
	groupName = strings.ToLower(groupName)

	// Find user
	_, userP := SearchUsersWithUserCreatedName(tracker.Users(), func(u *User) string { return u.GetUsername() }, username, false, false)
	if userP == nil || *userP == nil {
		return false
	}
	user := *userP

	// Already in group
	if user.Group == groupName {
		return false
	}

	// Get new group
	group := tracker.GetGroupByName(groupName)
	if group == nil {
		return false
	}

	// Add user to new group
	group.Users = append(group.Users, user)

	// Remove them from existing group
	res := tracker.GetGroupByName(user.Group).RemoveUserByID(user.ID)

	// Update user's group value
	user.Group = group.Name

	return res
}

// Removes a user from a group by username and group name.
func (tracker *PlayerTracker) RemoveUserByUsernameAndGroup(username string, groupName string) bool {
	groupName = strings.ToLower(groupName)
	group := tracker.GetGroupByName(groupName)
	if group == nil {
		return false
	}
	return group.RemoveUserByUsername(username)
}

// Get group by name
func (tracker *PlayerTracker) GetGroupByName(name string) *Group {
	name = strings.ToLower(name)
	for _, group := range tracker.Groups {
		if group.Name == name {
			return group
		}
	}
	return nil
}

// Get user by username
func (tracker *PlayerTracker) GetUserByUsername(username string) *User {
	users := tracker.Users()
	_, foundUser := SearchUsersWithUserCreatedName(users, func(u *User) string { return u.GetUsername() }, username, false, false)
	if foundUser != nil && *foundUser != nil {
		return *foundUser
	}

	return nil
}

// Get user by ID
func (tracker *PlayerTracker) GetUserByID(id string) *User {
	users := tracker.Users()
	for _, user := range users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

// Get list of users from all groups.
// Users.Group contains the group name still.
func (tracker *PlayerTracker) Users() []*User {
	users := []*User{}

	for _, group := range tracker.Groups {
		for _, user := range group.Users {
			users = append(users, user)
		}
	}

	return users
}

// Start the tracker loop
func (tracker *PlayerTracker) Start() {
	logger.Info("Starting tracker")
	tracker.Running = true
	go tracker.Loop()
}

// Stop the tracker loop, updates messenger as well
func (tracker *PlayerTracker) Stop() {
	tracker.Running = false
	// Allows messenger to know to stop
	tracker.Channel <- time.Now()
}

// Returns if the tracker is running
func (tracker *PlayerTracker) IsRunning() bool {
	return tracker.Running
}

// Main loop for the tracker. Update, signal messenger, wait, repeat
func (tracker *PlayerTracker) Loop() {
	for tracker.Running {
		tracker.Update()
		tracker.Channel <- time.Now()
		time.Sleep(time.Duration(tracker.Interval) * time.Second)
	}
}

// Fetches data from Battle Metrics and updates tracker with up to date information
func (tracker *PlayerTracker) Update() {
	// URL: https://api.battlemetrics.com/servers/10519728?include=player
	resp, err := http.Get(fmt.Sprintf("https://api.battlemetrics.com/servers/%s?include=player", tracker.BattleMetricsID))
	if err != nil {
		logger.Error("Failed to fetch Battle Metrics data", "error", err)
		return
	}

	bmRes := BattleMetricsResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&bmRes); err != nil {
		logger.Error("Failed to decode Battle Metrics response", "error", err)
		d := []byte{}
		resp.Body.Read(d)
		return
	}

	// Set server population
	tracker.Online[0] = bmRes.Data.Attributes.Players
	tracker.Online[1] = bmRes.Data.Attributes.MaxPlayers
	tracker.ServerName = bmRes.Data.Attributes.Name

	// Update status of players
	// If this is the first time we see the user, we add them to the tracker
	users := tracker.Users()
	for _, user := range users {
		var player *Player

		// Find player in BattleMetrics
		if user.Status == StatusUnknown && len(user.Usernames) == 0 && user.ID != "" {
			_, player = SearchUsersWithUserCreatedName(bmRes.Included, func(p Player) string { return p.ID }, user.ID, false, false)
		} else {
			_, player = SearchUsersWithUserCreatedName(bmRes.Included, func(p Player) string { return p.Attributes.Name }, user.GetUsername(), false, false)
		}

		// If found, update status
		if player != nil {
			user.SetStatus(StatusOnline)
			user.ID = player.ID
			user.ChangeUsername(player.Attributes.Name)
			continue
		}

		// If not found, set status to offline
		if user.Status == StatusOnline {
			user.SetStatus(StatusOffline)
		}
	}

	return
}
