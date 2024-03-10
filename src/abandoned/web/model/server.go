package model

import "net"

var (
	EVENT_ACCEPT     uint16 = 1
	EVENT_CLIENT     uint16 = 2
	EVENT_CLOSE      uint16 = 3
	EVENT_CREATE_JOB uint16 = 4
	EVENT_CANCEL_JOB uint16 = 5
)

type Server struct {
	IsRunning   int32                           /// serverRunningFlag
	ListenSock  *net.TCPListener                /// listenToSocket
	EvtQueue    chan *ServerEvent               /// clientEventQueue
	JobQueue    []*Job                          /// taskQueue
	EvtQueueVIP chan *ServerEvent               /// clientEventQueue
	JobQueueVIP []*Job                          /// taskQueue
	Clients     map[uint16]*Client              /// registeredClient
	JobMgr      *JobManager                     /// taskManagement
	NewClient   func(*Server, net.Conn) *Client /// createANewClientInstance
}

type ServerEvent struct {
	Id   uint16
	Data interface{}
}

func (srv *Server) PostEventVIP(evtId uint16, event interface{}) {
	evt := new(ServerEvent)
	evt.Id = evtId
	evt.Data = event

	srv.EvtQueueVIP <- evt
}

func (srv *Server) PostEvent(evtId uint16, event interface{}) {
	evt := new(ServerEvent)
	evt.Id = evtId
	evt.Data = event

	srv.EvtQueue <- evt
}
