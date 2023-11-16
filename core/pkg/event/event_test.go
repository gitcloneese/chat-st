package event

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type systemA struct {
}

func (s *systemA) Action1(num int, name string) {
	fmt.Printf("systemA Action1 num:%d name:%s\n", num, name)
}

func (s *systemA) Action2(num int) {
	fmt.Printf("systemA Action2 num:%d\n", num)
}

func (s *systemA) action3(num int) {
	fmt.Printf("systemA action3 num:%d\n", num)
}

type systemB struct {
}

func (s *systemB) Action1(num int, name string) {
	fmt.Printf("systemB Action1 num:%d name:%s\n", num, name)
}

func (s *systemB) Action2(num int) {
	fmt.Printf("systemB Action2 num:%d\n", num)
}

const (
	EVENT_TYPE_ACTION1 = "EVENT_TYPE_ACTION1"
	EVENT_TYPE_ACTION2 = "EVENT_TYPE_ACTION2"
)

func Test_EventDispatcher(t *testing.T) {

	Convey("Test_EventDispatcher", t, func() {
		var (
			err error
		)

		dispatcher := NewEventDispatcher()

		sysA := &systemA{}
		sysB := &systemB{}

		err = dispatcher.AddEventListener(EVENT_TYPE_ACTION1, sysA.Action1)
		So(err, ShouldBeNil)
		err = dispatcher.AddEventListener(EVENT_TYPE_ACTION2, sysA.Action2)
		So(err, ShouldBeNil)
		err = dispatcher.AddEventListener(EVENT_TYPE_ACTION2, sysA.Action2)
		So(err, ShouldNotBeNil)
		// err = dispatcher.AddEventListener(EVENT_TYPE_ACTION2, sysA.action3)
		// So(err, ShouldNotBeNil)

		err = dispatcher.AddEventListener(EVENT_TYPE_ACTION1, sysB.Action1)
		So(err, ShouldBeNil)
		err = dispatcher.AddEventListener(EVENT_TYPE_ACTION2, sysB.Action2)
		So(err, ShouldBeNil)

		dispatcher.DispatchEvent(EVENT_TYPE_ACTION1, 100, "hello")
		dispatcher.DispatchEvent(EVENT_TYPE_ACTION2, 200)
		dispatcher.DispatchEvent(EVENT_TYPE_ACTION2, "hello")
	})
}
