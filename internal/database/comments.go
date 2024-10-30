package database

func (s *service) DeleteComments(postId string) error {
	query := "DELETE FROM Comment WHERE post_id=?"
	_, err := s.db.Exec(query, postId)
	return err
}
