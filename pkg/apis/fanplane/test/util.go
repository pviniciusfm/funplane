package test

type ConvertionResult int
const (
	Fail ConvertionResult = iota
	Success
)

// FixtureTable defines an common structure to test model serialization
type FixtureTable struct {
	Title                 string
	InputFile             string
	ExpectedOutputMessage string
	ConversionResult      ConvertionResult
}
