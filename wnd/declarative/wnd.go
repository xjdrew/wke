package declarative

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/xjdrew/wke/wnd"
)

type WkeWnd struct {
	// Window

	Background       Brush
	ContextMenuItems []MenuItem
	Enabled          Property
	Font             Font
	MaxSize          Size
	MinSize          Size
	Name             string
	OnKeyDown        walk.KeyEventHandler
	OnKeyPress       walk.KeyEventHandler
	OnKeyUp          walk.KeyEventHandler
	OnMouseDown      walk.MouseEventHandler
	OnMouseMove      walk.MouseEventHandler
	OnMouseUp        walk.MouseEventHandler
	OnSizeChanged    walk.EventHandler
	Persistent       bool
	ToolTipText      Property
	Visible          Property

	// Widget

	AlwaysConsumeSpace bool
	Column             int
	ColumnSpan         int
	Row                int
	RowSpan            int
	StretchFactor      int

	// WebView

	AssignTo       **wnd.WkeWnd
	URL            Property
	Title          Property
	OnURLChanged   walk.EventHandler
	OnTitleChanged walk.EventHandler
}

func (ww WkeWnd) Create(builder *Builder) error {
	w, err := wnd.NewWkeWnd(builder.Parent())
	if err != nil {
		return err
	}

	return builder.InitWidget(ww, w, func() error {
		if ww.OnURLChanged != nil {
			w.URLChanged().Attach(ww.OnURLChanged)
		}

		if ww.OnTitleChanged != nil {
			w.TitleChanged().Attach(ww.OnTitleChanged)
		}
		if ww.AssignTo != nil {
			*ww.AssignTo = w
		}
		return nil
	})
}
