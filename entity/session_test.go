package entity

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSession(t *testing.T) {
	Convey("Returns new session instance", t, func() {
		s := NewSession()

		So(s, ShouldNotBeNil)

		Convey("Sessions subscribers is not nil", func() {
			So(s.Subscribers, ShouldNotBeNil)
		})
	})
}

func TestSession_AddParticipant(t *testing.T) {
	Convey("Add new participant to subscribers map", t, func() {
		s := NewSession()
		s.AddParticipant(&User{Name: "test user"})

		So(s.Subscribers["test user"], ShouldNotBeNil)

		Convey("with correct users data", func() {
			So(s.Subscribers["test user"].Name, ShouldEqual, "test user")
		})
	})
}

func TestSession_RemoveParticipant(t *testing.T) {
	Convey("Removes session participant", t, func() {
		user := &User{Name: "test user"}
		s := NewSession()
		s.AddParticipant(user)
		err := s.RemoveParticipant(user)

		So(err, ShouldBeNil)

		Convey("The session subscribers is empty", func() {
			So(s.Subscribers, ShouldBeEmpty)
		})

		Convey("Returns an error", func() {
			err := s.RemoveParticipant(&User{Name: "wrong user"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"subscriber: wrong user not found")
		})
	})
}
