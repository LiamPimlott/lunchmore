package scheduling

import (
	"fmt"
	"log"
	"strings"

	"github.com/robfig/cron"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/mail"
	"github.com/LiamPimlott/lunchmore/organizations"
	"github.com/LiamPimlott/lunchmore/users"
)

// Service interface to schedules service
type Service interface {
	CreateSchedule(sched ScheduleRequest, claimedID uint) (Schedule, error)
	ScheduleAll() error
	ScheduleJob(schedID uint) error
	GetScheduleUsers(schedID uint) ([]ScheduleUser, error)
	SaveMatches(lm []LunchMatch) error
}

type schedulingService struct {
	cron  *cron.Cron
	mail  mail.Service
	users users.Service
	orgs  organizations.Service
	repo  Repository
}

// NewSchedulingService will return a struct that implements the schedulesService interface
func NewSchedulingService(
	cron *cron.Cron,
	mail mail.Service,
	users users.Service,
	orgs organizations.Service,
	repo Repository,
) *schedulingService {
	return &schedulingService{
		cron,
		mail,
		users,
		orgs,
		repo,
	}
}

// CreateSchedule creates a new schedule
func (s *schedulingService) CreateSchedule(schedReq ScheduleRequest, claimedID uint) (Schedule, error) {

	org, err := s.orgs.GetByID(schedReq.OrgID)
	if err != nil {
		log.Printf("error getting organization: %s", err.Error())
		return Schedule{}, err
	}

	if org.AdminID != claimedID {
		log.Printf("error creating schedule: admin id does not equal claimed id")
		return Schedule{}, errs.ErrForbidden
	}

	days := strings.Join(schedReq.Days, ",")
	sched := Schedule{
		OrgID: schedReq.OrgID,
		Spec:  fmt.Sprintf("0 0 ? * %s", days),
	}

	sched, err = s.repo.CreateSchedule(sched)
	if err != nil {
		log.Printf("error creating schedule: %s", err.Error())
		return Schedule{}, err
	}

	err = s.ScheduleJob(sched.ID)
	if err != nil {
		log.Printf("error scheduling job: %s", err.Error())
		return Schedule{}, err
	}

	return sched, nil
}

// ScheduleAll schedules a cron job for all Schedules in the database
func (s *schedulingService) ScheduleAll() error {
	schdls, err := s.repo.GetSchedules()
	if err != nil {
		log.Printf("error scheduling all 1: %s", err.Error())
		return err
	}

	// Todo: check to see if already scheduled?
	for _, schd := range schdls {
		e, err := s.cron.AddJob(schd.Spec, MatchingJob{
			SchdlngSrvce: s,
			SchdlID:      schd.ID,
		})
		if err != nil {
			log.Printf("error scheduling all 2: %s", err.Error())
			return err
		}
		log.Printf("org %d scheduled (%s) as entry %d \n", schd.OrgID, schd.Spec, e)
	}
	return nil
}

// ScheduleJob schedules a cron job with the given schedule id
func (s *schedulingService) ScheduleJob(schedID uint) error {
	sched, err := s.repo.GetByID(schedID)
	if err != nil {
		log.Printf("error getting schedule by id: %s", err.Error())
		return err
	}

	e, err := s.cron.AddJob(sched.Spec, MatchingJob{
		SchdlngSrvce: s,
		SchdlID:      sched.ID,
	})
	if err != nil {
		log.Printf("error adding job to cron: %s", err.Error())
		return err
	}
	log.Printf("org %d scheduled (%s) as entry %d \n", sched.OrgID, sched.Spec, e)

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
