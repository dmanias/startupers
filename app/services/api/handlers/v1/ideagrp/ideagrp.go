package ideagrp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/aigrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/moderationgrp"
	"github.com/dmanias/startupers/business/core/idea"
	v1 "github.com/dmanias/startupers/business/web/v1"
	"github.com/dmanias/startupers/business/web/v1/paging"
	"github.com/dmanias/startupers/foundation/web"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DalleResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		B64JSON string `json:"b64_json"`
	} `json:"data"`
}

// Handlers manages the set of idea endpoints.
type Handlers struct {
	idea               *idea.Core
	log                *zap.SugaredLogger
	aiHandlers         *aigrp.Handlers
	moderationHandlers *moderationgrp.Handlers
	APIHost            string
}

// New constructs a handlers for route access.
func New(idea *idea.Core, log *zap.SugaredLogger, aiHandlers *aigrp.Handlers, moderationHandlers *moderationgrp.Handlers, APIHost string) *Handlers {
	return &Handlers{
		idea:               idea,
		log:                log,
		aiHandlers:         aiHandlers,
		moderationHandlers: moderationHandlers,
		APIHost:            APIHost,
	}
}

//func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//	var app AppNewIdea
//	if err := web.Decode(r, &app); err != nil {
//		return err
//	}
//
//	// Assume the question is the idea's description
//	ideaDescr := "Idea title:" + app.Title + ", " + "Idea description:" + app.Description + ", " + "Idea tags:" + fmt.Sprint(app.Tags)
//
//	// Get the instruction from the moderator
//	instruction, err := h.moderationHandlers.QueryByName(ctx, "avatar")
//	if err != nil {
//		return err
//	}
//
//	// Append the instruction to the question
//	question := instruction + ": " + ideaDescr
//	fmt.Printf("question.Query with moderator: %s\n", question)
//
//	// Call the DALLE function
//	aiResponse, err := h.aiHandlers.Dalle(ctx, question)
//	if err != nil {
//		if errors.Is(err, context.DeadlineExceeded) {
//			h.log.Errorf("DALL-E API request timed out")
//			return v1.NewRequestError(errors.New("DALL-E API request timed out"), http.StatusServiceUnavailable)
//		}
//		return err
//	}
//
//	// Download the image from the URL
//	h.log.Debugf("Downloading image from URL: %s", aiResponse)
//	resp, err := http.Get(aiResponse)
//	if err != nil {
//		h.log.Errorf("Error downloading image: %v", err)
//		return err
//	}
//	defer func() {
//		closeErr := resp.Body.Close()
//		if closeErr != nil {
//			h.log.Errorf("Error closing response body: %v", closeErr)
//		}
//	}()
//
//	// Read the image data
//	h.log.Debug("Reading image data")
//	imgData, err := io.ReadAll(resp.Body)
//	if err != nil {
//		h.log.Errorf("Error reading image data: %v", err)
//		return err
//	}
//
//	// Decode the image data
//	h.log.Debug("Decoding PNG image")
//	img, err := png.Decode(bytes.NewReader(imgData))
//	if err != nil {
//		h.log.Errorf("Error decoding PNG image: %v", err)
//		return err
//	}
//
//	// Compress the image as JPEG
//	var buf bytes.Buffer
//	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 70})
//	if err != nil {
//		h.log.Errorf("Error compressing image as JPEG: %v", err)
//		return err
//	}
//
//	compressedData := buf.Bytes()
//
//	// Convert the compressed image data to base64
//	base64Data := base64.StdEncoding.EncodeToString(compressedData)
//
//	// Set the AvatarURL field of AppNewIdea to the base64-encoded image data
//	app.AvatarURL = base64Data
//
//	nc, err := toCoreNewIdea(app)
//	if err != nil {
//		return v1.NewRequestError(err, http.StatusBadRequest)
//	}
//
//	newIdea, err := h.idea.Create(ctx, nc)
//	if err != nil {
//		return fmt.Errorf("create: idea[%+v]: %w", newIdea, err)
//	}
//
//	return web.Respond(ctx, w, toAppIdea(newIdea), http.StatusCreated)
//}

