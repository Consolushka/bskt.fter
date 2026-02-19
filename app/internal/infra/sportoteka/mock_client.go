package sportoteka

import (
	"reflect"
	"time"

	"github.com/golang/mock/gomock"
)

type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

func (m *MockClientInterface) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BoxScore", gameId)
	ret0, _ := ret[0].(GameBoxScoreResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockClientInterfaceMockRecorder) BoxScore(gameId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BoxScore", reflect.TypeOf((*MockClientInterface)(nil).BoxScore), gameId)
}

func (m *MockClientInterface) Calendar(tag string, season int, from time.Time, to time.Time) (CalendarResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Calendar", tag, season, from, to)
	ret0, _ := ret[0].(CalendarResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockClientInterfaceMockRecorder) Calendar(tag, season, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calendar", reflect.TypeOf((*MockClientInterface)(nil).Calendar), tag, season, from, to)
}
