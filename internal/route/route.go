package route

import (
	"database/sql"
	"log"

	"google.golang.org/grpc"

	classesDomain "lms-class-service/internal/domain/classes"
	subjectsDomain "lms-class-service/internal/domain/subjects"
	"lms-class-service/internal/pkg/db/redis"
	classesPb "lms-class-service/pb/classes"
	subjectsPb "lms-class-service/pb/subjects"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log *log.Logger, cache *redis.Cache) {
	// Subject service
	subjectServer := subjectsDomain.SubjectService{Db: db, Cache: cache, Log: log}
	subjectsPb.RegisterSubjectServiceServer(grpcServer, &subjectServer)

	// Class service
	classServer := classesDomain.ClassServiceServer{Db: db, Cache: cache, Log: log}
	classesPb.RegisterClassServiceServer(grpcServer, &classServer)

	// Subject class service
	subjectClassServer := classesDomain.SubjectClassServiceServer{Db: db, Cache: cache, Log: log}
	classesPb.RegisterSubjectClassServiceServer(grpcServer, &subjectClassServer)

	// Student class service
	studentClassServer := classesDomain.StudentClassServiceServer{Db: db, Cache: cache, Log: log}
	classesPb.RegisterStudentClassServiceServer(grpcServer, &studentClassServer)
}