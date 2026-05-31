package classes

import (
	"context"
	"database/sql"
	"fmt"

	classesPb "lms-class-service/pb/classes"
)

// ClassRepository struct
type ClassRepository struct {
	Db *sql.DB
}

// List classes with pagination
func (r *ClassRepository) List(ctx context.Context, limit, offset uint32, keyword, orderBy, sort string) ([]*classesPb.Class, uint32, error) {
	query := `
		SELECT id, university_id, university_name, COALESCE(faculty_id::text, ''), COALESCE(faculty_name, ''),
			programme_id, programme_name, code, name,
			COALESCE(updated_by::text, ''), updated_at, created_at
		FROM classes
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

	var classes []*classesPb.Class
	for rows.Next() {
		var class classesPb.Class
		err := rows.Scan(
			&class.Id, &class.UniversityId, &class.UniversityName,
			&class.FacultyId, &class.FacultyName, &class.ProgrammeId,
			&class.ProgrammeName, &class.Code, &class.Name,
			&class.UpdatedBy, &class.UpdatedAt, &class.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		classes = append(classes, &class)
	}

	countQuery := `SELECT COUNT(*) FROM classes WHERE name ILIKE $1 OR code ILIKE $1`
	var count uint32
	err = r.Db.QueryRowContext(ctx, countQuery, searchKeyword).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return classes, count, nil
}

// Get class by ID
func (r *ClassRepository) Get(ctx context.Context, id string) (*classesPb.Class, error) {
	query := `
		SELECT id, university_id, university_name, COALESCE(faculty_id::text, ''), COALESCE(faculty_name, ''),
			programme_id, programme_name, code, name,
			COALESCE(updated_by::text, ''), updated_at, created_at
		FROM classes WHERE id = $1
	`

	var class classesPb.Class
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&class.Id, &class.UniversityId, &class.UniversityName,
		&class.FacultyId, &class.FacultyName, &class.ProgrammeId,
		&class.ProgrammeName, &class.Code, &class.Name,
		&class.UpdatedBy, &class.UpdatedAt, &class.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &class, nil
}

// Create class
func (r *ClassRepository) Create(ctx context.Context, class *classesPb.Class) error {
	query := `
		INSERT INTO classes (university_id, university_name, faculty_id, faculty_name,
			programme_id, programme_name, code, name, updated_by)
		VALUES ($1, $2, NULLIF($3, '')::uuid, NULLIF($4, ''), $5, $6, $7, $8, NULLIF($9, '')::uuid)
		RETURNING id, updated_at, created_at
	`

	return r.Db.QueryRowContext(ctx, query,
		class.UniversityId, class.UniversityName, class.FacultyId,
		class.FacultyName, class.ProgrammeId, class.ProgrammeName,
		class.Code, class.Name, class.UpdatedBy,
	).Scan(&class.Id, &class.UpdatedAt, &class.CreatedAt)
}

// Update class
func (r *ClassRepository) Update(ctx context.Context, class *classesPb.Class) error {
	query := `
		UPDATE classes SET
			university_id = $1, university_name = $2, faculty_id = NULLIF($3, '')::uuid, faculty_name = NULLIF($4, ''),
			programme_id = $5, programme_name = $6, code = $7, name = $8,
			updated_by = NULLIF($9, '')::uuid, updated_at = NOW()
		WHERE id = $10
	`

	result, err := r.Db.ExecContext(ctx, query,
		class.UniversityId, class.UniversityName, class.FacultyId,
		class.FacultyName, class.ProgrammeId, class.ProgrammeName,
		class.Code, class.Name, class.UpdatedBy, class.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("class not found")
	}

	return nil
}

// Delete class
func (r *ClassRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM classes WHERE id = $1`

	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("class not found")
	}

	return nil
}

// SubjectClassRepository struct
type SubjectClassRepository struct {
	Db *sql.DB
}

