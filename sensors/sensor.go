package sensors

import "hostmonitor/measure"

type sensor interface {
	Poll(*measure.Measure)
}
