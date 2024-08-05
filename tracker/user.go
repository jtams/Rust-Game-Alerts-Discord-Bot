package tracker

import "time"

type UserStatus int

const (
	StatusOnline UserStatus = iota
	StatusOffline
	StatusUnknown
	StatusUntracked
)

// Playtime struct to store the start and end time of a user's online time
type Playtime struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// User struct to store user information
type User struct {
	// BattleMetrics ID
	ID string `json:"ID"`

	// List of usernames the user has had. The first username is the name
	// given to the tracker by the command issuer.
	Usernames []string `json:"usernames"`

	Group string `json:"group"`

	// Status of the user, Online, Offline, Unknown or Untracked
	Status UserStatus `json:"status"`

	// List of online times
	OnlineTimes []Playtime `json:"onlineTimes"`
}

// Adds a new username.
func (user *User) ChangeUsername(username string) {
	if user.Usernames == nil {
		user.Usernames = []string{}
	}

	if len(user.Usernames) == 0 {
		user.Usernames = append(user.Usernames, username)
		return
	}

	if user.Usernames[len(user.Usernames)-1] != username {
		user.Usernames = append(user.Usernames, username)
	}
}

// Get current username
func (user *User) GetUsername() string {
	if len(user.Usernames) == 0 {
		return ""
	}
	return user.Usernames[len(user.Usernames)-1]
}

// Set the status of the user. Automatically updates online times
func (user *User) SetStatus(status UserStatus) {
	if user.Status == status {
		return
	}

	if user.Status != StatusOnline && status == StatusOnline {
		user.OnlineTimes = append(user.OnlineTimes, Playtime{StartTime: time.Now()})
	}

	if user.Status == StatusOnline && status == StatusOffline {
		user.OnlineTimes[len(user.OnlineTimes)-1].EndTime = time.Now()
	}

	user.Status = status
}

// Get the last time the user was online
func (user *User) GetLastOnline() time.Time {
	if user.Status == StatusOffline {
		return user.OnlineTimes[len(user.OnlineTimes)-1].EndTime
	} else {
		return time.Now()
	}
}
