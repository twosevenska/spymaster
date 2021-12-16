package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
