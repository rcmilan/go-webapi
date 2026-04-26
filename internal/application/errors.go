package application

type ErrValidation struct{ Msg string }

func (e ErrValidation) Error() string { return e.Msg }

type ErrNotFound struct{ Msg string }

func (e ErrNotFound) Error() string { return e.Msg }

