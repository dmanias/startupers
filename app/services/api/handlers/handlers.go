package handlers

import (
	"context"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/aigrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/challengegrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/checkgrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/ideagrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/moderationgrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/postgrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/testgrp"
	"github.com/dmanias/startupers/app/services/api/handlers/v1/usergrp"
	"github.com/dmanias/startupers/business/core/ai"
	"github.com/dmanias/startupers/business/core/ai/stores/aidb"
	"github.com/dmanias/startupers/business/core/challenge"
	challengedb "github.com/dmanias/startupers/business/core/challenge/stores/challengedb"
	"github.com/dmanias/startupers/business/core/idea"
	"github.com/dmanias/startupers/business/core/idea/stores/ideadb"
	"github.com/dmanias/startupers/business/core/moderator"
	"github.com/dmanias/startupers/business/core/moderator/stores/moderatordb"
	"github.com/dmanias/startupers/business/core/post"
	postdb "github.com/dmanias/startupers/business/core/post/stores/postdb"
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
	Shutdown      chan os.Signal
	Log           *zap.SugaredLogger
	Auth          *auth.Auth
	AuthConfig    *auth.Config
	DB            *sqlx.DB
	APIKey        string
	Build         string
	ModeratorCore *moderator.Core
	ActiveKID     string
	AIType        string
	APIHost       string
	//GoogleOauthConfig *oauth2.Config
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())
	if cfg.Log == nil {
		panic("cfg.Log is nil")
	}
	if cfg.Auth == nil {
		panic("cfg.Auth is nil")
	}
	if cfg.DB == nil {
		panic("cfg.DB is nil")
	}
	if cfg.AuthConfig == nil {
		panic("cfg.AuthConfig is nil")
	}

	// Initialize the moderator.Core and moderationgrp.Handlers instances
	moderatorCore := moderator.NewCore(moderatordb.NewStore(cfg.Log, cfg.DB))

	aigrpCfg := aigrp.APIMuxConfig{
		Shutdown:      cfg.Shutdown,
		Log:           cfg.Log,
		DB:            cfg.DB,
		APIKey:        cfg.APIKey,
		Build:         cfg.Build,
		ModeratorCore: moderatorCore,
	}
	// Initialize the ai.Core and aigrp.Handlers instances
	aiCore := ai.NewCore(aidb.NewStore(cfg.Log, cfg.DB))
	mgh := moderationgrp.New(moderatorCore)
	cfg.Log.Info("cfg.Auth", cfg.Auth)
	// Initialize the post.Core and challengegrp.Handlers instances
	postCore := post.NewCore(postdb.NewStore(cfg.Log, cfg.DB))
	ideaCore := idea.NewCore(ideadb.NewStore(cfg.Log, cfg.DB))
	aiHandlers := aigrp.New(aiCore, aigrpCfg, mgh, postCore, ideaCore)
	ideaHandlers := ideagrp.New(ideaCore, cfg.Log, aiHandlers, mgh, cfg.APIHost)
	// Update the aigrp.New function call to include ideaCore and postCore
	postHandlers := postgrp.New(postCore, cfg.Log, aiHandlers, mgh)

	// Add the routes for idea-related operations
	app.Handle(http.MethodPost, "/ideas", ideaHandlers.Create, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/ideas/:idea_id", ideaHandlers.QueryByID, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	//app.Handle(http.MethodPut, "/ideas/:idea_id", ideaHandlers.Update, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodDelete, "/ideas/:idea_id", ideaHandlers.Delete, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/:user_id/ideas", ideaHandlers.Query, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/tags", ideaHandlers.QueryTags, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	// Add the routes for post-related operations
	app.Handle(http.MethodPost, "/ideas/posts", postHandlers.Create, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/posts/:post_id", postHandlers.QueryByID, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	//app.Handle(http.MethodPut, "/posts/:post_id", postHandlers.Update, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodDelete, "/posts/:post_id", postHandlers.Delete, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/ideas/:idea_id/posts", postHandlers.Query, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))

	// Add the routes for moderator-related operations
	app.Handle(http.MethodPost, "/moderators", mgh.Create, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))
	app.Handle(http.MethodGet, "/moderators/:name", mgh.QueryByNameHandler, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/moderators", mgh.Query, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))

	app.Handle(http.MethodGet, "/ask/:idea_id/:question_type/:description", aiHandlers.Ask, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))

	///app.Handle(http.MethodGet, "/ask/:scenario/idea_id", aiHandlers.Ask, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	// Create a handlers instance for checkgrp
	hdl := checkgrp.New(cfg.Build, cfg.Log, cfg.DB)

	// Add the readiness and liveness routes
	app.HandleNoMiddleware(http.MethodGet, "/readiness", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return hdl.Readiness(ctx, w, r)
	})
	app.HandleNoMiddleware(http.MethodGet, "/liveness", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return hdl.Liveness(ctx, w, r)
	})

	app.Handle(http.MethodGet, "/test", testgrp.Test, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))
	app.Handle(http.MethodGet, "/test/auth", testgrp.Test, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))

	// -------------------------------------------------------------------------
	usrCore := user.NewCore(userdb.NewStore(cfg.Log, cfg.DB))
	authConfig := auth.Config{
		Log:       cfg.Log, // or another *zap.SugaredLogger instance
		KeyLookup: cfg.AuthConfig.KeyLookup,
		Issuer:    cfg.AuthConfig.Issuer,
	}
	authInstance, err := auth.New(authConfig)
	if err != nil {
		cfg.Log.Errorf("Failed to create auth instance: %v", err)
		return nil
	}
	ugh := usergrp.New(usrCore, authInstance, cfg.ActiveKID, cfg.Log)
	app.Handle(http.MethodPost, "/users/login", ugh.Login)
	app.Handle(http.MethodPost, "/users/register", ugh.Create)
	//app.Handle(http.MethodGet, "/users", ugh.Query, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleAdminOnly))

	// Serve static files from the "uploads" directory
	fs := http.FileServer(http.Dir("./uploads"))
	app.Handle(http.MethodGet, "/uploads/*", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		http.StripPrefix("/uploads/", fs).ServeHTTP(w, r)
		return nil
	})

	//-------Challenge-------
	// Initialize the challenge.Core and challengegrp.Handlers instances
	challengeCore := challenge.NewCore(challengedb.NewStore(cfg.Log, cfg.DB))
	challengeHandlers := challengegrp.New(challengeCore, cfg.Log)

	// Add the routes for challenge-related operations
	app.Handle(http.MethodPost, "/ideas/challenges", challengeHandlers.Create, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	//app.Handle(http.MethodGet, "/challenges/:challenge_id", challengeHandlers.QueryByID, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodPut, "/challenges/:challenge_id", challengeHandlers.Update, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodDelete, "/challenges/:challenge_id", challengeHandlers.Delete, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))
	app.Handle(http.MethodGet, "/ideas/:idea_id/challenges", challengeHandlers.Query, mid.Authenticate(cfg.Auth), mid.Authorize(cfg.Auth, auth.RuleUserOnly))

	//----Auth-----
	// Initialize the authgrp.Handlers instance
	//authHandlers := authgrp.New(cfg.Auth)
	//
	//// Add the routes for token generation, authentication, and authorization
	//app.Handle(http.MethodGet, "/auth/token/:kid", authHandlers.Token)
	//app.Handle(http.MethodGet, "/auth/authenticate", authHandlers.Authenticate)
	//app.Handle(http.MethodPost, "/auth/authorize", authHandlers.Authorize)

	return app
}

//TODO add authentication to admin routes
