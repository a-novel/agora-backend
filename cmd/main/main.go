package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/a-novel/agora-backend/api"
	"github.com/a-novel/agora-backend/api/base"
	"github.com/a-novel/agora-backend/api/bookmark"
	"github.com/a-novel/agora-backend/api/forum"
	"github.com/a-novel/agora-backend/api/secrets"
	"github.com/a-novel/agora-backend/api/user"
	"github.com/a-novel/agora-backend/config"
	"github.com/a-novel/agora-backend/domains/bookmark/service/improve_post"
	"github.com/a-novel/agora-backend/domains/bookmark/storage/improve_post"
	"github.com/a-novel/agora-backend/domains/forum/service/improve_request"
	"github.com/a-novel/agora-backend/domains/forum/service/improve_suggestion"
	"github.com/a-novel/agora-backend/domains/forum/service/votes"
	"github.com/a-novel/agora-backend/domains/forum/storage/improve_request"
	"github.com/a-novel/agora-backend/domains/forum/storage/improve_suggestion"
	"github.com/a-novel/agora-backend/domains/forum/storage/votes"
	"github.com/a-novel/agora-backend/domains/generics"
	"github.com/a-novel/agora-backend/domains/keys/service/jwk"
	"github.com/a-novel/agora-backend/domains/keys/storage/jwk"
	"github.com/a-novel/agora-backend/domains/user/service/credentials"
	"github.com/a-novel/agora-backend/domains/user/service/identity"
	"github.com/a-novel/agora-backend/domains/user/service/profile"
	"github.com/a-novel/agora-backend/domains/user/service/token"
	"github.com/a-novel/agora-backend/domains/user/service/user"
	"github.com/a-novel/agora-backend/domains/user/storage/credentials"
	"github.com/a-novel/agora-backend/domains/user/storage/identity"
	"github.com/a-novel/agora-backend/domains/user/storage/profile"
	"github.com/a-novel/agora-backend/domains/user/storage/user"
	improve_post_bookmark "github.com/a-novel/agora-backend/environment/bookmark/improve_post"
	improve_post_forum "github.com/a-novel/agora-backend/environment/forum/improve_post"
	"github.com/a-novel/agora-backend/environment/secrets"
	"github.com/a-novel/agora-backend/environment/user/account"
	"github.com/a-novel/agora-backend/environment/user/authentication"
	"github.com/a-novel/agora-backend/environment/user/profile"
	"github.com/a-novel/agora-backend/framework/bunframework"
	"github.com/a-novel/agora-backend/framework/bunframework/pgconfig"
	"github.com/a-novel/agora-backend/framework/mailer"
	"github.com/a-novel/agora-backend/framework/security"
	"github.com/a-novel/agora-backend/migrations"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
	"io/fs"
	"os"
	"path"
	"time"
)

func FrontendURL(path string) generics.URL {
	cfg := config.Main()
	return generics.URL{
		Host: cfg.Frontend.URLs[0],
		Path: path,
	}
}

