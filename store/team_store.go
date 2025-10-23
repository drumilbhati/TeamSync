package store

/*
	APIs needed:
	GET:
	GetTeamByID
	GetTeamsByUserID

	POST:
	CreateTeam

	PUT:
	UpdateTeamByID

	DELETE:
	DeleteTeamByID
*/

import (
	"github.com/drumilbhati/teamsync/models"
)

/*
Given a team_id return the entire team
*/
func (s *Store) GetTeamByID(team_id int) (*models.Team, error) {
	var team models.Team
	err := s.db.QueryRow(
		"SELECT team_id, team_name, team_leader_id, created_at FROM teams WHERE team_id = $1",
		team_id,
	).Scan(&team.TeamID, &team.TeamName, &team.TeamLeaderID, &team.CreatedAt)

	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(
		"SELECT member_id, user_id, role FROM members WHERE team_id = $1", team_id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var member models.Member
		if err := rows.Scan(&member.MemberID, &member.UserID, &member.Role); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	team.Members = members
	return &team, nil
}

/*
Given a user_id return all teams the user is part
handles for team_leader_id as well, since team_leader_id references user_id
*/
func (s *Store) GetTeamsByUserID(user_id int) ([]models.Team, error) {
	var teams = make([]models.Team, 0)
	rows, err := s.db.Query(
		`SELECT team_id FROM members where user_id = $1
		UNION
		SELECT team_id FROM teams where team_leader_id = $1`, user_id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.TeamID); err != nil {
			return nil, err
		}
		full_team, err := s.GetTeamByID(team.TeamID)
		if err != nil {
			return nil, err
		}
		teams = append(teams, *full_team)
	}
	return teams, nil
}

/*
Given a team create a team in teams table
*/
func (s *Store) CreateTeam(t *models.Team) error {
	err := s.db.QueryRow(
		`INSERT INTO teams (team_name, team_leader_id)
		VALUES ($1, $2)
		RETURNING team_id, created_at`,
		t.TeamName, t.TeamLeaderID,
	).Scan(&t.TeamID, &t.CreatedAt)

	return err
}

/*
Given a team_id update its details
*/
func (s *Store) UpdateTeamByID(team_id int, t *models.Team) error {
	_, err := s.db.Exec(
		`UPDATE teams
		SET team_name = $1, team_leader_id = $2 
		WHERE team_id = $3`,
		t.TeamName, t.TeamLeaderID, team_id,
	)
	return err
}

/*
Given a team_id delete it from database
*/
func (s *Store) DeleteTeamByID(team_id int) error {
	_, err := s.db.Exec(
		"DELETE FROM teams WHERE team_id = $1",
		team_id,
	)
	return err
}