func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewIdea
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	// Assume the question is the idea's description
	ideaDescr := "Idea title:" + app.Title + ", " + "Idea description:" + app.Description + ", " + "Idea tags:" + fmt.Sprint(app.Tags)

	// Get the instruction from the moderator
	instruction, _, err := h.moderationHandlers.QueryByName(ctx, "avatar")
	if err != nil {
		return err
	}

	// Append the instruction to the question
	question := instruction + ": " + ideaDescr
	fmt.Printf("question.Query with moderator: %s\n", question)

	// Call the DALLE function
	aiResponse, err := h.aiHandlers.Dalle(ctx, question)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.log.Errorf("DALL-E API request timed out")
			return v1.NewRequestError(errors.New("DALL-E API request timed out"), http.StatusServiceUnavailable)
		}
		return err
	}

	// Download the image from the URL
	h.log.Debugf("Downloading image from URL: %s", aiResponse)
	resp, err := http.Get(aiResponse)
	if err != nil {
		h.log.Errorf("Error downloading image: %v", err)
		return err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			h.log.Errorf("Error closing response body: %v", closeErr)
		}
	}()

	// Read the image data
	h.log.Debug("Reading image data")
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Errorf("Error reading image data: %v", err)
		return err
	}

	// Generate a unique filename for the image
	filename := generateUniqueFilename()
	imageURL := ""
	// Save the image file locally
	imageURL, err = saveImageLocally(filename, imgData)
	if err != nil {
		h.log.Errorf("Error saving image locally: %v", err)
		return err
	}
	// Set the AvatarURL field of AppNewIdea to the image URL
	app.AvatarURL = imageURL

	nc, err := toCoreNewIdea(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}
	fmt.Printf("store idea2")
	newIdea, err := h.idea.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: idea[%+v]: %w", newIdea, err)
	}

	return web.Respond(ctx, w, toAppIdea(newIdea), http.StatusCreated)
}

func generateUniqueFilename() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("image_%d.png", timestamp)
}

func saveImageLocally(filename string, data []byte) (string, error) {
	// Specify the directory where you want to save the images
	saveDir := "uploads"

	// Create the directory if it doesn't exist
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}
	// Create the file path by joining the directory and filename
	filePath := filepath.Join(saveDir, filename)
	// Decode the image data
	//img, _, err := image.Decode(bytes.NewReader(data))
	//if err != nil {
	//	return fmt.Errorf("error decoding image: %v", err)
	//}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("error decoding PNG image: %v", err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(file)

	// Compress the image as JPEG with a quality of 75
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return "", fmt.Errorf("error compressing image as JPEG: %v", err)
	}

	return filePath, nil
}

// Update updates an existing idea in the system.
func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppUpdateIdea
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	uc, err := toCoreUpdateIdea(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	idea, err := h.idea.QueryByID(ctx, *uc.ID)
	if err != nil {
		return fmt.Errorf("query: ideaID[%s]: %w", uc.ID, err)
	}

	updatedIdea, err := h.idea.Update(ctx, idea, uc)
	if err != nil {
		return fmt.Errorf("update: idea[%+v]: %w", idea, err)
	}

	return web.Respond(ctx, w, toAppIdea(updatedIdea), http.StatusOK)
}

// Delete removes an idea from the system.
func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ideaID, err := uuid.Parse(web.Param(r, "idea_id"))
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	idea, err := h.idea.QueryByID(ctx, ideaID)
	if err != nil {
		return fmt.Errorf("query: ideaID[%s]: %w", ideaID, err)
	}

	err = h.idea.Delete(ctx, idea)
	if err != nil {
		// Log the actual error message
		h.log.Errorf("Error deleting idea: %v", err)

		return v1.NewRequestError(fmt.Errorf("Cannot delete idea with existing posts"), http.StatusConflict)
	}

	// Delete the corresponding image from the uploads directory
	err = os.Remove(idea.AvatarURL)
	if err != nil {
		h.log.Errorf("Error deleting image: %v", err)
		// You can decide whether you want to return an error here or just log it
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// QueryByID retrieves an idea by its ID.
func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Extract the idea_id parameter from the URL.
	ideaID := web.Param(r, "idea_id")

	// Parse the idea_id string into a uuid.UUID.
	id, err := uuid.Parse(ideaID)
	if err != nil {
		return v1.NewRequestError(fmt.Errorf("invalid idea ID: %w", err), http.StatusBadRequest)
	}

	// Retrieve the idea from the database.
	idea, err := h.idea.QueryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("query idea by ID: %w", err)
	}

	// Convert the idea to the AppIdea type.
	appIdea := toAppIdea(idea)

	// Respond with the AppIdea instance and a 200 status code.
	return web.Respond(ctx, w, appIdea, http.StatusOK)
}

// Query returns a list of ideas with paging.
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := paging.ParseRequest(r)
	if err != nil {
		return err
	}

	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	fmt.Printf("orderBy from idfeagrp.go: %+v\n", orderBy)

	ideas, err := h.idea.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppIdea, len(ideas))
	for i, idea := range ideas {
		items[i] = toAppIdea(idea)
	}

	total, err := h.idea.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

func (h *Handlers) QueryTags(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tags, err := h.idea.QueryTags(ctx)
	if err != nil {
		return fmt.Errorf("querying tags: %w", err)
	}

	return web.Respond(ctx, w, tags, http.StatusOK)
}
