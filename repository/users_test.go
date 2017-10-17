package repository

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUsersRepository(t *testing.T) {
	Convey("Returns new instance of repository", t, func() {
		r := NewUsersRepository()
		So(r, ShouldNotBeNil)

		Convey("Repository storage is not nil", func() {
			So(r.users, ShouldNotBeNil)
		})
	})
}

func TestUsers_Add(t *testing.T) {
	Convey("Add new user to repository", t, func() {
		r := NewUsersRepository()
		r.Add("test login", "test password", 1)
		So(r.users["test login"], ShouldNotBeNil)

		Convey("with correct user data", func() {
			So(r.users["test login"].Name, ShouldEqual, "test login")
			So(r.users["test login"].Password, ShouldEqual, "test password")
			So(r.users["test login"].Role, ShouldEqual, 1)
		})
	})
}

func TestUsers_Get(t *testing.T) {
	Convey("Returns user", t, func() {
		r := NewUsersRepository()
		r.Add("test login", "test password", 1)

		Convey("with correct user data", func() {
			user, _ := r.Get("test login")
			So(user.Name, ShouldEqual, "test login")
			So(user.Password, ShouldEqual, "test password")
			So(user.Role, ShouldEqual, 1)
		})

		Convey("Returns no error", func() {
			_, err := r.Get("test login")
			So(err, ShouldBeNil)
		})

		Convey("Returns an error", func() {
			_, err := r.Get("wrong login")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "login incorrect")
		})
	})
}
