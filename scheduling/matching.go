package scheduling

import (
	"log"
	"math/rand"
)

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
