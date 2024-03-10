package model

import "net"

type Client struct {
	Id            uint16          /// clientID
	Type          int             /// clientType
	Status        int             /// clientStatus
	ServiceType   int             /// serverType
	ProcessJob    int             /// countOfTasksBeingProcessed
	RemoteAddr    string          /// remoteAddress
	SupportFormat []SupportFormat /// supportedConversionFormats
	PayLoad       float32         /// CPULoad
	Conn          net.Conn        /// clientConnection
	Server        *Server         /// belongingToTheServer
	SendQueue     chan []byte
	Marks         []string /// markupDataEnvironment  dev developmentEnvironment  release isAFormalEnvironment
	MaxParallel   int64    /// maximumNumberOfConcurrences
	Parallel      int64    /// numberOfConcurrences
	UpdateTime    int64
}

type SupportFormat struct {
	In  []string /// supportedConversionFormats
	Out []string /// conversionRules
}
