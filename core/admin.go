package core

import (
	"grpc_micro/contract"
	"time"
)

func (i *InnerStruct) Logging(nothing *contract.Nothing, logServer contract.Admin_LoggingServer) error {
	listener := &ListenerLog{
		log:  make(chan *contract.Event, 0),
		stop: make(chan interface{}, 0),
	}

	i.RWMutex.Lock()
	i.LogListeners = append(i.LogListeners, listener)
	i.RWMutex.Unlock()

	for {
		select {
		case log := <-listener.log:
			{
				err := logServer.SendMsg(log)
				if err != nil {
					return err
				}
			}
		case <-listener.stop:
			{
				return nil
			}
		}
	}
}

func (i *InnerStruct) Statistics(s *contract.StatInterval, ass contract.Admin_StatisticsServer) error {
	newListener := newStatisticListener()

	i.RWMutex.Lock()
	i.StatisticListeners = append(i.StatisticListeners, newListener)
	i.RWMutex.Unlock()

	ticker := time.NewTicker(time.Second * time.Duration(s.IntervalSeconds))

	for {
		select {
		case inf := <-newListener.info:
			{
				newListener.accumulator.accept(inf.method, inf.consumer)
			}
		case <-ticker.C:
			{
				err := ass.Send(&contract.Stat{
					Timestamp:  time.Now().Unix(),
					ByMethod:   newListener.accumulator.methodInvoke,
					ByConsumer: newListener.accumulator.userInvoke,
				})
				if err != nil {
					return err
				}
				newListener.accumulator.drop()
			}
		case <-newListener.stop:
			{
				return nil
			}
		}
	}
}
