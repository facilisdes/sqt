package queue

import (
	"errors"
	"math"
	"sqt/config"
	"sqt/message"
	"sqt/queuetasks"
	"time"
)

const (
	MODE_QUEUE       = 0
	MODE_HEALTHCHECK = 1
)

func Run(keyToRead string, taskNum int, queueChannel chan message.Message, mode int) {
	if taskNum > config.Values.MaxStackSize {
		// refuse to execute command and return immediately if max queue size is exceeded
		result := message.Message{
			IsExecuted: false,
			Status:     message.STATUS_MAX_QUEUE_EXCEEDED,
			QueueSize:  taskNum,
		}

		queueChannel <- result
	}

	start := time.Now()

	var minTimeToExecute int
	var err error

	if mode == MODE_QUEUE {
		minTimeToExecute, err = getExecutionTime(taskNum)

		if err != nil {
			result := message.Message{
				IsExecuted: false,
				Status:     message.STATUS_WRONG_CONFIG,
				QueueSize:  taskNum,
			}
			queueChannel <- result
		}
	} else if mode == MODE_HEALTHCHECK {
		minTimeToExecute = 0
	} else {
		result := message.Message{
			IsExecuted: false,
			Status:     message.STATUS_WRONG_COMMAND_TYPE,
			QueueSize:  taskNum,
		}
		queueChannel <- result
	}

	dataRead, messageStatus := queuetasks.GetData(keyToRead)

	elapsed := int(time.Since(start).Milliseconds())

	if elapsed < minTimeToExecute {
		timeToSleep := minTimeToExecute - elapsed
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	}

	elapsedTotal := int(time.Since(start).Milliseconds())

	result := message.Message{
		IsExecuted:       true,
		Status:           messageStatus,
		Data:             dataRead,
		TimeElapsed:      elapsed,
		TimeQueuedMin:    minTimeToExecute,
		TimeElapsedTotal: elapsedTotal,
		QueueSize:        taskNum,
	}

	queueChannel <- result
}

func getExecutionTime(taskNum int) (int, error) {
	var minTimeToExecute int
	switch config.Values.ReadTimeGrowth {
	case "sum":
		// execution time is calculated as init time + read time step for every command in queue
		// linear fixed growth of read time
		minTimeToExecute = config.Values.ReadTimeInit + taskNum*config.Values.ReadTimeStep

	case "msum":
		// execution time is calculated as init time + read time step for every command in queue
		// impact of queue size is controlled by a first read parameter (ReadTimeParameter)
		// linear managed growth of read time
		minTimeToExecute = config.Values.ReadTimeInit + int(math.Round(float64(taskNum)*
			float64(config.Values.ReadTimeStep)*config.Values.ReadTimeParameter))

	case "exp":
		// execution time is calculated as init time + step to the power of count of commands in queue
		// fast exponential fixed growth of read time
		minTimeToExecute = config.Values.ReadTimeInit +
			int(math.Pow(float64(config.Values.ReadTimeStep), float64(taskNum)))

	case "mexp":
		// execution time is calculated as init time + step to the power of count of commands in queue
		// impact of queue size is controlled by a first read parameter (ReadTimeParameter)
		// fast exponential managed growth of read time
		minTimeToExecute = config.Values.ReadTimeInit +
			int(math.Pow(float64(config.Values.ReadTimeStep), float64(taskNum)*config.Values.ReadTimeParameter))

	case "log":
		// execution time is calculated as init time + step to the log of count of commands in queue
		// slow logarithmic fixed growth of read time
		minTimeToExecute = config.Values.ReadTimeInit + int(float64(config.Values.ReadTimeStep)*math.Log(float64(taskNum+1)))

	case "mlog":
		// execution time is calculated as init time + step to the log of count of commands in queue
		// impact of queue size is controlled by a first read parameter (ReadTimeParameter)
		// slow logarithmic managed growth of read time
		minTimeToExecute = config.Values.ReadTimeInit + int(float64(config.Values.ReadTimeStep)*math.Log(float64(taskNum+1)*config.Values.ReadTimeParameter))

	default:
		return 0, errors.New("Wrong config - no supported growth function supplied")
	}

	return minTimeToExecute, nil
}
