package queue

import (
	"math"
	"sqt/config"
	"sqt/message"
	"sqt/queuetasks"
	"time"
)

func Run(keyToRead string, taskNum int, queueChannel chan message.Message) {
	if taskNum > config.Values.MaxStackSize {
		result := message.Message{
			IsExecuted:    false,
			Status:        message.STATUS_MAX_QUEUE_EXCEEDED,
			QueueSize:     taskNum,
		};

		queueChannel <- result;
	} 
	
	start := time.Now()
	var minTimeToExecute int

	switch config.Values.ReadTimeGrowth {
	case "sum":
		minTimeToExecute = config.Values.ReadTimeInit + taskNum * config.Values.ReadTimeStep

	case "msum":
		minTimeToExecute = config.Values.ReadTimeInit + int(math.Round(float64(taskNum) *
			float64(config.Values.ReadTimeStep) * config.Values.ReadTimeParameter1))

	case "exp":
		minTimeToExecute = config.Values.ReadTimeInit +
			int(math.Pow(float64(config.Values.ReadTimeStep), float64(taskNum)))

	case "log":
		minTimeToExecute = config.Values.ReadTimeInit + int(float64(config.Values.ReadTimeStep) * math.Log(float64(taskNum + 1)))

	default:
		result := message.Message{
			IsExecuted:    false,
			Status:        message.STATUS_WRONG_CONFIG,
			QueueSize:     taskNum,
		};

		queueChannel <- result;
	}

	readedData, messageStatus := queuetasks.GetData(keyToRead)

	elapsed := int(time.Since(start).Milliseconds())

	if elapsed < minTimeToExecute {
		timeToSleep := minTimeToExecute - elapsed
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	}

	result := message.Message{
		IsExecuted:    true,
		Status:        messageStatus,
		Data:          readedData,
		TimeElapsed:   elapsed,
		TimeQueuedMin: minTimeToExecute,
		QueueSize:     taskNum,
	}

	queueChannel <- result

}

func task(msToSleep int) {
	dMsToSleep := time.Duration(msToSleep)
	time.Sleep(dMsToSleep * time.Millisecond)
}