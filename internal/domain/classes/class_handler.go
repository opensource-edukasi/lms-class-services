package classes

import (
	"context"
	"database/sql"
	"log"
	"time"

	"lms-class-service/internal/pkg/app"
	"lms-class-service/internal/pkg/db/redis"
	classesPb "lms-class-service/pb/classes"
	genericPb "lms-class-service/pb/generic"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ClassServiceServer struct
type ClassServiceServer struct {
	Db    *sql.DB
	Cache *redis.Cache
	Log   *log.Logger
	classesPb.UnimplementedClassServiceServer
}

// List classes
func (s *ClassServiceServer) List(ctx context.Context, in *classesPb.ClassListInput) (*classesPb.ClassList, error) {
	repo := ClassRepository{Db: s.Db}

	pagination := in.GetPagination()
	limit := uint32(10)
	offset := uint32(0)
	keyword := ""
	orderBy := "created_at"
	sort := "DESC"

	if pagination != nil {
		if pagination.Limit > 0 {
			limit = pagination.Limit
		}
		offset = pagination.Offset
		if pagination.Keyword != "" {
			keyword = pagination.Keyword
		}
		if pagination.Order != "" {
			orderBy = pagination.Order
		}
		if pagination.Sort != "" {
			sort = pagination.Sort
		}
	}

	classes, count, err := repo.List(ctx, limit, offset, keyword, orderBy, sort)
	if err != nil {
		s.Log.Printf("error listing classes: %v", err)
		return nil, status.Error(codes.Internal, "failed to list classes")
	}

	return &classesPb.ClassList{
		Classes: classes,
		Count:   count,
	}, nil
}

// Get class by ID
func (s *ClassServiceServer) Get(ctx context.Context, in *genericPb.Id) (*classesPb.Class, error) {
	repo := ClassRepository{Db: s.Db}

	class, err := repo.Get(ctx, in.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "class not found")
		}
		s.Log.Printf("error getting class: %v", err)
		return nil, status.Error(codes.Internal, "failed to get class")
	}

	return class, nil
}

// Create class
func (s *ClassServiceServer) Create(ctx context.Context, in *classesPb.ClassInput) (*classesPb.Class, error) {
	repo := ClassRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)

	class := &classesPb.Class{
		UniversityId:   in.GetUniversityId(),
		UniversityName: in.GetUniversityName(),
		FacultyId:      in.GetFacultyId(),
		FacultyName:    in.GetFacultyName(),
		ProgrammeId:    in.GetProgrammeId(),
		ProgrammeName:  in.GetProgrammeName(),
		Code:           in.GetCode(),
		Name:           in.GetName(),
		UpdatedBy:      userID,
	}

	err := repo.Create(ctx, class)
	if err != nil {
		s.Log.Printf("error creating class: %v", err)
		return nil, status.Error(codes.Internal, "failed to create class")
	}

	return class, nil
}

// Update class
func (s *ClassServiceServer) Update(ctx context.Context, in *classesPb.Class) (*classesPb.Class, error) {
	repo := ClassRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)
	in.UpdatedBy = userID
	in.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	err := repo.Update(ctx, in)
	if err != nil {
		s.Log.Printf("error updating class: %v", err)
		return nil, status.Error(codes.Internal, "failed to update class")
	}

	return in, nil
}

// Delete class
func (s *ClassServiceServer) Delete(ctx context.Context, in *genericPb.Id) (*genericPb.BoolMessage, error) {
	repo := ClassRepository{Db: s.Db}

	err := repo.Delete(ctx, in.GetId())
	if err != nil {
		s.Log.Printf("error deleting class: %v", err)
		return &genericPb.BoolMessage{IsTrue: false}, status.Error(codes.Internal, "failed to delete class")
	}

	return &genericPb.BoolMessage{IsTrue: true}, nil
}

// SubjectClassServiceServer struct
type SubjectClassServiceServer struct {
	Db    *sql.DB
	Cache *redis.Cache
	Log   *log.Logger
	classesPb.UnimplementedSubjectClassServiceServer
}

// List subject classes
func (s *SubjectClassServiceServer) List(ctx context.Context, in *classesPb.SubjectClassListInput) (*classesPb.SubjectClassList, error) {
	repo := SubjectClassRepository{Db: s.Db}

	pagination := in.GetPagination()
	limit := uint32(10)
	offset := uint32(0)
	keyword := ""
	orderBy := "created_at"
	sort := "DESC"

	if pagination != nil {
		if pagination.Limit > 0 {
			limit = pagination.Limit
		}
		offset = pagination.Offset
		if pagination.Keyword != "" {
			keyword = pagination.Keyword
		}
		if pagination.Order != "" {
			orderBy = pagination.Order
		}
		if pagination.Sort != "" {
			sort = pagination.Sort
		}
	}

	subjectClasses, count, err := repo.List(ctx, in.GetClassId(), limit, offset, keyword, orderBy, sort)
	if err != nil {
		s.Log.Printf("error listing subject classes: %v", err)
		return nil, status.Error(codes.Internal, "failed to list subject classes")
	}

	return &classesPb.SubjectClassList{
		SubjectClasses: subjectClasses,
		Count:          count,
	}, nil
}

