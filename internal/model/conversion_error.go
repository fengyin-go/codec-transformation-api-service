package model

type ConversionFailure struct {
	Kind    string
	Message string
}

func (e *ConversionFailure) Error() string { return e.Message }

func NewRejectedConversion(message string) error {
	return &ConversionFailure{Kind: "rejected", Message: message}
}

func NewTemporaryConversion(message string) error {
	return &ConversionFailure{Kind: "temporary", Message: message}
}
