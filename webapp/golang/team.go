package xsuportal

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/patrickmn/go-cache"
)

var (
	teamCache            = cache.New(5*time.Minute, 10*time.Minute)
	teamContestStartedAt time.Time
)

func GetTeamByID(teamID int64, db sqlx.Queryer, lock bool) (*Team, error) {
	if teamContestStartedAt.IsZero() {
		sqlx.Get(db, &teamContestStartedAt, "SELECT `contest_starts_at` FROM `contest_config`")
	}

	if teamContestStartedAt.IsZero() {
		var team Team
		query := "SELECT * FROM `teams` WHERE `id` = ?"
		if lock {
			query += " FOR UPDATE"
		}
		err := sqlx.Get(db, &team, query, teamID)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("query team: %w", err)
		}
		return &team, nil
	} else if !teamContestStartedAt.IsZero() && time.Now().Before(teamContestStartedAt) {
		var team Team
		query := "SELECT * FROM `teams` WHERE `id` = ?"
		if lock {
			query += " FOR UPDATE"
		}
		err := sqlx.Get(db, &team, query, teamID)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("query team: %w", err)
		}
		return &team, nil
	}
	if lock {
		var team Team
		query := "SELECT * FROM `teams` WHERE `id` = ? FOR UPDATE"
		err := sqlx.Get(db, &team, query, teamID)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("query team: %w", err)
		}
		return &team, nil
	}
	if val, found := teamCache.Get(strconv.FormatInt(teamID, 10)); found {
		team := val.(Team)
		return &team, nil
	}
	var team Team
	query := "SELECT * FROM `teams` WHERE `id` = ?"
	err := sqlx.Get(db, &team, query, teamID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query team: %w", err)
	}
	teamCache.Set(strconv.FormatInt(teamID, 10), team, cache.DefaultExpiration)
	return &team, nil
}