// Get subject class by ID
func (s *SubjectClassServiceServer) Get(ctx context.Context, in *genericPb.Id) (*classesPb.SubjectClass, error) {
	repo := SubjectClassRepository{Db: s.Db}

	sc, err := repo.Get(ctx, in.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "subject class not found")
		}
		s.Log.Printf("error getting subject class: %v", err)
		return nil, status.Error(codes.Internal, "failed to get subject class")
	}

	return sc, nil
}

// Create subject class
func (s *SubjectClassServiceServer) Create(ctx context.Context, in *classesPb.SubjectClassInput) (*classesPb.SubjectClass, error) {
	repo := SubjectClassRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)

	sc := &classesPb.SubjectClass{
		SubjectId:     in.GetSubjectId(),
		ClassId:       in.GetClassId(),
		Period:        in.GetPeriod(),
		TeacherId:     in.GetTeacherId(),
		TeacherName:   in.GetTeacherName(),
		Name:          in.GetName(),
		TimetableDay:  in.GetTimetableDay(),
		TimetableTime: in.GetTimetableTime(),
		UpdatedBy:     userID,
	}

	err := repo.Create(ctx, sc)
	if err != nil {
		s.Log.Printf("error creating subject class: %v", err)
		return nil, status.Error(codes.Internal, "failed to create subject class")
	}

	return sc, nil
}

// Update subject class
func (s *SubjectClassServiceServer) Update(ctx context.Context, in *classesPb.SubjectClass) (*classesPb.SubjectClass, error) {
	repo := SubjectClassRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)
	in.UpdatedBy = userID
	in.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	err := repo.Update(ctx, in)
	if err != nil {
		s.Log.Printf("error updating subject class: %v", err)
		return nil, status.Error(codes.Internal, "failed to update subject class")
	}

	return in, nil
}

// Delete subject class
func (s *SubjectClassServiceServer) Delete(ctx context.Context, in *genericPb.Id) (*genericPb.BoolMessage, error) {
	repo := SubjectClassRepository{Db: s.Db}

	err := repo.Delete(ctx, in.GetId())
	if err != nil {
		s.Log.Printf("error deleting subject class: %v", err)
		return &genericPb.BoolMessage{IsTrue: false}, status.Error(codes.Internal, "failed to delete subject class")
	}

	return &genericPb.BoolMessage{IsTrue: true}, nil
}

// StudentClassServiceServer struct
type StudentClassServiceServer struct {
	Db    *sql.DB
	Cache *redis.Cache
	Log   *log.Logger
	classesPb.UnimplementedStudentClassServiceServer
}

// List student classes
func (s *StudentClassServiceServer) List(ctx context.Context, in *classesPb.StudentClassListInput) (*classesPb.StudentClassList, error) {
	repo := StudentClassRepository{Db: s.Db}

	pagination := in.GetPagination()
	limit := uint32(10)
	offset := uint32(0)
	keyword := ""
	orderBy := "created_at"
	sort := "DESC"

	if pagination != nil {
		if pagination.Limit > 0 {
			limit = pagination.Limit
		}
		offset = pagination.Offset
		if pagination.Keyword != "" {
			keyword = pagination.Keyword
		}
		if pagination.Order != "" {
			orderBy = pagination.Order
		}
		if pagination.Sort != "" {
			sort = pagination.Sort
		}
	}

	studentClasses, count, err := repo.List(ctx, in.GetSubjectClassId(), limit, offset, keyword, orderBy, sort)
	if err != nil {
		s.Log.Printf("error listing student classes: %v", err)
		return nil, status.Error(codes.Internal, "failed to list student classes")
	}

	return &classesPb.StudentClassList{
		StudentClasses: studentClasses,
		Count:          count,
	}, nil
}

// Get student class by ID
func (s *StudentClassServiceServer) Get(ctx context.Context, in *genericPb.Id) (*classesPb.StudentClass, error) {
	repo := StudentClassRepository{Db: s.Db}

	sc, err := repo.Get(ctx, in.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "student class not found")
		}
		s.Log.Printf("error getting student class: %v", err)
		return nil, status.Error(codes.Internal, "failed to get student class")
	}

	return sc, nil
}

// Create student class
func (s *StudentClassServiceServer) Create(ctx context.Context, in *classesPb.StudentClassInput) (*classesPb.StudentClass, error) {
	repo := StudentClassRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)

	sc := &classesPb.StudentClass{
		SubjectClassId: in.GetSubjectClassId(),
		StudentId:      in.GetStudentId(),
		StudentName:    in.GetStudentName(),
		UpdatedBy:      userID,
	}

	err := repo.Create(ctx, sc)
	if err != nil {
		s.Log.Printf("error creating student class: %v", err)
		return nil, status.Error(codes.Internal, "failed to create student class")
	}

	return sc, nil
}

// Delete student class
func (s *StudentClassServiceServer) Delete(ctx context.Context, in *genericPb.Id) (*genericPb.BoolMessage, error) {
	repo := StudentClassRepository{Db: s.Db}

	err := repo.Delete(ctx, in.GetId())
	if err != nil {
		s.Log.Printf("error deleting student class: %v", err)
		return &genericPb.BoolMessage{IsTrue: false}, status.Error(codes.Internal, "failed to delete student class")
	}

	return &genericPb.BoolMessage{IsTrue: true}, nil
}