// List subject classes
func (r *SubjectClassRepository) List(ctx context.Context, classID string, limit, offset uint32, keyword, orderBy, sort string) ([]*classesPb.SubjectClass, uint32, error) {
	query := `
		SELECT id, subject_id, class_id, period, COALESCE(teacher_id::text, ''), teacher_name, name,
			timetable_day, timetable_time::text, COALESCE(updated_by::text, ''), updated_at, created_at
		FROM subjects_classes
		WHERE ($1 = '' OR class_id = NULLIF($1, '')::uuid) AND name ILIKE $2
		ORDER BY ` + orderBy + ` ` + sort + `
		LIMIT $3 OFFSET $4
	`

	searchKeyword := "%" + keyword + "%"
	rows, err := r.Db.QueryContext(ctx, query, classID, searchKeyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subjectClasses []*classesPb.SubjectClass
	for rows.Next() {
		var sc classesPb.SubjectClass
		err := rows.Scan(
			&sc.Id, &sc.SubjectId, &sc.ClassId, &sc.Period,
			&sc.TeacherId, &sc.TeacherName, &sc.Name,
			&sc.TimetableDay, &sc.TimetableTime, &sc.UpdatedBy,
			&sc.UpdatedAt, &sc.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subjectClasses = append(subjectClasses, &sc)
	}

	countQuery := `SELECT COUNT(*) FROM subjects_classes WHERE ($1 = '' OR class_id = NULLIF($1, '')::uuid) AND name ILIKE $2`
	var count uint32
	err = r.Db.QueryRowContext(ctx, countQuery, classID, searchKeyword).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return subjectClasses, count, nil
}

// Get subject class by ID
func (r *SubjectClassRepository) Get(ctx context.Context, id string) (*classesPb.SubjectClass, error) {
	query := `
		SELECT id, subject_id, class_id, period, COALESCE(teacher_id::text, ''), teacher_name, name,
			timetable_day, timetable_time::text, COALESCE(updated_by::text, ''), updated_at, created_at
		FROM subjects_classes WHERE id = $1
	`

	var sc classesPb.SubjectClass
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&sc.Id, &sc.SubjectId, &sc.ClassId, &sc.Period,
		&sc.TeacherId, &sc.TeacherName, &sc.Name,
		&sc.TimetableDay, &sc.TimetableTime, &sc.UpdatedBy,
		&sc.UpdatedAt, &sc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

// Create subject class
func (r *SubjectClassRepository) Create(ctx context.Context, sc *classesPb.SubjectClass) error {
	query := `
		INSERT INTO subjects_classes (subject_id, class_id, period, teacher_id, teacher_name, name, timetable_day, timetable_time, updated_by)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid, $5, $6, $7, NULLIF($8, '')::time, NULLIF($9, '')::uuid)
		RETURNING id, updated_at, created_at
	`

	return r.Db.QueryRowContext(ctx, query,
		sc.SubjectId, sc.ClassId, sc.Period, sc.TeacherId,
		sc.TeacherName, sc.Name, sc.TimetableDay, sc.TimetableTime,
		sc.UpdatedBy,
	).Scan(&sc.Id, &sc.UpdatedAt, &sc.CreatedAt)
}

// Update subject class
func (r *SubjectClassRepository) Update(ctx context.Context, sc *classesPb.SubjectClass) error {
	query := `
		UPDATE subjects_classes SET
			subject_id = $1, class_id = $2, period = $3, teacher_id = NULLIF($4, '')::uuid,
			teacher_name = $5, name = $6, timetable_day = $7, timetable_time = NULLIF($8, '')::time,
			updated_by = NULLIF($9, '')::uuid, updated_at = NOW()
		WHERE id = $10
	`

	result, err := r.Db.ExecContext(ctx, query,
		sc.SubjectId, sc.ClassId, sc.Period, sc.TeacherId,
		sc.TeacherName, sc.Name, sc.TimetableDay, sc.TimetableTime,
		sc.UpdatedBy, sc.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("subject class not found")
	}

	return nil
}

// Delete subject class
func (r *SubjectClassRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subjects_classes WHERE id = $1`

	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("subject class not found")
	}

	return nil
}

// StudentClassRepository struct
type StudentClassRepository struct {
	Db *sql.DB
}

// List student classes
func (r *StudentClassRepository) List(ctx context.Context, subjectClassID string, limit, offset uint32, keyword, orderBy, sort string) ([]*classesPb.StudentClass, uint32, error) {
	query := `
		SELECT id, subject_class_id, student_id, student_name,
			COALESCE(updated_by::text, ''), updated_at, created_at
		FROM student_classes
		WHERE subject_class_id = $1 AND student_name ILIKE $2
		ORDER BY ` + orderBy + ` ` + sort + `
		LIMIT $3 OFFSET $4
	`

	searchKeyword := "%" + keyword + "%"
	rows, err := r.Db.QueryContext(ctx, query, subjectClassID, searchKeyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var studentClasses []*classesPb.StudentClass
	for rows.Next() {
		var sc classesPb.StudentClass
		err := rows.Scan(
			&sc.Id, &sc.SubjectClassId, &sc.StudentId, &sc.StudentName,
			&sc.UpdatedBy, &sc.UpdatedAt, &sc.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		studentClasses = append(studentClasses, &sc)
	}

	countQuery := `SELECT COUNT(*) FROM student_classes WHERE subject_class_id = $1 AND student_name ILIKE $2`
	var count uint32
	err = r.Db.QueryRowContext(ctx, countQuery, subjectClassID, searchKeyword).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return studentClasses, count, nil
}

// Get student class by ID
func (r *StudentClassRepository) Get(ctx context.Context, id string) (*classesPb.StudentClass, error) {
	query := `
		SELECT id, subject_class_id, student_id, student_name,
			COALESCE(updated_by::text, ''), updated_at, created_at
		FROM student_classes WHERE id = $1
	`

	var sc classesPb.StudentClass
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&sc.Id, &sc.SubjectClassId, &sc.StudentId, &sc.StudentName,
		&sc.UpdatedBy, &sc.UpdatedAt, &sc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

// Create student class
func (r *StudentClassRepository) Create(ctx context.Context, sc *classesPb.StudentClass) error {
	query := `
		INSERT INTO student_classes (subject_class_id, student_id, student_name, updated_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, updated_at, created_at
	`

	return r.Db.QueryRowContext(ctx, query,
		sc.SubjectClassId, sc.StudentId, sc.StudentName, sc.UpdatedBy,
	).Scan(&sc.Id, &sc.UpdatedAt, &sc.CreatedAt)
}

// Delete student class
func (r *StudentClassRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM student_classes WHERE id = $1`

	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("student class not found")
	}

	return nil
}
