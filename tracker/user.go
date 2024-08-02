package tracker

import "time"

type UserStatus int

const (
	StatusOnline UserStatus = iota
	StatusOffline
	StatusUnknown
	StatusUntracked
)

type Playtime struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type User struct {
	ID          string     `json:"ID"`
	Usernames   []string   `json:"username"`
	Group       string     `json:"group"`
	Status      UserStatus `json:"status"`
	OnlineTimes []Playtime `json:"onlineTimes"`
}

func (user *User) ChangeUsername(username string) {
	if user.Usernames == nil {
		user.Usernames = []string{}
	}

	if user.Usernames[len(user.Usernames)-1] != username {
		user.Usernames = append(user.Usernames, username)
	}
}

func (user *User) GetUsername() string {
	return user.Usernames[len(user.Usernames)-1]
}

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

func (user *User) GetLastOnline() time.Time {
	if user.Status == StatusOffline {
		return user.OnlineTimes[len(user.OnlineTimes)-1].EndTime
	} else {
		return time.Now()
	}
}
