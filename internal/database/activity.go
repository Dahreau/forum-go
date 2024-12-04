package database

import (
	"forum-go/internal/models"
)

func (s *service) GetActivities(user models.User) ([]models.Activity, error) {
	activities := make([]models.Activity, 0)
	query := `
		SELECT 
			a.activity_id, a.user_id, a.action_user_id, a.action_type, a.post_id, a.comment_id, a.creation_date, a.details, a.is_read,
			u.username AS action_user_name
		FROM 
			Activity a
		JOIN 
			User u ON a.action_user_id = u.user_id
		WHERE 
			a.user_id=? 
		ORDER BY 
			a.creation_date DESC`
	rows, err := s.db.Query(query, user.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var activity models.Activity
		var actionUserName string
		err := rows.Scan(
			&activity.ActivityId,
			&activity.UserId,
			&activity.ActionUserId,
			&activity.ActionType,
			&activity.PostId,
			&activity.CommentId,
			&activity.CreationDate,
			&activity.Details,
			&activity.IsRead,
			&actionUserName)
		if err != nil {
			return nil, err
		}
		activity.ActionUsername = actionUserName
		activity.FormattedCreationDate = activity.CreationDate.Format("Jan 02, 2006 - 15:04:05")
		activities = append(activities, activity)
	}
	return activities, nil
}

func (s *service) CreateActivity(activity models.Activity) error {
	query := "INSERT INTO Activity (activity_id, user_id, action_user_id, action_type, post_id, comment_id, creation_date, details, is_read) VALUES (?,?,?,?,?,?,?,?,?)"
	_, err := s.db.Exec(query, activity.ActivityId, activity.UserId, activity.ActionUserId, activity.ActionType, activity.PostId, activity.CommentId, activity.CreationDate, activity.Details, &activity.IsRead)
	return err
}

func (s *service) UpdateActivity(activity models.Activity) error {
	query := "UPDATE Activity SET user_id=?, action_user_id=?, action_type=?, post_id=?, comment_id=?, creation_date=?, details=?, is_read=? WHERE activity_id=?"
	_, err := s.db.Exec(query, activity.UserId, activity.ActionUserId, activity.ActionType, activity.PostId, activity.CommentId, activity.CreationDate, activity.Details, activity.IsRead, activity.ActivityId)
	return err
}

func (s *service) ReadActivites(userId string) error {
	query := "UPDATE Activity SET is_read=1 WHERE user_id=?"
	_, err := s.db.Exec(query, userId)
	return err
}
