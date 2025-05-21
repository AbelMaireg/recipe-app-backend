package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"app/services"
	"app/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadRecipePictureHandler struct {
	recipeService services.RecipeService
	s3Client      *s3.Client
	bucketName    string
}

func NewUploadRecipePictureHandler(recipeService services.RecipeService, s3Client *s3.Client, bucketName string) *UploadRecipePictureHandler {
	return &UploadRecipePictureHandler{
		recipeService: recipeService,
		s3Client:      s3Client,
		bucketName:    bucketName,
	}
}

type UploadRecipePictureResponse struct {
	ID        string    `json:"id"`
	RecipeID  string    `json:"recipe_id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *UploadRecipePictureHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_FORM", "Failed to parse form: "+err.Error())
		return
	}

	_, err := utils.ParseJWT(r.Header.Get("Authorization"))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization token")
		return
	}

	recipeIDStr := r.FormValue("recipe_id")
	if recipeIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_RECIPE_ID", "recipe_id is required")
		return
	}

	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_RECIPE_ID", "Invalid recipe_id: "+err.Error())
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_FILE", "Failed to get file: "+err.Error())
		return
	}
	defer file.Close()

	if handler.Size > 5<<20 { // 5MB limit
		utils.WriteError(w, http.StatusBadRequest, "FILE_TOO_LARGE", "File size exceeds 5MB")
		return
	}
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_FILE_TYPE", "Only JPG and PNG files are allowed")
		return
	}

	fileID := uuid.New()
	objectKey := fmt.Sprintf("recipe/%s/%s%s", recipeID, fileID, ext)

	_, err = h.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      &h.bucketName,
		Key:         &objectKey,
		Body:        file,
		ContentType: &handler.Header["Content-Type"][0],
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_FAILED", "Failed to upload to MinIO: "+err.Error())
		return
	}

	picture, err := h.recipeService.SaveRecipePicture(recipeID, objectKey)
	if err != nil {
		h.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: &h.bucketName,
			Key:    &objectKey,
		})
		if strings.Contains(err.Error(), "foreign key") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_RECIPE", "Recipe does not exist or not owned by user")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "DB_ERROR", "Failed to save picture: "+err.Error())
		}
		return
	}

	response := UploadRecipePictureResponse{
		ID:        picture.ID.String(),
		RecipeID:  picture.RecipeId.String(),
		Path:      picture.Path,
		CreatedAt: picture.CreatedAt,
	}
	utils.EncodeJSON(w, response)
}
