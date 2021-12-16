package controllers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"

	"internal/mongo"
	"internal/server"
	"types"
)

// errorResponse represents an API response error
type errorResponse struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

var (
	mc *mongo.Client
	r  *gin.Engine
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestPing(t *testing.T) {
	Convey("When a ping is sent to the API", t, func() {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/ping", nil)
		So(err, ShouldBeNil)

		var result errorResponse
		serveAndUnmarshal(recorder, req, &result)
		resp := recorder.Result()

		Convey("The response should be a JSON-formatted pong response", func() {
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(result, ShouldResemble, errorResponse{Message: "pong"})
		})
	})
}

func TestHealth(t *testing.T) {
	Convey("When a /health is sent to the API", t, func() {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/health", nil)
		So(err, ShouldBeNil)

		var result errorResponse
		serveAndUnmarshal(recorder, req, &result)
		resp := recorder.Result()

		Convey("The response should be a JSON-formatted UP response", func() {
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(result, ShouldResemble, errorResponse{Message: "UP"})
		})
	})
}

// Utils
func setup() {

	// Use a test database as out experimental playground
	confMongo := mongo.Config{
		Hosts:    []string{"0.0.0.0:27017"},
		Database: "test_spymaster",
		User:     "test",
		Password: "test",
	}
	var err error

	mc, err = mongo.Connect(confMongo)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	r = server.CreateRouter(mc)
	cleanUp()
}

func shutdown() {
	cleanUp()
	mc.Close()
}

func cleanUp() {
	// Clean up the MongoDB collections
	_, err := mc.Database.C("users").RemoveAll(bson.M{})
	if err != nil {
		log.Fatalf("Failed cleaning up MongoDB for tests")
		return
	}
}

func withCleanup(f func()) func() {
	return func() {

		Reset(func() {
			cleanUp()
		})

		f()
	}
}

func serveAndUnmarshal(rec *httptest.ResponseRecorder, req *http.Request, result interface{}) {
	r.ServeHTTP(rec, req)
	b := rec.Body.Bytes()
	if err := json.Unmarshal(b, result); err != nil {
		panic(err.Error())
	}
}

func createUser(u types.User) (nu *types.User, err error) {
	now := time.Now().UTC()
	u.ID = bson.NewObjectId()
	u.CreatedAt = now
	u.UpdatedAt = now

	err = mc.Database.C("users").Insert(u)
	if err == nil {
		nu = &u
	} else {
		log.Printf("Error inserting user: %s", err)
	}

	return
}

func getDBUser(userID string) (u types.User, err error) {
	id := bson.ObjectIdHex(userID)

	criteria := bson.M{"_id": id}
	err = mc.Database.C("users").Find(criteria).One(&u)
	if err != nil {
		return
	}

	return
}
