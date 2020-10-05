package xsuportal

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/patrickmn/go-cache"
)

var (
	contestantCache            = cache.New(5*time.Minute, 10*time.Minute)
	contestantContestStartedAt time.Time
)

func GetContestantByID(contestantID string, db sqlx.Queryer, lock bool) (Contestant, error) {
	if contestantContestStartedAt.IsZero() {
		err := sqlx.Get(db, &contestantContestStartedAt, "SELECT `contest_starts_at` FROM `contest_config`")
		if err != nil {
			fmt.Errorf("SELECT `contest_starts_at` FROM `contest_config`: %w", err)
		}
	}

	if contestantContestStartedAt.IsZero() {
		var contestant Contestant
		query := "SELECT * FROM `contestants` WHERE `id` = ? LIMIT 1"
		if lock {
			query += " FOR UPDATE"
		}
		err := sqlx.Get(db, &contestant, query, contestantID)
		return contestant, err
	} else if !contestantContestStartedAt.IsZero() && time.Now().Before(contestantContestStartedAt) {
		var contestant Contestant
		query := "SELECT * FROM `contestants` WHERE `id` = ? LIMIT 1"
		if lock {
			query += " FOR UPDATE"
		}
		err := sqlx.Get(db, &contestant, query, contestantID)
		return contestant, err
	}
	if lock {
		var contestant Contestant
		query := "SELECT * FROM `contestants` WHERE `id` = ? FOR UPDATE"
		err := sqlx.Get(db, &contestant, query, contestantID)
		return contestant, err
	}
	if val, found := contestantCache.Get(contestantID); found {
		contestant := val.(Contestant)
		return contestant, nil
	}
	var contestant Contestant
	query := "SELECT * FROM `contestants` WHERE `id` = ?"
	err := sqlx.Get(db, &contestant, query, contestantID)
	if err != nil {
		return contestant, err
	}
	contestantCache.Set(contestantID, contestant, cache.DefaultExpiration)
	return contestant, nil
}
