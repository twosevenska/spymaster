package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"

	"types"
)

func TestListUsers(t *testing.T) {
	Convey("When a list of Users is requested from the API...", t, withCleanup(func() {
		recorder := httptest.NewRecorder()

		Convey("And there are no users.", func() {
			req, err := http.NewRequest("GET", "/users", nil)
			So(err, ShouldBeNil)

			var result types.UsersResult
			serveAndUnmarshal(recorder, req, &result)
			resp := recorder.Result()

			Convey("The response should be an empty JSON array", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(result.Page, ShouldEqual, 1)
				So(result.PerPage, ShouldEqual, 100)
				So(result.TotalCount, ShouldEqual, 0)
				So(result.Users, ShouldHaveLength, 0)
			})
		})

		Convey("And there are users...", func() {

			u1 := types.User{
				FirstName: "Yellow",
				LastName:  "King",
				Nickname:  "hastur",
				Password:  "Carcosa",
				Email:     "hastur@lost.space",
				Country:   "UK",
			}
			createUser(u1)

			u2 := types.User{
				FirstName: "Blue",
				LastName:  "King",
				Nickname:  "fake_hastur",
				Password:  "Carcosa",
				Email:     "fake_hastur@lost.space",
				Country:   "UK",
			}
			createUser(u2)

			u3 := types.User{
				FirstName: "Robin",
				LastName:  "Williams",
				Nickname:  "genie",
				Password:  "Jumanji",
				Email:     "rwilliams@hollywood.fake",
				Country:   "US",
			}
			createUser(u3)

			Convey("And we don't filter...", func() {
				Convey("Nor paginate", func() {
					req, err := http.NewRequest("GET", "/users", nil)
					So(err, ShouldBeNil)

					var result types.UsersResult
					serveAndUnmarshal(recorder, req, &result)
					resp := recorder.Result()

					Convey("The response should be an empty JSON array", func() {
						So(resp.StatusCode, ShouldEqual, http.StatusOK)
						So(result.Page, ShouldEqual, 1)
						So(result.PerPage, ShouldEqual, 100)
						So(result.TotalCount, ShouldEqual, 3)
						So(result.Users, ShouldHaveLength, 3)
					})
				})

				Convey("Paginate", func() {
					Convey("First page", func() {
						req, err := http.NewRequest("GET", "/users?per_page=1", nil)
						So(err, ShouldBeNil)

						var result types.UsersResult
						serveAndUnmarshal(recorder, req, &result)
						resp := recorder.Result()

						Convey("The response should be an empty JSON array", func() {
							So(resp.StatusCode, ShouldEqual, http.StatusOK)
							So(result.Page, ShouldEqual, 1)
							So(result.PerPage, ShouldEqual, 1)
							So(result.TotalCount, ShouldEqual, 3)
							So(result.Users, ShouldHaveLength, 1)
						})
					})
					Convey("Second page", func() {
						req, err := http.NewRequest("GET", "/users?per_page=1&page=2", nil)
						So(err, ShouldBeNil)

						var result types.UsersResult
						serveAndUnmarshal(recorder, req, &result)
						resp := recorder.Result()

						Convey("The response should be an empty JSON array", func() {
							So(resp.StatusCode, ShouldEqual, http.StatusOK)
							So(result.Page, ShouldEqual, 2)
							So(result.PerPage, ShouldEqual, 1)
							So(result.TotalCount, ShouldEqual, 3)
							So(result.Users, ShouldHaveLength, 1)
						})
					})
				})
			})

			Convey("And we filter...", func() {
				Convey("with partial field.", func() {
					req, err := http.NewRequest("GET", "/users?nickname=hast", nil)
					So(err, ShouldBeNil)

					var result types.UsersResult
					serveAndUnmarshal(recorder, req, &result)
					resp := recorder.Result()

					Convey("The response should be an empty JSON array", func() {
						So(resp.StatusCode, ShouldEqual, http.StatusOK)
						So(result.Page, ShouldEqual, 1)
						So(result.PerPage, ShouldEqual, 100)
						So(result.TotalCount, ShouldEqual, 2)
						So(result.Users, ShouldHaveLength, 2)
					})
				})

				Convey("with exact field.", func() {
					req, err := http.NewRequest("GET", "/users?country=US", nil)
					So(err, ShouldBeNil)

					var result types.UsersResult
					serveAndUnmarshal(recorder, req, &result)
					resp := recorder.Result()

					Convey("The response should be an empty JSON array", func() {
						So(resp.StatusCode, ShouldEqual, http.StatusOK)
						So(result.Page, ShouldEqual, 1)
						So(result.PerPage, ShouldEqual, 100)
						So(result.TotalCount, ShouldEqual, 1)
						So(result.Users, ShouldHaveLength, 1)
					})
				})

				Convey("with both an exact and partial field...", func() {
					req, err := http.NewRequest("GET", "/users?nickname=ge&country=US", nil)
					So(err, ShouldBeNil)

					var result types.UsersResult
					serveAndUnmarshal(recorder, req, &result)
					resp := recorder.Result()

					Convey("The response should be an empty JSON array", func() {
						So(resp.StatusCode, ShouldEqual, http.StatusOK)
						So(result.Page, ShouldEqual, 1)
						So(result.PerPage, ShouldEqual, 100)
						So(result.TotalCount, ShouldEqual, 1)
						So(result.Users, ShouldHaveLength, 1)
					})
				})
			})
		})
	}))
}

