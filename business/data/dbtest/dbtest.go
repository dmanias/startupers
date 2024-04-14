// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"math/rand"
	"net/mail"
	"testing"
	"time"

	"github.com/dmanias/startupers/business/core/user"
	"github.com/dmanias/startupers/business/core/user/stores/userdb"
	database "github.com/dmanias/startupers/business/sys/database/pgx"
	"github.com/dmanias/startupers/business/web/auth"
	"github.com/dmanias/startupers/foundation/docker"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StartDB starts a database instance.
func StartDB() (*docker.Container, error) {
	image := "postgres:15.3"
	port := "5432"
	args := []string{"-e", "POSTGRES_PASSWORD=postgres"}

	c, err := docker.StartContainer(image, port, args...)
	if err != nil {
		return nil, fmt.Errorf("starting container: %w", err)
	}

	fmt.Printf("Image:       %s\n", image)
	fmt.Printf("ContainerID: %s\n", c.ID)
	fmt.Printf("Host:        %s\n", c.Host)

	return c, nil
}

// StopDB stops a running database instance.
func StopDB(c *docker.Container) {
	docker.StopContainer(c.ID)
	fmt.Println("Stopped:", c.ID)
}

// =============================================================================

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *sqlx.DB
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	CoreAPIs CoreAPIs
	Teardown func()
	t        *testing.T
}

// NewTest creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewTest(t *testing.T, c *docker.Container) *Test {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbM, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	if err := database.StatusCheck(ctx, dbM); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	dbName := string(b)

	if _, err := dbM.ExecContext(context.Background(), "CREATE DATABASE "+dbName); err != nil {
		t.Fatalf("creating database %s: %v", dbName, err)
	}
	dbM.Close()

	t.Log("Database ready")

	// -------------------------------------------------------------------------

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       dbName,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Migrate and seed database ...")

	//if err := dbmigrate.Migrate(ctx, db); err != nil {
	//	t.Logf("Logs for %s\n%s:", c.ID, docker.DumpContainerLogs(c.ID))
	//	t.Fatalf("Migrating error: %s", err)
	//}
	//
	//if err := dbmigrate.Seed(ctx, db); err != nil {
	//	t.Logf("Logs for %s\n%s:", c.ID, docker.DumpContainerLogs(c.ID))
	//	t.Fatalf("Seeding error: %s", err)
	//}

	// -------------------------------------------------------------------------

	var buf bytes.Buffer
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buf)
	log := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel),
		zap.WithCaller(true),
	).Sugar()

	coreAPIs := newCoreAPIs(log, db)

	t.Log("Ready for testing ...")

	// -------------------------------------------------------------------------

	cfg := auth.Config{
		Log:       log,
		KeyLookup: &keyStore{},
	}
	a, err := auth.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()

		log.Sync()

		writer.Flush()
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	test := Test{
		DB:       db,
		Log:      log,
		Auth:     a,
		CoreAPIs: coreAPIs,
		Teardown: teardown,
		t:        t,
	}

	return &test
}

// Token generates an authenticated token for a user.
func (test *Test) Token(email string, pass string) string {
	test.t.Log("Generating token for test ...")

	addr, _ := mail.ParseAddress(email)

	store := userdb.NewStore(test.Log, test.DB)
	dbUsr, err := store.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return ""
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: dbUsr.Roles,
	}

	token, err := test.Auth.GenerateToken(kid, claims)
	if err != nil {
		test.t.Fatal(err)
	}

	return token
}

// =============================================================================

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// FloatPointer is a helper to get a *float64 from a float64. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func FloatPointer(f float64) *float64 {
	return &f
}

// =============================================================================

// CoreAPIs represents all the core api's needed for testing.
type CoreAPIs struct {
	User *user.Core
}

func newCoreAPIs(log *zap.SugaredLogger, db *sqlx.DB) CoreAPIs {
	usrCore := user.NewCore(userdb.NewStore(log, db))

	return CoreAPIs{
		User: usrCore,
	}
}

// =============================================================================

type keyStore struct{}

func (ks *keyStore) PrivateKey(kid string) (string, error) {
	return privateKeyPEM, nil
}

func (ks *keyStore) PublicKey(kid string) (string, error) {
	return publicKeyPEM, nil
}

// =============================================================================

const (
	kid = ""

	privateKeyPEM = ``
	publicKeyPEM  = ``
)
