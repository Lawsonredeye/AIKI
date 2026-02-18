package service

import (
	"aiki/internal/domain"
	"aiki/internal/repository"
	"context"
)

//go:generate mockgen -source=job_service.go -destination=mocks/mock_job_service.go -package=mocks

type JobService interface {
	Create(ctx context.Context, job *domain.Job) (int32, error)
	Update(ctx context.Context, jobId int32, job *domain.Job) error
	Delete(ctx context.Context, jobId int32) error
	GetByID(ctx context.Context, jobId int32) (*domain.Job, error)
	GetAllByUserID(ctx context.Context, userId int32) ([]domain.Job, error)
}

type jobService struct {
	jobRepo repository.JobRepository
}

func NewJobService(jobRepo repository.JobRepository) JobService {
	return &jobService{
		jobRepo: jobRepo,
	}
}

func (s *jobService) Create(ctx context.Context, job *domain.Job) (int32, error) {
	jobId, err := s.jobRepo.Create(ctx, job)
	if err != nil {
		return 0, err
	}
	return jobId, nil
}

func (s *jobService) Update(ctx context.Context, jobId int32, job *domain.Job) error {
	// if job exists before updating
	_, err := s.jobRepo.GetJobByID(ctx, jobId)
	if err != nil {
		return err
	}

	err = s.jobRepo.Update(ctx, jobId, job)
	if err != nil {
		return err
	}
	return nil
}

func (s *jobService) Delete(ctx context.Context, jobId int32) error {
	// Verify job exists before deleting
	_, err := s.jobRepo.GetJobByID(ctx, jobId)
	if err != nil {
		return err
	}

	err = s.jobRepo.DeleteJob(ctx, jobId)
	if err != nil {
		return err
	}
	return nil
}

func (s *jobService) GetByID(ctx context.Context, jobId int32) (*domain.Job, error) {
	job, err := s.jobRepo.GetJobByID(ctx, jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *jobService) GetAllByUserID(ctx context.Context, userId int32) ([]domain.Job, error) {
	dbJobs, err := s.jobRepo.GetAllJobs(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Convert db.Job to domain.Job
	jobs := make([]domain.Job, 0, len(dbJobs))
	for _, dbJob := range dbJobs {
		job := domain.Job{
			ID:        dbJob.ID,
			UserId:    dbJob.UserID,
			Title:     dbJob.Title,
			Status:    dbJob.Status,
			CreatedAt: dbJob.CreatedAt.Time,
		}

		// Handle nullable fields
		if dbJob.CompanyName != nil {
			job.CompanyName = *dbJob.CompanyName
		}
		if dbJob.Notes != nil {
			job.Notes = *dbJob.Notes
		}
		if dbJob.Location != nil {
			job.Location = *dbJob.Location
		}
		if dbJob.Platform != nil {
			job.Platform = *dbJob.Platform
		}
		if dbJob.Link != nil {
			job.Link = *dbJob.Link
		}
		if dbJob.DateApplied.Valid {
			job.DateApplied = dbJob.DateApplied.Time.Format("2006-01-02")
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
