package subjects

import (
	"context"
	"database/sql"
	"log"
	"time"

	"lms-class-service/internal/pkg/app"
	"lms-class-service/internal/pkg/db/redis"
	genericPb "lms-class-service/pb/generic"
	subjectsPb "lms-class-service/pb/subjects"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SubjectService struct
type SubjectService struct {
	Db    *sql.DB
	Cache *redis.Cache
	Log   *log.Logger
	subjectsPb.UnimplementedSubjectServiceServer
}

// List subjects
func (s *SubjectService) List(ctx context.Context, in *subjectsPb.SubjectListInput) (*subjectsPb.SubjectList, error) {
	repo := SubjectRepository{Db: s.Db}

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

	subjects, count, err := repo.List(ctx, limit, offset, keyword, orderBy, sort)
	if err != nil {
		s.Log.Printf("error listing subjects: %v", err)
		return nil, status.Error(codes.Internal, "failed to list subjects")
	}

	return &subjectsPb.SubjectList{
		Subjects: subjects,
		Count:    count,
	}, nil
}

// Get subject by ID
func (s *SubjectService) Get(ctx context.Context, in *genericPb.Id) (*subjectsPb.Subject, error) {
	repo := SubjectRepository{Db: s.Db}

	subject, err := repo.Get(ctx, in.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "subject not found")
		}
		s.Log.Printf("error getting subject: %v", err)
		return nil, status.Error(codes.Internal, "failed to get subject")
	}

	return subject, nil
}

// Create subject
func (s *SubjectService) Create(ctx context.Context, in *subjectsPb.SubjectInput) (*subjectsPb.Subject, error) {
	repo := SubjectRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)

	subject := &subjectsPb.Subject{
		UniversityId:    in.GetUniversityId(),
		UniversityName:  in.GetUniversityName(),
		FacultyId:       in.GetFacultyId(),
		FacultyName:     in.GetFacultyName(),
		ProgrammeId:     in.GetProgrammeId(),
		ProgrammeName:   in.GetProgrammeName(),
		Code:            in.GetCode(),
		Name:            in.GetName(),
		Sks:             in.GetSks(),
		DefaultSemester: in.GetDefaultSemester(),
		UpdatedBy:       userID,
	}

	err := repo.Create(ctx, subject)
	if err != nil {
		s.Log.Printf("error creating subject: %v", err)
		return nil, status.Error(codes.Internal, "failed to create subject")
	}

	// Create topic subjects
	for _, topic := range in.GetTopics() {
		topicSubject := &subjectsPb.TopicSubject{
			SubjectId: subject.Id,
			Name:      topic.GetName(),
			UpdatedBy: userID,
		}
		err := repo.CreateTopic(ctx, topicSubject)
		if err != nil {
			s.Log.Printf("error creating topic subject: %v", err)
			return nil, status.Error(codes.Internal, "failed to create topic subject")
		}
		subject.Topics = append(subject.Topics, topicSubject)
	}

	return subject, nil
}

// Update subject
func (s *SubjectService) Update(ctx context.Context, in *subjectsPb.Subject) (*subjectsPb.Subject, error) {
	repo := SubjectRepository{Db: s.Db}

	userID := ctx.Value(app.Ctx("user_id")).(string)
	in.UpdatedBy = userID
	in.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	err := repo.Update(ctx, in)
	if err != nil {
		s.Log.Printf("error updating subject: %v", err)
		return nil, status.Error(codes.Internal, "failed to update subject")
	}

	return in, nil
}

// Delete subject
func (s *SubjectService) Delete(ctx context.Context, in *genericPb.Id) (*genericPb.BoolMessage, error) {
	repo := SubjectRepository{Db: s.Db}

	err := repo.Delete(ctx, in.GetId())
	if err != nil {
		s.Log.Printf("error deleting subject: %v", err)
		return &genericPb.BoolMessage{IsTrue: false}, status.Error(codes.Internal, "failed to delete subject")
	}

	return &genericPb.BoolMessage{IsTrue: true}, nil
}
