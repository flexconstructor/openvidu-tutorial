package action

import (
	"testing"

	"github.com/flexconstructor/openvidu-tutorial/repository"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUser_Get(t *testing.T) {
	a := &User{
		UsersRepo: repository.NewUsersRepository(),
	}

	Convey("Returns user value object", t, func() {
		a.UsersRepo.Add("test user", "test user password", 1)
		user, err := a.Get("test user")

		Convey("Returns no errors", func() {
			So(err, ShouldBeNil)
		})

		Convey("Returns user", func() {
			So(user, ShouldNotBeNil)

			Convey("with correct user data", func() {
				So(user.Name, ShouldEqual, "test user")
				So(user.Password, ShouldEqual, "test user password")
				So(user.Role, ShouldEqual, 1)
			})
		})
	})

	Convey("Returns an error", t, func() {
		_, err := a.Get("wrong user")

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "login incorrect")
	})
}
