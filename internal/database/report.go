package database

import "forum-go/internal/models"

func (s *service) CreateReport(report models.Report) error {
	// Insert a new report
	_, err := s.db.Exec(`
		INSERT INTO Report (report_id, user_id, post_id, creation_date, content, reason, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?);`, report.ReportId, report.UserId, report.PostId, report.CreationDate, report.Content, report.Reason, report.Status)
	return err
}

func (s *service) GetReports() ([]models.Report, error) {
	// Get all reports
	rows, err := s.db.Query(`
		SELECT r.report_id, r.user_id, r.post_id, r.creation_date, r.content, r.reason, r.status, u.username
		FROM Report r
		JOIN User u ON r.user_id = u.user_id
		ORDER BY r.creation_date DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var report models.Report
		var username string
		err := rows.Scan(&report.ReportId, &report.UserId, &report.PostId, &report.CreationDate, &report.Content, &report.Reason, &report.Status, &username)
		if err != nil {
			return nil, err
		}
		report.Username = username
		report.FormattedCreationDate = report.CreationDate.Format("2006-01-02 15:04:05")
		reports = append(reports, report)
	}
	return reports, nil
}

func (s *service) UpdateReportStatus(reportId, status string) error {
	// Update the status of a report
	_, err := s.db.Exec(`
		UPDATE Report SET status = ? WHERE report_id = ?;`, status, reportId)
	return err
}
