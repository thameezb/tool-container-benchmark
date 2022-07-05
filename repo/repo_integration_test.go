//go:build integration_test
// +build integration_test

package repo_test

import (
	"fmt"
	"net/url"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gray.net/tool-container-benchmark/repo"
)

func connect(dbHost, dbPort, dbName, dbUsername, dbPassword string) *sqlx.DB {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&app+name=tool-container-benchmark",
		url.PathEscape(dbUsername),
		url.PathEscape(dbPassword),
		dbHost,
		dbPort,
		url.QueryEscape(dbName))

	return sqlx.MustConnect("mssql", dsn)
}

var _ = Describe("Project Repository Integration", func() {
	var (
		r repo.Interface
	)

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	db := connect(dbHost, dbPort, dbName, dbUsername, dbPassword)

	proj := repo.Project{
		GitlabID:     1,
		Name:         "test",
		TeamGitlabID: 10002,
		Path:         "/test/test",
		AvatarURL:    "test",
		IsActive:     true,
		ToSync:       false,
	}

	BeforeEach(func() {
		r = repo.NewRepository(dbHost, dbPort, dbName, dbUsername, dbPassword)
		Expect(r.Wipe()).To(BeNil())
		Expect(r.CreateProject(&proj)).To(BeNil())
	})

	It("creates a new project", func() {
		var gotName string
		Expect(db.Get(
			&gotName,
			`SELECT TOP 1 name FROM [ContainerBenchmark].[Project] WHERE name=?`,
			"test")).To(BeNil())
		Expect(proj.Name).To(Equal(gotName))
	})

	It("gets a repo", func() {
		p, err := r.GetProject(proj.GitlabID)
		Expect(err).To(BeNil())
		Expect(p.Name).To(Equal(proj.Name))
	})

	It("updates a project", func() {
		proj.AvatarURL = "pre"
		err := r.UpdateProject(&proj)
		Expect(err).To(BeNil())

		var gotURL string
		err = db.Get(
			&gotURL,
			`SELECT TOP 1 avatarURL FROM [ContainerBenchmark].[Project] WHERE Name=?`,
			proj.Name)
		Expect(err).To(BeNil())
		Expect(gotURL).To(Equal("pre"))
	})

	It("marks deleted project as inactive", func() {
		err := r.CreateProject(&repo.Project{
			GitlabID:     100,
			Name:         "toDelete",
			TeamGitlabID: 10002,
			Path:         "/test/toDelete",
			AvatarURL:    "test",
			IsActive:     true,
			ToSync:       false,
		})
		Expect(err).To(BeNil())

		err = r.DeleteProject(100)
		Expect(err).To(BeNil())

		var isActive bool
		err = db.Get(
			&isActive,
			`SELECT TOP 1 isActive FROM [ContainerBenchmark].[Project] WHERE Name=?`,
			"toDelete")
		Expect(err).To(BeNil())
		Expect(isActive).To(BeFalse())
	})

	It("drops table data", func() {
		err := r.Wipe()
		Expect(err).To(BeNil())

		data := []repo.Project{}
		err = db.Select(
			&data,
			`SELECT * FROM [ContainerBenchmark].[Project]`)
		Expect(err).To(BeNil())
		Expect(len(data)).To(Equal(0))
	})
})
