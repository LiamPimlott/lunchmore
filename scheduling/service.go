package scheduling

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/robfig/cron"

	"github.com/LiamPimlott/lunchmore/mail"
	"github.com/LiamPimlott/lunchmore/users"
)

// Service interface to schedules service
type Service interface {
	ScheduleAll() error
	GetScheduleUsers(schedID uint) ([]ScheduleUser, error)
	SaveMatches(lm []LunchMatch) error
}

type schedulingService struct {
	mail  mail.Service
	users users.Service
	repo  Repository
}

// NewSchedulingService will return a struct that implements the schedulesService interface
func NewSchedulingService(mail mail.Service, users users.Service, repo Repository) *schedulingService {
	return &schedulingService{
		mail:  mail,
		users: users,
		repo:  repo,
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

// sendBasicMatchEmail
func (s *schedulingService) sendBasicMatchEmail(ufn, fn, ln string, rcpts []string) error {
	msg := fmt.Sprintf("Hi %s, you have been matched with %s %s. Have fun! :)", ufn, fn, ln)
	err := s.mail.SendText(rcpts, "New lunch match!", msg)
	if err != nil {
		log.Printf("error sending match email: %s", err.Error())
		return err
	}
	return nil
}

// SendMatchEmails
func (s *schedulingService) sendMatchEmails(ms []LunchMatch) error {
	usrIDs := []uint{}
	for _, m := range ms {
		usrIDs = append(usrIDs, m.UserID1, m.UserID2)
	}

	usrs, err := s.users.GetUsersMap(usrIDs)
	if err != nil {
		log.Printf("error sending match emails: %s", err.Error())
		return err
	}

	for _, m := range ms {
		usr1 := usrs[m.UserID1]
		usr2 := usrs[m.UserID2]
		s.sendBasicMatchEmail(usr1.FirstName, usr2.FirstName, usr2.LastName, []string{usr1.Email})
		log.Printf("Match email sent to User %d of Schedule: %d\n", usr1.ID, m.ScheduleID)
		s.sendBasicMatchEmail(usr2.FirstName, usr1.FirstName, usr1.LastName, []string{usr2.Email})
		log.Printf("Match email sent to User %d of Schedule: %d\n", usr2.ID, m.ScheduleID)
	}
	return nil
}

// SendOddOutEmail
func (s *schedulingService) sendOddOutEmail(su ScheduleUser) error {
	u, err := s.users.GetByID(su.UserID, su.UserID)
	if err != nil {
		log.Printf("error sending odd out email: %s", err.Error())
		return err
	}

	msg := fmt.Sprintf("Sorry %s, your the odd one out this time. :(", u.FirstName)
	err = s.mail.SendText([]string{u.Email}, "Were sorry!", msg)
	if err != nil {
		log.Printf("error sending match email: %s", err.Error())
		return err
	}

	log.Printf("Odd-Out email sent to User %+v of Schedule: %d\n", u.ID, su.ScheduleID)
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
	// TODO: a way to re-register before killing?

	schdlUsrs, err := j.SchdlngSrvce.GetScheduleUsers(j.SchdlID)
	if err != nil {
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	} else if len(schdlUsrs) < 2 {
		log.Printf("Less than 2 users have signed up, exiting Schedule %d.\n", j.SchdlID)
		return
	}

	mtchs, odd, err := matchRandom(schdlUsrs)
	if err != nil {
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	err = j.SchdlngSrvce.SaveMatches(mtchs)
	if err != nil {
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	err = j.SchdlngSrvce.sendMatchEmails(mtchs)
	if err != nil {
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	err = j.SchdlngSrvce.sendOddOutEmail(odd)
	if err != nil {
		log.Printf("error running matching job for schedule %d: %s", j.SchdlID, err.Error())
		panic(err)
	}

	log.Printf("Finished cron job for Schedule ID: %d\n", j.SchdlID)
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
