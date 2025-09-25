package uid

import (
	"github.com/sony/sonyflake/v2"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	sonyInstance *sonyflake.Sonyflake
	once         sync.Once
)

func Init() {
	once.Do(func() {
		startTimeStr := os.Getenv("SONYFLAKE_START_TIME")
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			panic("Invalid SONYFLAKE_START_TIME format")
		}

		machineIDStr := os.Getenv("SONYFLAKE_MACHINE_ID")
		machineID, err := strconv.Atoi(machineIDStr)
		if err != nil {
			panic("Invalid SONYFLAKE_MACHINE_ID format")
		}

		sony, err := sonyflake.New(sonyflake.Settings{
			StartTime: startTime,
			MachineID: func() (int, error) {
				return machineID, nil
			},
		})
		if err != nil {
			panic(err)
		}
		sonyInstance = sony
	})
}

func Gen() (int64, error) {
	return sonyInstance.NextID()
}
