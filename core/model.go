package core

import (
	"grpc_micro/contract"
	"sync"
)

type InnerStruct struct {
	sync.RWMutex
	Acl map[string][]string

	LogChannel     chan *contract.Event
	StopLogChannel chan interface{}
	LogListeners   []*ListenerLog

	StatisticChannel     chan *Info
	StopStatisticChannel chan interface{}
	StatisticListeners   []*StatisticListener
}

type Info struct {
	method   string
	consumer string
}

type ListenerLog struct {
	log  chan *contract.Event
	stop chan interface{}
}

type StatisticListener struct {
	info        chan *Info
	stop        chan interface{}
	accumulator *accumulator
}

func newStatisticListener() *StatisticListener {
	return &StatisticListener{
		info: make(chan *Info, 0),
		stop: make(chan interface{}, 0),
		accumulator: &accumulator{
			methodInvoke: make(map[string]uint64, 0),
			userInvoke:   make(map[string]uint64, 0),
		},
	}
}

type accumulator struct {
	methodInvoke map[string]uint64
	userInvoke   map[string]uint64
}

func (s *accumulator) accept(nameMethod, nameUser string) {
	incCounterInMap(s.methodInvoke, nameMethod)
	incCounterInMap(s.userInvoke, nameUser)
}

func incCounterInMap(counterMap map[string]uint64, name string) {
	if v, ok := counterMap[name]; ok {
		counterMap[name] = v + 1
		return
	}
	counterMap[name] = 1
}

func (s *accumulator) drop() {
	s.userInvoke = make(map[string]uint64, 0)
	s.methodInvoke = make(map[string]uint64, 0)
}
