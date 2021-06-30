package core

func (i *InnerStruct) SendLogEvent() {
	for {
		select {
		case event := <-i.LogChannel:
			{
				i.RWMutex.RLock()
				for _, listener := range i.LogListeners {
					listener.log <- event
				}
				i.RWMutex.RUnlock()
			}
		case <-i.StopLogChannel:
			{
				i.RWMutex.RLock()
				for _, listener := range i.LogListeners {
					listener.stop <- struct{}{}
				}
				i.RWMutex.RUnlock()
				return
			}
		}
	}
}

func (i *InnerStruct) SendStatisticEvent() {
	for {
		select {
		case inf := <-i.StatisticChannel:
			{
				i.RWMutex.RLock()
				for _, listener := range i.StatisticListeners {
					listener.info <- inf
				}
				i.RWMutex.RUnlock()
			}
		case <-i.StopStatisticChannel:
			{
				i.RWMutex.RLock()
				for _, listener := range i.StatisticListeners {
					listener.stop <- struct{}{}
				}
				i.RWMutex.RUnlock()
				return
			}
		}
	}
}
