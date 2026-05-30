package subjects

import (
	"context"
	"database/sql"
	"fmt"

	subjectsPb "lms-class-service/pb/subjects"
)

// SubjectRepository struct
type SubjectRepository struct {
	Db *sql.DB
}

// List subjects with pagination
func (r *SubjectRepository) List(ctx context.Context, limit, offset uint32, keyword, orderBy, sort string) ([]*subjectsPb.Subject, uint32, error) {
	query := `
		SELECT id, university_id, university_name, faculty_id, faculty_name, 
			programme_id, programme_name, code, name, sks, default_semester,
			COALESCE(updated_by::text, ''), updated_at, created_at
		FROM subjects
		WHERE name ILIKE $1 OR code ILIKE $1
		ORDER BY ` + orderBy + ` ` + sort + `
		LIMIT $2 OFFSET $3
	`

	searchKeyword := "%" + keyword + "%"
	rows, err := r.Db.QueryContext(ctx, query, searchKeyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subjects []*subjectsPb.Subject
	for rows.Next() {
		var subject subjectsPb.Subject
		err := rows.Scan(
			&subject.Id, &subject.UniversityId, &subject.UniversityName,
			&subject.FacultyId, &subject.FacultyName, &subject.ProgrammeId,
			&subject.ProgrammeName, &subject.Code, &subject.Name,
			&subject.Sks, &subject.DefaultSemester, &subject.UpdatedBy,
			&subject.UpdatedAt, &subject.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subjects = append(subjects, &subject)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM subjects WHERE name ILIKE $1 OR code ILIKE $1`
	var count uint32
	err = r.Db.QueryRowContext(ctx, countQuery, searchKeyword).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return subjects, count, nil
}

// Get subject by ID with topics
func (r *SubjectRepository) Get(ctx context.Context, id string) (*subjectsPb.Subject, error) {
	query := `
		SELECT id, university_id, university_name, faculty_id, faculty_name,
			programme_id, programme_name, code, name, sks, default_semester,
			COALESCE(updated_by::text, ''), updated_at, created_at
		FROM subjects WHERE id = $1
	`

	var subject subjectsPb.Subject
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&subject.Id, &subject.UniversityId, &subject.UniversityName,
		&subject.FacultyId, &subject.FacultyName, &subject.ProgrammeId,
		&subject.ProgrammeName, &subject.Code, &subject.Name,
		&subject.Sks, &subject.DefaultSemester, &subject.UpdatedBy,
		&subject.UpdatedAt, &subject.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get topics
	topicQuery := `
		SELECT id, subject_id, name, COALESCE(updated_by::text, ''), updated_at, created_at
		FROM topic_subjects WHERE subject_id = $1
	`
	rows, err := r.Db.QueryContext(ctx, topicQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var topic subjectsPb.TopicSubject
		err := rows.Scan(&topic.Id, &topic.SubjectId, &topic.Name, &topic.UpdatedBy, &topic.UpdatedAt, &topic.CreatedAt)
		if err != nil {
			return nil, err
		}
		subject.Topics = append(subject.Topics, &topic)
	}

	return &subject, nil
}

// Create subject
func (r *SubjectRepository) Create(ctx context.Context, subject *subjectsPb.Subject) error {
	query := `
		INSERT INTO subjects (university_id, university_name, faculty_id, faculty_name,
			programme_id, programme_name, code, name, sks, default_semester, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, updated_at, created_at
	`

	return r.Db.QueryRowContext(ctx, query,
		subject.UniversityId, subject.UniversityName, subject.FacultyId,
		subject.FacultyName, subject.ProgrammeId, subject.ProgrammeName,
		subject.Code, subject.Name, subject.Sks, subject.DefaultSemester,
		subject.UpdatedBy,
	).Scan(&subject.Id, &subject.UpdatedAt, &subject.CreatedAt)
}

// CreateTopic creates a topic for a subject
func (r *SubjectRepository) CreateTopic(ctx context.Context, topic *subjectsPb.TopicSubject) error {
	query := `
		INSERT INTO topic_subjects (subject_id, name, updated_by)
		VALUES ($1, $2, $3)
		RETURNING id, updated_at, created_at
	`

	return r.Db.QueryRowContext(ctx, query,
		topic.SubjectId, topic.Name, topic.UpdatedBy,
	).Scan(&topic.Id, &topic.UpdatedAt, &topic.CreatedAt)
}

// Update subject
func (r *SubjectRepository) Update(ctx context.Context, subject *subjectsPb.Subject) error {
	query := `
		UPDATE subjects SET 
			university_id = $1, university_name = $2, faculty_id = $3, faculty_name = $4,
			programme_id = $5, programme_name = $6, code = $7, name = $8,
			sks = $9, default_semester = $10, updated_by = $11, updated_at = NOW()
		WHERE id = $12
	`

	result, err := r.Db.ExecContext(ctx, query,
		subject.UniversityId, subject.UniversityName, subject.FacultyId,
		subject.FacultyName, subject.ProgrammeId, subject.ProgrammeName,
		subject.Code, subject.Name, subject.Sks, subject.DefaultSemester,
		subject.UpdatedBy, subject.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("subject not found")
	}

	return nil
}

// Delete subject
func (r *SubjectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subjects WHERE id = $1`

	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("subject not found")
	}

	return nil
}
