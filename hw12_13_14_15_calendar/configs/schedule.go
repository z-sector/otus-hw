package configs

import "fmt"

type ScheduleConf struct {
	DeleteEventsAgeDay   int `validate:"required,gt=0"`
	DeleteEventsCronbeat int `validate:"required,gt=0"`
	SendNotifyCronbeat   int `validate:"required,gt=0"`
}

func (s ScheduleConf) String() string {
	return fmt.Sprintf(
		"{DeleteEvents(days=%d, beat=%d), SendNotify(beat=%d}",
		s.DeleteEventsAgeDay,
		s.DeleteEventsCronbeat,
		s.SendNotifyCronbeat,
	)
}
