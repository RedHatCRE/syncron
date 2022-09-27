package mocks

import "github.com/stretchr/testify/mock"

type CLIParserMock struct {
	mock.Mock
}

func (mock *CLIParserMock) GetDebug() bool {
	args := mock.Called()

	return args.Bool(0)
}
