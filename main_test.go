package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_main(t *testing.T) {
	Convey("Test convey", t, func() {
		So(2, ShouldEqual, 1+1)
	})
}
