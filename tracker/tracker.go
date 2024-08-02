package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type PlayerTracker struct {
	BattleMetricsID string         `json:"battleMetricsID"`
	ServerName      string         `json:"serverName"`
	Groups          []*Group       `json:"groups"`
	Interval        int            `json:"interval"`
	Running         bool           `json:"running"`
	Channel         chan time.Time `json:"-"`
	Online          [2]int         `json:"online"`
}

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
	tracker.AddNewGroup("Squad")
	tracker.AddNewGroup("Allies")
	tracker.AddNewGroup("Neighbors")
	tracker.AddNewGroup("Enemies")
	tracker.AddNewGroup("Others")

	return tracker
}

func (tracker *PlayerTracker) AddGroup(group *Group) error {
	if tracker.GetGroupByName(group.Name) != nil {
		return errors.New("group with that name already exists")
	}

	tracker.Groups = append(tracker.Groups, group)
	return nil
}

func (tracker *PlayerTracker) AddNewGroup(name string) error {
	newGroup := NewGroup(name)
	return tracker.AddGroup(newGroup)
}

func (tracker *PlayerTracker) RemoveGroup(groupName string) {
	for i, group := range tracker.Groups {
		if group.Name == groupName {
			tracker.Groups = append(tracker.Groups[:i], tracker.Groups[i+1:]...)
			break
		}
	}
}

func (tracker *PlayerTracker) AddUserToGroup(username string, groupName string) error {
	group := tracker.GetGroupByName(groupName)
	if group == nil {
		return errors.New("group not found")
	}

	group.AddUser(username)
	return nil
}

func (tracker *PlayerTracker) RemoveUserByUsername(username string) bool {
	for _, group := range tracker.Groups {
		return group.RemoveUserByUsername(username)
	}

	return false
}

func (tracker *PlayerTracker) GetGroupByName(name string) *Group {
	for _, group := range tracker.Groups {
		if group.Name == name {
			return group
		}
	}
	return nil
}

func (tracker *PlayerTracker) Users() []*User {
	users := []*User{}

	for _, group := range tracker.Groups {
		for _, user := range group.Users {
			users = append(users, user)
		}
	}

	return users
}

func (tracker *PlayerTracker) Start() {
	log.Println("Starting tracker")
	tracker.Running = true
	go tracker.Loop()
}

func (tracker *PlayerTracker) Stop() {
	tracker.Running = false
	// Allows messenger to know to stop
	tracker.Channel <- time.Now()
}

func (tracker *PlayerTracker) IsRunning() bool {
	return tracker.Running
}

func (tracker *PlayerTracker) Loop() {
	for tracker.Running {
		tracker.Update()
		tracker.Channel <- time.Now()
		time.Sleep(time.Duration(tracker.Interval) * time.Second)
	}
}

// Finds a user by username
// Todo: Do search in a single loop
func (tracker *PlayerTracker) GetUserByName(username string, exact bool, casesensitive bool) *User {
	users := tracker.Users()
	var matchedUser *User
	lowerUsername := strings.ToLower(username)

	// Exact match
	for _, user := range users {
		if user.GetUsername() == username {
			return user
		}
	}

	// Case insensitive exact match
	for _, user := range users {
		if strings.ToLower(user.GetUsername()) == lowerUsername {
			return user
		}
	}

	// Match the start of the name
	for _, user := range users {
		if strings.HasPrefix(user.GetUsername(), username) {
			return user
		}
	}

	// Case insensitive match the start of the name
	for _, user := range users {
		if strings.HasPrefix(strings.ToLower(user.GetUsername()), lowerUsername) {
			return user
		}
	}

	// Substring match (more than 5 characters)
	for _, user := range users {
		if len(username) > 5 && strings.Contains(user.GetUsername(), username) {
			return user
		}
	}

	// Case insensitive substring match (more than 5 characters)
	for _, user := range users {
		if len(username) > 5 && strings.Contains(strings.ToLower(user.GetUsername()), lowerUsername) {
			return user
		}
	}

	return matchedUser
}

// Todo: Do search in a single loop
func MatchUserToBattleMetricsPlayer(user *User, bm BattleMetricsResponse, exact bool, casesensitive bool) *Player {
	players := bm.Included
	username := user.GetUsername()
	lowerUsername := strings.ToLower(username)

	for _, player := range players {
		if player.ID == user.ID {
			return &player
		}
	}

	// Exact match
	for _, player := range players {
		if player.Attributes.Name == username {
			return &player
		}
	}

	// Case insensitive exact match
	for _, player := range players {
		if strings.ToLower(player.Attributes.Name) == lowerUsername {
			return &player
		}
	}

	// Match the start of the name
	for _, player := range players {
		if strings.HasPrefix(player.Attributes.Name, username) {
			return &player
		}
	}

	// Case insensitive match the start of the name
	for _, player := range players {
		if strings.HasPrefix(strings.ToLower(player.Attributes.Name), lowerUsername) {
			return &player
		}
	}

	// Substring match (more than 5 characters)
	for _, player := range players {
		if len(username) > 5 && strings.Contains(player.Attributes.Name, username) {
			return &player
		}
	}

	// Case insensitive substring match (more than 5 characters)
	for _, player := range players {
		if len(username) > 5 && strings.Contains(strings.ToLower(player.Attributes.Name), lowerUsername) {
			return &player
		}
	}

	return nil
}

func (tracker *PlayerTracker) GetUserByID(id string) *User {
	users := tracker.Users()
	for _, user := range users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

func (tracker *PlayerTracker) Update() {
	// URL: https://api.battlemetrics.com/servers/10519728?include=player
	resp, err := http.Get(fmt.Sprintf("https://api.battlemetrics.com/servers/%s?include=player", tracker.BattleMetricsID))
	if err != nil {
		log.Println(err)
		return
	}

	bmRes := BattleMetricsResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&bmRes); err != nil {
		log.Println("Failed to decode Battle Metrics response")
		log.Println(err)
		d := []byte{}
		resp.Body.Read(d)
		log.Println(d)
		return
	}

	tracker.Online[0] = bmRes.Data.Attributes.Players
	tracker.Online[1] = bmRes.Data.Attributes.MaxPlayers
	tracker.ServerName = bmRes.Data.Attributes.Name

	// Update status of players
	// If this is the first time we see the user, we add them to the tracker
	users := tracker.Users()
	for _, user := range users {
		var player *Player

		// Find player in BattleMetrics
		player = MatchUserToBattleMetricsPlayer(user, bmRes, false, false)

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

func (tracker *PlayerTracker) Save() {

}
