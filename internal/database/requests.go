package database

import "forum-go/internal/models"

func (s *service) GetRequests() ([]models.Request, error) {
	rows, err := s.db.Query(`
		SELECT r.*, u.username 
		FROM request r 
		JOIN user u ON r.user_id = u.user_id
		ORDER BY r.creation_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []models.Request
	for rows.Next() {
		var request models.Request
		err := rows.Scan(&request.RequestId, &request.UserId, &request.Status, &request.Content, &request.CreationDate, &request.Username)
		if err != nil {
			return nil, err
		}
		request.FormattedCreationDate = request.CreationDate.Format("2006-01-02 15:04:05")
		requests = append(requests, request)
	}
	return requests, nil
}

func (s *service) CreateRequest(request models.Request) error {
	_, err := s.db.Exec(`
		INSERT INTO request (request_id, user_id, content, creation_date, status) 
		VALUES (?, ?, ?, ?, ?)
	`, request.RequestId, request.UserId, request.Content, request.CreationDate, request.Status)
	return err
}

func (s *service) DeleteRequest(requestId string) error {
	_, err := s.db.Exec(`
		DELETE FROM request WHERE request_id = ?
	`, requestId)
	return err
}

func (s *service) UpdateRequestStatus(requestId, status string) error {
	_, err := s.db.Exec(`
		UPDATE request SET status = ? WHERE request_id = ?
	`, status, requestId)
	return err
}
