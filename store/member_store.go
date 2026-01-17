package store

import (
	"github.com/drumilbhati/teamsync/models"
)

/*
	APIs
	GET:
	GetMemberByID
	GetMembersByTeamID

	POST:
	CreateMember

	PUT:
	UpdateMemberByID

	DELETE:
	DeleteMemberByID
*/

func (s *Store) GetMemberByID(member_id int) (*models.Member, error) {
	var member models.Member
	err := s.db.QueryRow(
		`SELECT m.member_id, m.user_id, m.team_id, m.role, u.user_name, u.email 
		FROM members m
		JOIN users u ON m.user_id = u.user_id
		WHERE m.member_id = $1`,
		member_id,
	).Scan(&member.MemberID, &member.UserID, &member.TeamID, &member.Role, &member.UserName, &member.Email)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (s *Store) GetMembersByTeamID(team_id int) ([]models.Member, error) {
	var members []models.Member
	rows, err := s.db.Query(
		`SELECT m.member_id, m.user_id, m.team_id, m.role, u.user_name, u.email 
		FROM members m
		JOIN users u ON m.user_id = u.user_id
		WHERE m.team_id = $1`,
		team_id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var member models.Member
		if err := rows.Scan(&member.MemberID, &member.UserID, &member.TeamID, &member.Role, &member.UserName, &member.Email); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (s *Store) CreateMember(m *models.Member) error {
	err := s.db.QueryRow(
		`INSERT INTO members (user_id, team_id, role)
		VALUES ($1, $2, $3)
		RETURNING member_id, created_at`,
		m.UserID, m.TeamID, m.Role,
	).Scan(&m.MemberID, &m.CreatedAt)

	return err
}

func (s *Store) UpdateMemberByID(member_id int, m *models.Member) error {
	_, err := s.db.Exec(
		`UPDATE members
		SET role = $1
		WHERE member_id = $2`,
		m.Role, member_id,
	)

	return err
}

func (s *Store) DeleteMemberByID(member_id int) error {
	_, err := s.db.Exec(
		"DELETE FROM members WHERE member_id = $1",
		member_id,
	)
	return err
}

func (s *Store) IsTeamMember(userID int, teamID int) (bool, error) {
	var exists bool
	err := s.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM members WHERE user_id = $1 AND team_id = $2
		) OR EXISTS(
			SELECT 1 FROM teams WHERE team_leader_id = $1 AND team_id = $2
		)`,
		userID, teamID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}
