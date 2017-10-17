package action

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flexconstructor/openvidu-tutorial/repository"
)

func TestLogin_Do(t *testing.T) {
	r := repository.NewUsersRepository()
	r.Add("test login", "test password", 1)
	a := &Login{UserRepo: r}

	Convey("Returns no error", t, func() {
		err := a.Do("test login", "test password")

		So(err, ShouldBeNil)
	})

	Convey("Returns login error", t, func() {
		err := a.Do("wrong login", "test password")

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "login incorrect")
	})

	Convey("Returns password error", t, func() {
		err := a.Do("test login", "wrong password")

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "password incorrect")
	})
}
