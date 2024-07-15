package scheme

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Create uuid extension",
		Script:      `CREATE EXTENSION "uuid-ossp";`,
	},
	{
		Version:     2,
		Description: "Create subjects table ",
		Script:`
			CREATE TABLE subjects (
				id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
				university_id uuid NOT NULL,
				university_name varchar(45) NOT NULL,
				faculty_id uuid,
				faculty_name varchar(45),
				programme_id uuid NOT NULL,
				programme_name varchar(45) NOT NULL,
				code char(6) NOT NULL,
				name varchar(45) NOT NULL,
				sks smallint NOT NULL,
				default_semester smallint NOT NULL,
				updated_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				updated_by uuid,
				created_at timestamptz NOT NULL DEFAULT timezone('utc', NOW())
			);
		`,
	},
	{
		Version:     3,
		Description: "Create classes table",
		Script:`
			CREATE TABLE classes (
				id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
				university_id uuid NOT NULL,
				university_name varchar(45) NOT NULL,
				faculty_id uuid,
				faculty_name varchar(45),
				programme_id uuid NOT NULL,
				programme_name varchar(45) NOT NULL,
				code char(6) NOT NULL,
				name varchar(45) NOT NULL,
				updated_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				updated_by uuid,
				created_at timestamptz NOT NULL DEFAULT timezone('utc', NOW())
			);

		`,
	},
	{
		Version:     4,
		Description: "Create subjects_classes table",
		Script:`
			CREATE TABLE subjects_classes (
				id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
				subject_id uuid NOT NULL,
				class_id uuid NOT NULL,
				period char(16) NOT NULL,
				teacher_id uuid NOT NULL,
				teacher_name varchar(45) NOT NULL,
				name varchar(45) NOT NULL,
				timetable_day smallint NOT NULL DEFAULT EXTRACT(ISODOW FROM timezone('utc', NOW())), -- Default day of the week (1-7)
				timetable_time time NOT NULL DEFAULT timezone('utc', NOW())::time, -- Default time (current time)
				updated_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				updated_by uuid,
				created_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				FOREIGN KEY (subject_id) REFERENCES subjects (id) ON UPDATE CASCADE ON DELETE CASCADE,
				FOREIGN KEY (class_id) REFERENCES classes (id) ON UPDATE CASCADE ON DELETE CASCADE
			);
		`,
	},
	{
		Version:     5,
		Description: "Create student_classes table",
		Script:`
			CREATE TABLE student_classes (
				id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
				subject_class_id uuid NOT NULL,
				student_id uuid NOT NULL,
				student_name varchar(45) NOT NULL,
				updated_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				updated_by uuid,
				created_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				FOREIGN KEY (subject_class_id) REFERENCES subjects_classes (id) ON UPDATE CASCADE ON DELETE CASCADE
			);
		`,
	},
	{
		Version:     6,
		Description: "Create topic_subject table",
		Script:`
			CREATE TABLE topic_subjects (
				id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
				subject_id uuid NOT NULL,
				name varchar(45) NOT NULL,
				updated_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				updated_by uuid,
				created_at timestamptz NOT NULL DEFAULT timezone('utc', NOW()),
				FOREIGN KEY (subject_id) REFERENCES subjects (id) ON UPDATE CASCADE ON DELETE CASCADE
			);
		`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}