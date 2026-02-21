package candidates

import (
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	domainerrors "resume/internal/domain/errors"
	"resume/internal/middleware"
	"resume/internal/server/response"

	"github.com/gin-gonic/gin"
)

// ResumeUploader uploads file content and returns a storage URL (e.g. Supabase Storage).
type ResumeUploader interface {
	Upload(ctx context.Context, fileContent []byte, fileName string) (storageURL string, err error)
}

// Handler exposes candidate and resume HTTP endpoints.
type Handler struct {
	svc      *Service
	uploader ResumeUploader
}

// NewHandler creates a candidates handler.
func NewHandler(svc *Service, uploader ResumeUploader) *Handler {
	return &Handler{svc: svc, uploader: uploader}
}

// GetProfile returns the candidate profile for the authenticated user.
func (h *Handler) GetProfile(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	candidateID := c.Param("candidate_id")
	candidate, err := h.svc.GetCandidateByUserID(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Error(c, err)
		return
	}
	if candidate.ID != candidateID {
		response.Error(c, &domainerrors.ForbiddenError{Message: "forbidden"})
		return
	}
	response.JSON(c, http.StatusOK, candidate)
}

// UploadResumeRequest multipart or JSON with file upload; for simplicity we accept JSON with path after upload.
type UploadResumeRequest struct {
	FileName   string `json:"file_name" binding:"required"`
	StoragePath string `json:"storage_path" binding:"required"`
	FileSize   int64  `json:"file_size"`
	MimeType   string `json:"mime_type"`
}

// UploadResume stores resume metadata (file assumed uploaded to storagePath by caller).
func (h *Handler) UploadResume(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	candidateID := c.Param("candidate_id")
	candidate, err := h.svc.GetCandidateByUserID(c.Request.Context(), claims.UserID)
	if err != nil || candidate.ID != candidateID {
		response.Error(c, &domainerrors.ForbiddenError{Message: "forbidden"})
		return
	}
	var req UploadResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	resume, err := h.svc.AddResume(c.Request.Context(), candidateID, req.FileName, req.StoragePath, req.MimeType, req.FileSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, resume)
}

const maxResumeSize = 100 << 20 // 100 MB

// UploadResumePDF accepts a multipart form file (PDF only), uploads to Supabase Storage, and creates resume record.
func (h *Handler) UploadResumePDF(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	candidateID := c.Param("candidate_id")
	candidate, err := h.svc.GetCandidateByUserID(c.Request.Context(), claims.UserID)
	if err != nil || candidate.ID != candidateID {
		response.Error(c, &domainerrors.ForbiddenError{Message: "forbidden"})
		return
	}
	if h.uploader == nil {
		response.Error(c, &domainerrors.ValidationError{Field: "upload", Message: "resume upload not configured (Supabase Storage)"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ValidationError(c, err)
		return
	}
	defer file.Close()

	fileName := header.Filename
	if fileName == "" {
		fileName = "resume.pdf"
	}
	if !strings.EqualFold(path.Ext(fileName), ".pdf") {
		response.Error(c, &domainerrors.ValidationError{Field: "file", Message: "only PDF format is allowed"})
		return
	}
	contentType := header.Header.Get("Content-Type")
	if contentType != "" && !strings.EqualFold(contentType, "application/pdf") {
		response.Error(c, &domainerrors.ValidationError{Field: "file", Message: "only PDF format is allowed"})
		return
	}

	fileContent, err := io.ReadAll(io.LimitReader(file, maxResumeSize+1))
	if err != nil {
		response.Error(c, err)
		return
	}
	if len(fileContent) > maxResumeSize {
		response.Error(c, &domainerrors.ValidationError{Field: "file", Message: "file too large (max 100 MB)"})
		return
	}

	storageURL, err := h.uploader.Upload(c.Request.Context(), fileContent, fileName)
	if err != nil {
		response.Error(c, err)
		return
	}

	resume, err := h.svc.AddResume(c.Request.Context(), candidateID, fileName, storageURL, "application/pdf", int64(len(fileContent)))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, resume)
}

// ListResumes returns resumes for the candidate.
func (h *Handler) ListResumes(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	candidateID := c.Param("candidate_id")
	candidate, err := h.svc.GetCandidateByUserID(c.Request.Context(), claims.UserID)
	if err != nil || candidate.ID != candidateID {
		response.Error(c, &domainerrors.ForbiddenError{Message: "forbidden"})
		return
	}
	list, err := h.svc.ListResumes(c.Request.Context(), candidateID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, list)
}

// GetResume returns a resume by ID (must belong to candidate).
func (h *Handler) GetResume(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil {
		return
	}
	candidateID := c.Param("candidate_id")
	resumeID := c.Param("resume_id")
	candidate, err := h.svc.GetCandidateByUserID(c.Request.Context(), claims.UserID)
	if err != nil || candidate.ID != candidateID {
		response.Error(c, &domainerrors.ForbiddenError{Message: "forbidden"})
		return
	}
	resume, err := h.svc.GetResumeByID(c.Request.Context(), resumeID, candidateID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, resume)
}
