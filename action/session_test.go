package action

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flexconstructor/openvidu-tutorial/entity"
	"github.com/flexconstructor/openvidu-tutorial/repository"
)

func TestSession_Add(t *testing.T) {
	Convey("Adds new session to repository", t, func() {
		a := Session{
			SessionRepo: repository.NewSessionsRepository(),
			UserRepo:    repository.NewUsersRepository(),
		}
		a.UserRepo.Add("test user", "test password", 1)
		a.UserRepo.Add("test participant", "test password", 0)
		err := a.Add("test session id", "test session name", "test user")
		So(err, ShouldBeNil)

		Convey("New Session was be added to repository", func() {
			s, _ := a.SessionRepo.Get("test session name")
			So(s, ShouldNotBeNil)
			So(s.Name, ShouldEqual, "test session name")
			So(s.ID, ShouldEqual, "test session id")
			So(s.Owner, ShouldNotBeNil)
			So(s.Owner.Name, ShouldEqual, "test user")
		})

		Convey("Add user to subscribers", func() {
			err := a.Add("test session id", "test session name",
				"test participant")
			So(err, ShouldBeNil)
			s, _ := a.SessionRepo.Get("test session name")
			So(s.Subscribers, ShouldNotBeEmpty)
			So(s.Subscribers["test participant"].Name, ShouldEqual,
				"test participant")
		})

		Convey("Returns an user error", func() {
			err := a.Add("test session id", "test session name", "wrong user")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "login incorrect")
		})

	})
}

func TestSession_Delete(t *testing.T) {
	Convey("Deletes session from repository", t, func() {
		a := Session{
			SessionRepo: repository.NewSessionsRepository(),
			UserRepo:    repository.NewUsersRepository(),
		}
		a.UserRepo.Add("test user", "test password", 1)
		a.UserRepo.Add("test participant", "test password", 0)
		a.SessionRepo.Add("test session id", "test session name",
			&entity.User{Name: "test user"})

		err := a.Delete("test session name", "test user")
		So(err, ShouldBeNil)

		Convey("test session was be removed from repository", func() {
			_, err := a.SessionRepo.Get("test session name")
			So(err, ShouldNotBeNil)
		})

		Convey("Remove participant from session", func() {
			s, _ := a.SessionRepo.Add("test session id", "test session name",
				&entity.User{Name: "test user"})
			s.AddParticipant(&entity.User{Name: "test participant"})
			err := a.Delete("test session name", "test participant")
			So(err, ShouldBeNil)

			Convey("Session subscribers should be empty", func() {
				So(s.Subscribers, ShouldBeEmpty)
			})
		})

		Convey("Returns a session error", func() {
			err := a.Delete("wrong session name", "test user")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"session wrong session name does not exists")
		})

		Convey("Returns a user error", func() {
			err := a.Delete("test session name", "wrong user")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "login incorrect")
		})
	})
}

func TestSession_GetID(t *testing.T) {
	Convey("Returns session ID", t, func() {
		a := Session{
			SessionRepo: repository.NewSessionsRepository(),
			UserRepo:    repository.NewUsersRepository(),
		}
		a.UserRepo.Add("test user", "test password", 1)
		a.Add("test session id", "test session name", "test user")

		sID, err := a.GetID("test session name")
		Convey("Returns no errors", func() {
			So(err, ShouldBeNil)
		})

		Convey("Session ID equal created session ID", func() {
			So(sID, ShouldEqual, "test session id")
		})

		Convey("Returns an error", func() {
			_, err := a.GetID("wrong session name")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring,
				"session wrong session name does not exists")
		})
	})
}

func TestSession_addParticipant(t *testing.T) {
	a := Session{
		SessionRepo: repository.NewSessionsRepository(),
		UserRepo:    repository.NewUsersRepository(),
	}
	a.UserRepo.Add("test user", "test password", 1)
	a.Add("test session id", "test session name", "test user")

	Convey("Returns user error", t, func() {
		_, err := a.addParticipant("test session id", "wrong user name")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "login incorrect")
	})

	Convey("Returns session error", t, func() {
		_, err := a.addParticipant("wrong session id", "test user")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring,
			"session wrong session id does not exists")
	})

	Convey("Session owner cannot be added as subscriber", t, func() {
		_, err := a.addParticipant("test session name", "test user")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring,
			"owner test user can not subscribe session test session name")
	})

	Convey("Duplicate of participant is impossible", t, func() {
		a.UserRepo.Add("test participant", "test password", 0)
		a.Add("test session id", "test session name", "test participant")
		err := a.Add("test session id", "test session name", "test participant")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring,
			"user test participant already subscribed to the session test session name")

	})
}