func main() {
	cfg, env := config.Main(), config.Env()

	fmt.Println(cfg.Frontend.URLs)

	corsConfig := cors.Config{
		AllowOrigins: cfg.Frontend.URLs,
		AllowMethods: cors.DefaultConfig().AllowMethods,
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Access-Control-Allow-Origin",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	// Setup logger.
	logger := zerolog.New(os.Stdout).
		With().
		Dict(
			"application", zerolog.Dict().
				Str("name", cfg.App).
				Str("env", env).
				Dict(
					"cors", zerolog.Dict().
						Strs("allowed_origins", cfg.Frontend.URLs).
						Strs("allowed_methods", corsConfig.AllowMethods).
						Strs("allowed_headers", corsConfig.AllowHeaders).
						Strs("exposed_headers", corsConfig.ExposeHeaders).
						Bool("allow_credentials", corsConfig.AllowCredentials).
						Dur("max_age", corsConfig.MaxAge),
				).
				Str("host", cfg.API.Host).
				Int("port", cfg.API.Port),
		).
		Dict(
			"mailer", zerolog.Dict().
				Str("sender_email", cfg.Mailer.Sender.Email).
				Str("sender_name", cfg.Mailer.Sender.Name).
				Bool("sandbox", cfg.Mailer.Sandbox),
		).
		Logger()
	switch env {
	case config.ENVProduction:
		logger = logger.With().Timestamp().Logger()
	default:
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Setup postgresql.
	postgres, sql, err := bunframework.NewClient(context.Background(), bunConfig(cfg))
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = postgres.Close()
		_ = sql.Close()
	}()

	// Setup mailer.
	mailSender := mail.NewEmail(cfg.Mailer.Sender.Name, cfg.Mailer.Sender.Email)
	mailClient := mailer.NewMailer(cfg.Mailer.APIKey, mailSender, cfg.Mailer.Sandbox, logger)

	// Setup repositories.
	keysRepository, logger, err := startKeysRepository(cfg, logger, env)
	if err != nil {
		panic(err.Error())
	}

	userCredentialsRepository := credentials_storage.NewRepository(postgres)
	userIdentityRepository := identity_storage.NewRepository(postgres)
	userProfileRepository := profile_storage.NewRepository(postgres)
	userRepository := user_storage.NewRepository(postgres)

	forumImproveRequestRepository := improve_request_storage.NewRepository(postgres, cfg.Forum.Search.CropContent)
	forumImproveSuggestionRepository := improve_suggestion_storage.NewRepository(postgres, cfg.Forum.Search.CropContent)
	forumVotesRepository := votes_storage.NewRepository(postgres)

	bookmarkImprovePostRepository := improve_post_storage.NewRepository(postgres)

	// Setup services.
	keysService := jwk_service.NewService(keysRepository)
	keysServiceCached := keysService.ReadOnly()
	tokenService := token_service.NewService()

	userCredentialsService := credentials_service.NewService(
		userCredentialsRepository,
		security.GenerateCode,
		security.VerifyCode,
		bcrypt.GenerateFromPassword,
		bcrypt.CompareHashAndPassword,
	)
	userIdentityService := identity_service.NewService(userIdentityRepository)
	userProfileService := profile_service.NewService(userProfileRepository)
	userService := user_service.NewService(
		userRepository,
		userCredentialsService,
		userIdentityService,
		userProfileService,
	)

	forumImproveRequestService := improve_request_service.NewService(forumImproveRequestRepository)
	forumImproveSuggestionService := improve_suggestion_service.NewService(forumImproveSuggestionRepository)
	forumVotesService := votes_service.NewService(forumVotesRepository)

	bookmarkImprovePostService := improve_post_service.NewService(bookmarkImprovePostRepository)

	// Setup providers.
	secretsProvider := secrets.NewProvider(secrets.Config{
		KeysService:    keysService,
		KeyGen:         security.JWKKeyGen,
		Now:            time.Now,
		ID:             uuid.New,
		MaxBackups:     cfg.Secrets.Backups,
		UpdateInterval: cfg.Secrets.UpdateInterval,
	})

	accountProvider := account.NewProvider(account.Config{
		CredentialsService:         userCredentialsService,
		IdentityService:            userIdentityService,
		ProfileService:             userProfileService,
		UserService:                userService,
		TokenService:               tokenService,
		KeysService:                keysServiceCached,
		Mailer:                     mailClient,
		Time:                       time.Now,
		ID:                         uuid.New,
		TokenTTL:                   cfg.Tokens.TTL,
		TokenRenewDelta:            cfg.Tokens.RenewDelta,
		EmailValidationLink:        FrontendURL(cfg.Frontend.Routes.ValidateEmail),
		NewEmailValidationLink:     FrontendURL(cfg.Frontend.Routes.ValidateNewEmail),
		PasswordResetLink:          FrontendURL(cfg.Frontend.Routes.ResetPassword),
		EmailValidationTemplate:    cfg.Mailer.Templates.EmailValidation,
		NewEMailValidationTemplate: cfg.Mailer.Templates.EmailUpdate,
		PasswordResetTemplate:      cfg.Mailer.Templates.PasswordReset,
	})
	authenticationProvider := authentication.NewProvider(authentication.Config{
		CredentialsService: userCredentialsService,
		TokenService:       tokenService,
		KeysService:        keysServiceCached,
		Time:               time.Now,
		ID:                 uuid.New,
		TokenTTL:           cfg.Tokens.TTL,
		TokenRenewDelta:    cfg.Tokens.RenewDelta,
	})
	profileProvider := profile.NewProvider(profile.Config{
		UserService: userService,
	})

	forumImprovePostProvider := improve_post_forum.NewProvider(improve_post_forum.Config{
		ImproveRequestService:    forumImproveRequestService,
		ImproveSuggestionService: forumImproveSuggestionService,
		VotesService:             forumVotesService,
		TokenService:             tokenService,
		KeysService:              keysServiceCached,
		UserService:              userService,
		Time:                     time.Now,
		ID:                       uuid.New,
	})

	bookmarkImprovePostProvider := improve_post_bookmark.NewProvider(improve_post_bookmark.Config{
		BookmarkService: bookmarkImprovePostService,
		TokenService:    tokenService,
		KeysService:     keysServiceCached,
		UserService:     userService,
		Time:            time.Now,
	})

	// Refresh cache once at startup, to have keys loaded (otherwise the handler will be empty and unable to
	// authenticate the first request).
	if err := keysServiceCached.RefreshCache(context.Background()); err != nil {
		panic(err.Error())
	}

	// Setup API.
	router := gin.New()
	loadGinMiddlewares(cfg, corsConfig, logger, keysServiceCached, router)
	apiRouter := router.Group("/api")

	if env == config.ENVProduction {
		gin.SetMode(gin.ReleaseMode)
		router.TrustedPlatform = gin.PlatformGoogleAppEngine
	}

	baseapi.API(apiRouter)

	secretsapi.JWKsAPI("/secrets", apiRouter, secretsProvider, cfg.IAM.ServiceAccounts.Scheduler)

	userapi.AuthenticationAPI("/user/auth", apiRouter, authenticationProvider)
	userapi.AccountAPI("/user/account", apiRouter, accountProvider)
	userapi.ProfileAPI("/user/profile", apiRouter, profileProvider)

	forumapi.ImproveRequestAPI("/forum/improve-request", apiRouter, forumImprovePostProvider)
	forumapi.ImproveSuggestionAPI("/forum/improve-suggestion", apiRouter, forumImprovePostProvider)
	forumapi.VotesAPI("/forum/votes", apiRouter, forumImprovePostProvider)

	bookmarkapi.ImprovePostAPI("/bookmark/improve-post", apiRouter, bookmarkImprovePostProvider)

	// Run API.
	if err := router.Run(fmt.Sprintf(":%d", cfg.API.Port)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running API, and the server had to shut down")
	}
}

func bunConfig(cfg *config.Config) bunframework.Config {
	return bunframework.Config{
		Driver: pgconfig.Driver{
			DSN:         cfg.Postgres.DSN,
			AppName:     cfg.App,
			DialTimeout: 120 * time.Second,
		},
		DiscardUnknownColumns: true,
		Migrations: &bunframework.MigrateConfig{
			Files: []fs.FS{migrations.Migrations},
		},
	}
}

func startKeysRepository(cfg *config.Config, logger zerolog.Logger, env string) (jwk_storage.Repository, zerolog.Logger, error) {
	if env == config.ENVProduction {
		client, err := storage.NewClient(context.Background())
		if err != nil {
			return nil, logger, err
		}

		logger = logger.With().
			Dict(
				"secrets_manager",
				zerolog.Dict().
					Int("backups", cfg.Secrets.Backups).
					Dur("update_interval", cfg.Secrets.UpdateInterval).
					Str("type", "GCP Datastore"),
			).
			Logger()
		return jwk_storage.NewGoogleDatastoreRepository(client.Bucket(cfg.Buckets.SecretKeys)), logger, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, logger, err
	}

	keysPath := path.Join(wd, ".secrets")
	logger = logger.With().
		Dict(
			"secrets_manager",
			zerolog.Dict().
				Int("backups", cfg.Secrets.Backups).
				Dur("update_interval", cfg.Secrets.UpdateInterval).
				Str("type", "local storage").
				Str("path", keysPath).
				Str("prefix", cfg.Secrets.Prefix),
		).
		Logger()
	return jwk_storage.NewFileSystemRepository(keysPath, cfg.Secrets.Prefix), logger, nil
}

func loadGinMiddlewares(cfg *config.Config, corsConfig cors.Config, logger zerolog.Logger, keysService jwk_service.ServiceCached, r gin.IRouter) {
	r.Use(
		gin.RecoveryWithWriter(logger),
		api.Logger(logger, cfg.ProjectID),
		cors.New(corsConfig),
		func(c *gin.Context) {
			c.Next()

			// Update keys after response was sent, to avoid impact on performances.
			if err := keysService.RefreshCache(c); err != nil {
				_ = c.Error(err)
			}
		},
	)
}
