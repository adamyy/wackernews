package view

import (
	"fmt"
	"sync"

	"github.com/adamyy/wackernews/view/text"
	"github.com/jroimartin/gocui"
)

type MessageView struct {
	*Prop

	message string

	tainted bool
	mutex   sync.Mutex
}

func NewMessageView(opts ...PropOption) *MessageView {
	v := &MessageView{Prop: DefaultProp()}
	_ = v.Set(opts...)
	return v
}

func (mv *MessageView) SetMessage(message string) {
	mv.mutex.Lock()
	defer mv.mutex.Unlock()

	mv.message = message
	mv.tainted = true
}

func (mv *MessageView) KeyBindings() KeyBindings {
	return KeyBindings{} // no key bindings
}

func (mv *MessageView) Draw(v *gocui.View) error {
	mv.mutex.Lock()
	defer mv.mutex.Unlock()

	if mv.message == "" { // content is missing, skipping rendering
		return nil
	}

	if !mv.tainted { // no need re-rendered
		return nil
	}
	defer func() { mv.tainted = false }()

	v.Clear()
	v.Frame = true
	for _, line := range mv.Render() {
		if _, err := fmt.Fprintln(v, line); err != nil {
			return err
		}
	}

	return nil
}

func (mv *MessageView) Name() string {
	return mv.name
}

func (mv *MessageView) Render() []string {
	width, _ := mv.Size()
	return text.Justify(mv.message, width, true)
}
