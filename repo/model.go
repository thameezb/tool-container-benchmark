package repo

type Project struct {
	GitlabID     int    `db:"gitlabId"`
	Name         string `db:"name"`
	TeamGitlabID int    `db:"teamGitlabId"`
	Path         string `db:"path"`
	AvatarURL    string `db:"avatarURL"`
	IsActive     bool   `db:"isActive"`
	ToSync       bool   `db:"toSync"`
}
