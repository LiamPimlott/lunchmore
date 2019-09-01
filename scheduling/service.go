package scheduling

import (
	"log"

	"github.com/robfig/cron"

	"github.com/LiamPimlott/lunchmore/mail"
)

// Service interface to schedules service
type Service interface {
	ScheduleAll() error
	GetScheduleUsers(schedID uint) ([]ScheduleUser, error)
	SaveMatches(lm []LunchMatch) error
}

type schedulingService struct {
	mail mail.Service
	repo Repository
}

// MatchingJob
type MatchingJob struct {
	SchedulingService *schedulingService
	ScheduleID        uint
}

// Run satisfies the cron Job interface
func (j *MatchingJob) Run() error {
	// get schedule users

	// run matching algorithm

	// save matches

	// send emails

	return nil
}

// NewSchedulingService will return a struct that implements the schedulesService interface
func NewSchedulingService(cron cron.Cron, mail mail.Service, repo Repository) *schedulingService {
	return &schedulingService{
		mail: mail,
		repo: repo,
	}
}

// ScheduleAll schedules a cron job for all Schedules in the database
func (s *schedulingService) ScheduleAll(c *cron.Cron) error {

	scheds, err := s.repo.GetSchedules()
	if err != nil {
		log.Printf("error scheduling all: %s", err.Error())
		return err
	}

	return nil
}

// GetScheduleUsers return all
func (s *schedulingService) GetScheduleUsers(schedID uint) ([]ScheduleUser, error) {
	return nil, nil
}

// SaveMatches
func (s *schedulingService) SaveMatches(lm []LunchMatch) error {
	return nil
}
