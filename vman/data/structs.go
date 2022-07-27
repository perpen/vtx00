package data

type ApiRequest struct {
	Cmd AbstractCmd
	Out chan interface{}
}

type AbstractCmd struct {
	GetStateCmd    *GetStateCmd
	CreatePanelCmd *CreatePanelCmd
	KillCmd        *KillCmd
	LayoutCmd      *LayoutCmd
}

type GetStateCmd struct {
}

type CreatePanelCmd struct {
	Argv []string
	Cwd  string
	Meta string
}

type KillCmd struct {
	PanelID int
}

type LayoutCmd struct {
	FocusID int
	Panels  []Panel
}

type State struct {
	FocusID int
	Size    []int
	Panels  []Panel
}

type Panel struct {
	ID      int
	Pos     []int
	Border  Border
	Content string
	Z       int
	Meta    string
}

type Border struct {
	Style      Style
	Title      BarComp
	Components []BarComp
}

type Style struct {
	Fg, Bg string
	Bold   bool
}

type BarComp struct {
	String string
	Style  Style
}

type Event struct {
	Type    string
	Target  int
	Details interface{}
	Time    string
}
