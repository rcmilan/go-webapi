package book

type ErrConflict struct{ Msg string }

func (e ErrConflict) Error() string { return e.Msg }
