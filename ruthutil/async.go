package ruthutil

import "errors"

func WaitChannel[T any](toUseChan chan T) (T, error) {
	for {

		select {
		case val, ok := <-toUseChan:
			if !ok {
				var empty T
				return empty, errors.New("channel closed")
			}
			return val, nil
		default:
			//this is a nil channel i'm not sure how
		}
	}
}

func IsChannelClosed[T any](toUseChan chan T) bool {
	select {
	case _, ok := <-toUseChan:
		if !ok {

			return true
		}
		return false
	default:
		//this is a nil channel i'm not sure how
	}
	return false
}