func TestCreateUser(t *testing.T) {
	Convey("When a new user is created through the API...", t, withCleanup(func() {
		recorder := httptest.NewRecorder()

		Convey("And the user does not exist", func() {
			Convey("And the payload was valid", func() {

				payload := map[string]interface{}{
					"first_name": "Robin",
					"last_name":  "Williams",
					"nickname":   "genie",
					"password":   "Jumanji",
					"email":      "rwilliams@hollywood.fake",
					"country":    "US",
				}
				p, err := json.Marshal(payload)
				if err != nil {
					log.Printf("server_test: Error creating POST request: %s", err)
				}

				req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(p))
				So(err, ShouldBeNil)

				var result types.User
				serveAndUnmarshal(recorder, req, &result)
				resp := recorder.Result()

				Convey("Ensure the user was properly created", func() {
					So(resp.StatusCode, ShouldEqual, http.StatusCreated)
					u, err := getDBUser(result.ID.Hex())
					So(err, ShouldBeNil)
					So(u.Nickname, ShouldEqual, payload["nickname"])
				})
			})

			Convey("And the payload was not valid", func() {

				payload := map[string]interface{}{
					"first_name": "Robin",
					"last_name":  "Williams",
					"OPS":        "OPS",
				}
				p, err := json.Marshal(payload)
				if err != nil {
					log.Printf("server_test: Error creating POST request: %s", err)
				}

				req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(p))
				So(err, ShouldBeNil)

				r.ServeHTTP(recorder, req)
				resp := recorder.Result()
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("And the user exists", func() {
			u2 := types.User{
				FirstName: "Blue",
				LastName:  "King",
				Nickname:  "fake_hastur",
				Password:  "Carcosa",
				Email:     "fake_hastur@lost.space",
				Country:   "UK",
			}
			createUser(u2)

			payload := map[string]interface{}{
				"first_name": "Blue",
				"last_name":  "King",
				"nickname":   "fake_hastur",
				"password":   "Carcosa",
				"email":      "fake_hastur@lost.space",
				"country":    "UK",
			}
			p, err := json.Marshal(payload)
			if err != nil {
				log.Printf("server_test: Error creating POST request: %s", err)
			}

			req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(p))
			So(err, ShouldBeNil)

			var result types.User
			serveAndUnmarshal(recorder, req, &result)
			resp := recorder.Result()
			So(resp.StatusCode, ShouldEqual, http.StatusConflict)
		})
	}))
}

func TestUpdateUser(t *testing.T) {
	Convey("When a list of Users is requested from the API...", t, withCleanup(func() {
		recorder := httptest.NewRecorder()

		Convey("And there are no users.", func() {
			req, err := http.NewRequest("GET", "/users", nil)
			So(err, ShouldBeNil)

			var result types.UsersResult
			serveAndUnmarshal(recorder, req, &result)
			resp := recorder.Result()

			Convey("The response should be an empty JSON array", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(result.Page, ShouldEqual, 1)
				So(result.PerPage, ShouldEqual, 100)
				So(result.TotalCount, ShouldEqual, 0)
				So(result.Users, ShouldHaveLength, 0)
			})
		})
	}))
}

func TestDeleteUser(t *testing.T) {
	Convey("When a delete request comes in...", t, withCleanup(func() {
		recorder := httptest.NewRecorder()

		Convey("And the user does not exist", func() {
			req, err := http.NewRequest("DELETE", "/users?id=61ba6382df4bec585cf60e60", nil)
			So(err, ShouldBeNil)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()
			So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
		})

		Convey("And the user exists", func() {
			user1 := types.User{
				FirstName: "Yellow",
				LastName:  "King",
				Nickname:  "hastur",
				Password:  "Carcosa",
				Email:     "hastur@lost.space",
				Country:   "UK",
			}
			dbUser, _ := createUser(user1)
			userId := dbUser.ID.Hex()

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/users?id=%s", userId), nil)
			So(err, ShouldBeNil)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()
			So(resp.StatusCode, ShouldEqual, http.StatusNoContent)

			Convey("Ensure the user was deleted in the DB", func() {
				_, err := getDBUser(userId)
				So(err, ShouldEqual, mgo.ErrNotFound)
			})
		})
	}))
}
