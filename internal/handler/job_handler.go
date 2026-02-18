package handler

import (
	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"aiki/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type JobHandler struct {
	jobService service.JobService
	validator  echo.Validator
}

func NewJobHandler(jobService service.JobService, validator echo.Validator) *JobHandler {
	return &JobHandler{
		jobService: jobService,
		validator:  validator,
	}
}

// CreateJob godoc
// @Summary Create a new job
// @Description Create a new job application for the authenticated user
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.JobRequest true "Job details"
// @Success 201 {object} response.Response{data=map[string]int32}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /jobs [post]
func (h *JobHandler) CreateJob(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	var req domain.JobRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	job := req.ToDomain(userID)

	jobID, err := h.jobService.Create(c.Request().Context(), &job)
	if err != nil {
		c.Logger().Errorf("failed to create job: %v", err)
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, "job created successfully", map[string]int32{"id": jobID})
}

// GetJob godoc
// @Summary Get a job by ID
// @Description Get a specific job application by ID
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response{data=domain.Job}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /jobs/{id} [get]
func (h *JobHandler) GetJob(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	jobID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return response.ValidationError(c, "invalid job ID")
	}

	job, err := h.jobService.GetByID(c.Request().Context(), int32(jobID))
	if err != nil {
		return response.Error(c, err)
	}

	// Verify the job belongs to the authenticated user
	if job.UserId != userID {
		return response.Error(c, domain.ErrUnauthorized)
	}

	return response.Success(c, http.StatusOK, "job retrieved successfully", job)
}

// GetAllJobs godoc
// @Summary Get all jobs
// @Description Get all job applications for the authenticated user
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]domain.Job}
// @Failure 401 {object} response.Response
// @Router /jobs [get]
func (h *JobHandler) GetAllJobs(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	jobs, err := h.jobService.GetAllByUserID(c.Request().Context(), userID)
	if err != nil {
		c.Logger().Errorf("failed to get jobs: %v", err)
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "jobs retrieved successfully", jobs)
}

// UpdateJob godoc
// @Summary Update a job
// @Description Update an existing job application
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Param request body domain.JobRequest true "Updated job details"
// @Success 200 {object} response.Response{data=domain.Job}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /jobs/{id} [put]
func (h *JobHandler) UpdateJob(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	jobID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return response.ValidationError(c, "invalid job ID")
	}

	// Verify ownership before update
	existingJob, err := h.jobService.GetByID(c.Request().Context(), int32(jobID))
	if err != nil {
		return response.Error(c, err)
	}

	if existingJob.UserId != userID {
		return response.Error(c, domain.ErrUnauthorized)
	}

	var req domain.JobRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	job := req.ToDomain(userID)

	err = h.jobService.Update(c.Request().Context(), int32(jobID), &job)
	if err != nil {
		c.Logger().Errorf("failed to update job: %v", err)
		return response.Error(c, err)
	}

	// Get updated job to return
	updatedJob, err := h.jobService.GetByID(c.Request().Context(), int32(jobID))
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "job updated successfully", updatedJob)
}

// DeleteJob godoc
// @Summary Delete a job
// @Description Delete a job application
// @Tags jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /jobs/{id} [delete]
func (h *JobHandler) DeleteJob(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	jobID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return response.ValidationError(c, "invalid job ID")
	}

	// Verify ownership before delete
	job, err := h.jobService.GetByID(c.Request().Context(), int32(jobID))
	if err != nil {
		return response.Error(c, err)
	}

	if job.UserId != userID {
		return response.Error(c, domain.ErrUnauthorized)
	}

	err = h.jobService.Delete(c.Request().Context(), int32(jobID))
	if err != nil {
		c.Logger().Errorf("failed to delete job: %v", err)
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "job deleted successfully", nil)
}
