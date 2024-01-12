package job

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type JobScheduler struct {
	scheduler gocron.Scheduler
	dcSession *discordgo.Session
}

type Job struct {
	Name          string
	JobDefinition gocron.JobDefinition
	Task          gocron.Task
}

func (s *JobScheduler) getJobs() []Job {
	return []Job{
		{
			Name: "test",
			JobDefinition: gocron.DurationJob(
				10 * time.Second,
			),
			Task: gocron.NewTask(
				s.hello,
				"hello",
				1,
			),
		},
	}
}

func Initialize(dcSession *discordgo.Session) *JobScheduler {
	logger := zap.L().Sugar()
	defer logger.Sync()

	logger.Debugln("Initializing job scheduler")

	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Panicf("error creating scheduler: %v", err)
	}

	return &JobScheduler{
		scheduler: s,
		dcSession: dcSession,
	}
}

func (js *JobScheduler) Start() {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Debugln("Starting job scheduler")

	// Add jobs
	for _, job := range js.getJobs() {
		j, err := js.scheduler.NewJob(job.JobDefinition, job.Task)
		if err != nil {
			logger.Errorf("error creating job: %v", err)
		}
		logger.Debugf("Job %v created", j.ID())
	}

	js.scheduler.Start()
}

func (s *JobScheduler) Stop() {
	logger := zap.L().Sugar()
	defer logger.Sync()

	logger.Debugln("Stopping job scheduler")
	err := s.scheduler.Shutdown()
	if err != nil {
		logger.Errorf("error shutting down scheduler: %v", err)
	}
}
