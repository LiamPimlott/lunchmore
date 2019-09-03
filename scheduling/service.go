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
		log.Printf("error scheduling all 1: %s", err.Error())
		return err
	}

	// Todo: check to see if already scheduled?
	for _, schd := range schdls {
		e, err := c.AddJob(schd.Spec, MatchingJob{
			SchdlngSrvce: s,
			SchdlID:      schd.ID,
		})
		if err != nil {
			log.Printf("error scheduling all 2: %s", err.Error())
			return err
		}
		log.Printf("Entry %d Scheduled: %s\n", e, schd.Spec)
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

// MatchingJob models matching job
type MatchingJob struct {
	SchdlngSrvce *schedulingService
	SchdlID      uint
}

// Run satisfies the cron Job interface
func (j MatchingJob) Run() {

	log.Printf("Running cron job for Schedule ID: %d\n", j.SchdlID)

	// get schedule users
	schdlUsrs, err := j.SchdlngSrvce.GetScheduleUsers(j.SchdlID)
	if err != nil {
		// a way to re-register before killing?
		log.Printf("1 error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	log.Printf("Sched %d Users: %+v\n", j.SchdlID, schdlUsrs)

	if len(schdlUsrs) < 2 {
		log.Printf("Less than 2 users have signed up, exiting Schedule %d.\n", j.SchdlID)
		return
	}

	// run matching algorithm
	mtchs, odd, err := matchRandom(schdlUsrs)
	if err != nil {
		// a way to re-register before killing?
		log.Printf("2 error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	log.Printf("Odd User Out: %+v\n", odd.ID)
	log.Printf("Sched %d Matches: %+v\n", j.SchdlID, mtchs)

	// save matches
	err = j.SchdlngSrvce.SaveMatches(mtchs)
	if err != nil {
		// a way to re-register before killing?
		log.Printf("3 error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	log.Printf("Finished cron job for Schedule ID: %d\n", j.SchdlID)

	// TODO: send match emails
	// TODO: send apology email for remainder
}

// matchRandom returns slice of randomely matched schedule users the odd user out
func matchRandom(su []ScheduleUser) ([]LunchMatch, ScheduleUser, error) {
	lm := []LunchMatch{}

	for len(su) > 1 {
		r1 := rand.Intn(len(su))
		r2 := rand.Intn(len(su))
		for r1 == r2 {
			r2 = rand.Intn(len(su))
		}

		lm = append(lm, LunchMatch{
			UserID1:    su[r1].UserID,
			UserID2:    su[r2].UserID,
			ScheduleID: su[r1].ScheduleID,
		})

		if su[r2] == su[len(su)-1] {
			su[r1], su[len(su)-2] = su[len(su)-2], su[r1]
		} else {
			su[r1], su[len(su)-1] = su[len(su)-1], su[r1]
			su[r2], su[len(su)-2] = su[len(su)-2], su[r2]
		}
		su = su[:len(su)-2]
	}

	if len(su) != 0 {
		return lm, su[0], nil
	}
	return lm, ScheduleUser{}, nil
}
