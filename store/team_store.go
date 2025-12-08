package store

/*
	APIs needed:
	GET:
	GetTeamByID
	GetTeamsByTeamLeaderID
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
func (s *Store) GetTeamsByTeamLeaderID(team_leader_id int) ([]models.Team, error) {
	var teams = make([]models.Team, 0)
	rows, err := s.db.Query(
		`SELECT t.team_id, t.team_name, t.team_leader_id, t.created_at
			FROM teams t WHERE t.team_leader_id = $1
		`,
		team_leader_id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.TeamID, &t.TeamName, &t.TeamLeaderID, &t.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (s *Store) GetTeamsByUserID(user_id int) ([]models.Team, error) {
	teams := []models.Team{}
	rows, err := s.db.Query(
		`SELECT DISTINCT t.team_id, t.team_name, t.team_leader_id, t.created_at
		FROM teams t
		LEFT JOIN members m ON t.team_id = m.team_id
		WHERE m.user_id = $1 OR t.team_leader_id = $1`,
		user_id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t models.Team
		t.Members = []models.Member{}
		if err := rows.Scan(&t.TeamID, &t.TeamName, &t.TeamLeaderID, &t.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
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
