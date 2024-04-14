package handlers

import (
	"encoding/json"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/aigrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/moderationgrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/testgrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/usergrp"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/core/user"
	"github.com/dmanias/startupers/business/core/user/stores/userdb"
	"github.com/dmanias/startupers/business/web/auth"
	"github.com/dmanias/startupers/business/web/v1/mid"
	"github.com/dmanias/startupers/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	DB       *sqlx.DB
	APIKey   string
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())
	// Pass the APIKey to the handleAsk function
	// Pass the APIKey to the handleAsk function
	app.Handle(http.MethodPost, "/ask", aigrp.HandleAskWrapper(cfg.APIKey))

	app.Handle(http.MethodGet, "/test", testgrp.Test)
	app.Handle(http.MethodGet, "/test/auth", testgrp.Test, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))

	// -------------------------------------------------------------------------

	usrCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))

	ugh := usergrp.New(usrCore)

	app.Handle(http.MethodGet, "/users", ugh.Query)

	return app
}

func handleAddModeration(db *sqlx.DB, log *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) {
	CheckPOSTMethod(w, r)

	var mdr moderator.Moderator
	if err := DecodeRequestBody(r, w, &mdr); err != nil {
		return
	}

	err := moderationgrp.AddModerationToDB(log, db, mdr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendSuccess(w)
	w.Write([]byte("Moderation data added successfully"))
}

// DecodeRequestBody decodes the request body into the provided interface.
func DecodeRequestBody[T any](r *http.Request, w http.ResponseWriter, v *T) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return err
}

func CheckPOSTMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
	}
}

// SendSuccess sends a successful response.
func SendSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
