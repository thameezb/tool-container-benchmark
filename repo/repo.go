//go:generate mockgen -destination mock_repo/repo.go gray.net/tool-container-benchmark/repo Interface

package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	// MSSQL driver.
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type Interface interface {
	CreateProject(p *Project) error
	GetProject(pID int) (Project, error)
	DeleteProject(projectID int) error
	UpdateProject(p *Project) error
	Wipe() error
}

type repoImpl struct {
	db *sqlx.DB
}

// NewRepository configures a new database connection and returns the associated repository.
func NewRepository(host, port, dbName, username, password string) Interface {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&app+name=tool-container-benchmark",
		url.PathEscape(username),
		url.PathEscape(password),
		host,
		port,
		url.QueryEscape(dbName))
	db := sqlx.MustConnect("mssql", dsn)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &repoImpl{db}
}

func (r *repoImpl) CreateProject(p *Project) error {
	return insert(r.db, "ContainerBenchmark.Project", p)
}

func (r *repoImpl) GetProject(pID int) (Project, error) {
	var p Project
	query := `
		SELECT gitlabId, name, teamGitlabId, path, avatarURL, isActive, toSync
		FROM [ContainerBenchmark].[Project]
	  	WHERE [gitlabId] = ?`
	err := r.db.Get(&p, query, pID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return Project{}, nil
	}

	return p, err
}

func (r *repoImpl) DeleteProject(projectID int) error {
	query := `
	UPDATE [ContainerBenchmark].[Project]
	SET
		[isActive]=0
	WHERE [gitlabId]=:gId`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"gId": projectID,
	})
	return err
}

func (r *repoImpl) UpdateProject(p *Project) error {
	query := `
	UPDATE [ContainerBenchmark].[Project]
	SET
		[name]=:n,
		[path]=:p,
		[avatarURL]=:au,
		[teamGitlabId]=:tGId,
		[toSync]=0
	WHERE [gitlabId]=:gId`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"n":    p.Name,
		"gId":  p.GitlabID,
		"p":    p.Path,
		"au":   p.AvatarURL,
		"tGId": p.TeamGitlabID,
	})
	return err
}

func (r *repoImpl) Wipe() error {
	_, err := r.db.Exec("DELETE FROM [ContainerBenchmark].[Project]")
	return err
}
