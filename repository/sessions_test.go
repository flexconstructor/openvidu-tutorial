package repository

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

func TestNewSessionsRepository(t *testing.T) {
	Convey("Returns new instance of repository", t, func() {
		r := NewSessionsRepository()

		So(r, ShouldNotBeNil)

		Convey("The repository storage is not nil", func() {
			So(r.storage, ShouldNotBeNil)
		})
	})
}

func TestSessions_Add(t *testing.T) {
	Convey("Add new session to repository", t, func() {
		r := NewSessionsRepository()
		s, err := r.Add("test session ID", "test session name",
			&entity.User{Name: "test user", Password: "test password", Role: 1})

		Convey("Returns not nil session", func() {
			So(s, ShouldNotBeNil)
		})

		Convey("Returns no errors", func() {
			So(err, ShouldBeNil)
		})

		Convey("The repository storage has not nil session", func() {
			So(r.storage["test session name"], ShouldNotBeNil)
		})

		Convey("Returned session equal stored session", func() {
			So(s, ShouldEqual, r.storage["test session name"])
		})

		Convey("The session has correct session parameters", func() {
			So(s.Name, ShouldEqual, "test session name")
			So(s.ID, ShouldEqual, "test session ID")
			So(s.Subscribers, ShouldNotBeNil)
			So(s.Owner, ShouldNotBeNil)
		})

		Convey("Session owner has correct user parameters", func() {
			So(r.storage["test session name"].Owner.Name,
				ShouldEqual, "test user")
			So(r.storage["test session name"].Owner.Password,
				ShouldEqual, "test password")
			So(r.storage["test session name"].Owner.Role, ShouldEqual, 1)
		})

		Convey("Returns an error", func() {
			_, err := r.Add("test session ID", "test session name",
				&entity.User{Name: "test user",
					Password: "test password", Role: 1})

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"session test session name already exists")
		})
	})
}

func TestSessions_Delete(t *testing.T) {
	Convey("Deletes session from repository", t, func() {
		r := NewSessionsRepository()
		s, _ := r.Add("test session ID", "test session name",
			&entity.User{Name: "test user", Password: "test password", Role: 1})

		Convey("Returns no error", func() {
			err := r.Delete(s.Name)
			So(err, ShouldBeNil)
		})

		Convey("Returns an error", func() {
			err := r.Delete("wrong session name")
			So(err, ShouldNotBeNil)
			So(err.Error(),
				ShouldContainSubstring,
				"session wrong session name does not exists")
		})
	})
}

func TestSessions_Get(t *testing.T) {
	Convey("Returns a session", t, func() {
		r := NewSessionsRepository()
		s, _ := r.Add("test session ID", "test session name",
			&entity.User{Name: "test user", Password: "test password", Role: 1})
		session, err := r.Get(s.Name)

		Convey("Returns no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("Returns not nil session", func() {
			So(session, ShouldNotBeNil)
		})

		Convey("Returns session that was be added", func() {
			So(session, ShouldEqual, s)
		})

		Convey("Returns an error", func() {
			_, err := r.Get("wrong name")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"session wrong name does not exists")
		})
	})
}

func TestSessions_Leave(t *testing.T) {
	Convey("Removes session participant", t, func() {
		r := NewSessionsRepository()
		s, _ := r.Add("test session ID", "test session name",
			&entity.User{Name: "test user", Password: "test password", Role: 1})
		s.AddParticipant(&entity.User{Name: "test participant"})
		err := r.Leave("test session name", "test participant")

		So(err, ShouldBeNil)

		Convey("Session subscribers has no any participants", func() {
			So(s.Subscribers, ShouldBeEmpty)
		})

		Convey("Returns a session error", func() {
			err := r.Leave("wrong session name", "test participant")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"session wrong session name does not exists")
		})

		Convey("Returns participant error", func() {
			err := r.Leave("test session name", "wrong participant")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"user wrong participant does not exists")
		})
	})
}
