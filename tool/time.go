package tool

import "time"

/*
	1. Current Local Time: 2023-10-13 10:08:23.123456789 +0800 CST m=+0.000123456
	2. Current UTC Time: 2023-10-13 02:08:23.123456789 +0000 UTC
	3. Formatted Time: 2023-10-13 10:08:23
	4. Current DateClock:
		date:	2023-10-13
		clock:	10:08:23
	5. Unix Time (seconds):
		s		second
		ms		milli second
		us		micro second
		ns		nano second
*/

// 1. 获取当前的本地时间
func CurrentLocalTime() time.Time {
	return time.Now()
}

// 2. 获取当前的UTC时间
func CurrentUTCTime() time.Time {
	return time.Now().UTC()
}

// 3. 获取当前的本地时间(格式化)
func CurrentLocalTimeFormatted() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 4. 获取当前时间的各个组成部分
// Date: 年，月，日
// Clock: 时，分，秒
func CurrentDateClock() (year int, month time.Month, day, hour, minute, second int) {
	currentTime := time.Now()
	year, month, day = currentTime.Date()
	hour, minute, second = currentTime.Clock()

	return
}

// 5. 获取当前unix时间戳
/*
	Accuracy:
		s		second
		ms		milli second
		us		micro second
		ns		nano second
*/
func CurrentUnixTime(accuracy string) int64 {
	currentTime := time.Now()

	switch accuracy {
	case "s":
		return currentTime.Unix()
	case "ms":
		return currentTime.UnixMilli()
	case "us":
		return currentTime.UnixMicro()
	case "ns":
		return currentTime.UnixNano()
	default:
		return currentTime.Unix()
	}
}
