package benchmark

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"gray.net/tool-container-benchmark/repo"
)

type Interface interface {
	Run(numberOfEvents int)
	createProject(gitlabID int) error
	modifyProject(gitlabID int) error
	deleteProject(gitlabID int) error
}

type Benchmark struct {
	repo             repo.Interface
	time             Time
	numberOfProjects int
}

func NewBenchmark(repo repo.Interface, numberOfProjects int) Interface {
	return &Benchmark{
		repo:             repo,
		numberOfProjects: numberOfProjects,
	}
}

func (bm *Benchmark) Run(numberOfEvents int) {
	var err error

	log.Info("Wiping DB")
	if err := bm.repo.Wipe(); err != nil {
		panic(err)
	}

	log.Info("Creating projects")
	bm.time.projectCreationStartTime = time.Now()
	for i := 0; i < bm.numberOfProjects; i++ {
		err := bm.createProject(i)
		if err != nil {
			log.Errorf("failed to create project %d - %s", i, err.Error())
		}
	}
	bm.time.projectCreationEndTime = time.Now()
	log.Info("Projects creation complete")

	log.Info("STARTING BENCHMARK")
	bm.time.eventStartTime = time.Now()
	for i := 0; i < numberOfEvents; i++ {
		choice := rand.Intn(3)
		projectID := rand.Intn(bm.numberOfProjects)

		switch choice {
		case 0:
			err = bm.createProject(bm.numberOfProjects + i)
		case 1:
			err = bm.modifyProject(projectID)
		case 2:
			err = bm.deleteProject(projectID)
		}

		if err != nil {
			log.Errorf("failed to run event type %d (event number %d) for project %d - %s", choice, i, projectID, err.Error())
		}
	}
	bm.time.eventEndTime = time.Now()
	log.Info("BENCHMARK COMPLETE")

	bm.displayTimes()
}

func (bm *Benchmark) createProject(gitlabID int) error {
	pr := repo.Project{
		GitlabID:     gitlabID,
		Name:         GenerateRandomString(15),
		TeamGitlabID: rand.Intn(1000),
		Path:         GenerateRandomString(30),
		AvatarURL:    GenerateRandomString(30),
	}

	return bm.repo.CreateProject(&pr)
}

func (bm *Benchmark) modifyProject(gitlabID int) error {
	pr, err := bm.repo.GetProject(gitlabID)
	if err != nil {
		return err
	}

	pr.Name = GenerateRandomString(15)
	pr.Path = GenerateRandomString(30)
	pr.AvatarURL = GenerateRandomString(30)

	return bm.repo.UpdateProject(&pr)
}

func (bm *Benchmark) deleteProject(gitlabID int) error {
	return bm.repo.DeleteProject(gitlabID)
}

func (bm *Benchmark) displayTimes() {
	log.Infof(`
==========================================================================	
RESULTS
Total Runtime: %s 
==========================================================================
Time taken to create projects: %s
Time taken to run events: %s
==========================================================================
	`,
		bm.time.eventEndTime.Sub(bm.time.projectCreationStartTime),
		bm.time.projectCreationEndTime.Sub(bm.time.projectCreationStartTime),
		bm.time.eventEndTime.Sub(bm.time.eventStartTime),
	)
}
