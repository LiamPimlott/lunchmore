package scheduling

import (
	"log"
	"math/rand"

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
	SchdlngSrvce *schedulingService
	SchdlID      uint
}

// Run satisfies the cron Job interface
func (j MatchingJob) Run() {
	// get schedule users
	schdlUsrs, err := j.SchdlngSrvce.GetScheduleUsers(j.SchdlID)
	if err != nil {
		// a way to re-register before killing?
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	// run matching algorithm
	mtchs, _, err := matchRandom(schdlUsrs)
	if err != nil {
		// a way to re-register before killing?
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	// save matches
	err = j.SchdlngSrvce.SaveMatches(mtchs)
	if err != nil {
		// a way to re-register before killing?
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	// TODO: send match emails
	// TODO: send apology email for remainder
}

// NewSchedulingService will return a struct that implements the schedulesService interface
func NewSchedulingService(mail mail.Service, repo Repository) *schedulingService {
	return &schedulingService{
		mail: mail,
		repo: repo,
	}
}

// ScheduleAll schedules a cron job for all Schedules in the database
func (s *schedulingService) ScheduleAll(c *cron.Cron) error {
	schdls, err := s.repo.GetSchedules()
	if err != nil {
		log.Printf("error scheduling all: %s", err.Error())
		return err
	}

	// Todo: check to see if already scheduled?
	for _, schd := range schdls {
		_, err = c.AddJob(schd.Spec, MatchingJob{
			SchdlngSrvce: s,
			SchdlID:      schd.ID,
		})
		if err != nil {
			log.Printf("error scheduling all: %s", err.Error())
			return err
		}
	}
	return nil
}

// GetScheduleUsers return all
func (s *schedulingService) GetScheduleUsers(schdID uint) ([]ScheduleUser, error) {
	su, err := s.repo.GetScheduleUsers(schdID)
	if err != nil {
		log.Printf("error getting schedule users: %s", err.Error())
		return []ScheduleUser{}, err
	}
	return su, nil
}

// SaveMatches
func (s *schedulingService) SaveMatches(lm []LunchMatch) error {
	err := s.repo.SaveLunchMatches(lm)
	if err != nil {
		log.Printf("error saving lunch matches: %s", err.Error())
		return err
	}
	return nil
}

// matchRandom returns slice of randomely matched schedule users the odd user out
func matchRandom(su []ScheduleUser) ([]LunchMatch, ScheduleUser, error) {
	lm := []LunchMatch{}

	m := len(su) / 2

	for i := 0; i < m; i++ {
		r1 := rand.Intn(len(su))
		r2 := rand.Intn(len(su))
		for r1 == r2 {
			r2 = rand.Intn(len(su))
		}

		lm = append(lm, LunchMatch{
			UserID1: su[r1].ID,
			UserID2: su[r2].ID,
		})

		su[r1], su[len(su)-1] = su[len(su)-1], su[r1]
		su[r2], su[len(su)-2] = su[len(su)-2], su[r2]
		su = su[:len(su)-1]
	}

	if len(su) != 0 {
		return lm, su[0], nil
	}

	return lm, ScheduleUser{}, nil
}
