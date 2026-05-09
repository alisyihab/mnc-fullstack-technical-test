package handlers

import (
	"mnc-fullstack-technical-test/tahap-2/internal/delivery/http/response"
	"mnc-fullstack-technical-test/tahap-2/internal/infrastructure/worker"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler exposes monitoring endpoints for the transfer worker.
type DashboardHandler struct {
	transferWorker *worker.TransferWorker
}

func NewDashboardHandler(tw *worker.TransferWorker) *DashboardHandler {
	return &DashboardHandler{transferWorker: tw}
}

// GetQueueStats godoc
// @Summary Get queue statistics
// @Description Returns current transfer worker queue statistics
// @Tags dashboard
// @Produce json
// @Success 200 {object} response.JSONResponse
// @Router /dashboard/queue-stats [get]
func (h *DashboardHandler) GetQueueStats(c *gin.Context) {
	response.Success(c, http.StatusOK, h.transferWorker.Stats())
}

// GetJobs godoc
// @Summary Get recent transfer jobs
// @Description Returns the last 100 processed transfer jobs with their status
// @Tags dashboard
// @Produce json
// @Success 200 {object} response.JSONResponse
// @Router /dashboard/jobs [get]
func (h *DashboardHandler) GetJobs(c *gin.Context) {
	jobs := h.transferWorker.RecentJobs()
	if jobs == nil {
		jobs = []worker.JobRecord{}
	}
	response.Success(c, http.StatusOK, jobs)
}
