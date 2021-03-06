package pewpew

import (
	"fmt"
	"time"
)

type requestStatSummary struct {
	avgRPS               float64 //requests per nanoseconds
	avgDuration          time.Duration
	maxDuration          time.Duration
	minDuration          time.Duration
	statusCodes          map[int]int //counts of each code
	startTime            time.Time   //start of first request
	endTime              time.Time   //end of last request
	avgDataTransferred   int         //bytes
	maxDataTransferred   int         //bytes
	minDataTransferred   int         //bytes
	totalDataTransferred int         //bytes
}

//create statistical summary of all requests
func CreateRequestsStats(requestStats []RequestStat) requestStatSummary {
	if len(requestStats) == 0 {
		return requestStatSummary{}
	}

	requestCodes := make(map[int]int)
	summary := requestStatSummary{maxDuration: requestStats[0].Duration,
		minDuration:          requestStats[0].Duration,
		minDataTransferred:   requestStats[0].DataTransferred,
		statusCodes:          requestCodes,
		startTime:            requestStats[0].StartTime,
		endTime:              requestStats[0].EndTime,
		totalDataTransferred: 0,
	}
	var totalDurations time.Duration //total time of all requests (concurrent is counted)
	nonErrCount := 0
	for i := 0; i < len(requestStats); i++ {
		if requestStats[i].Duration > summary.maxDuration {
			summary.maxDuration = requestStats[i].Duration
		}
		if requestStats[i].Duration < summary.minDuration {
			summary.minDuration = requestStats[i].Duration
		}
		if requestStats[i].StartTime.Before(summary.startTime) {
			summary.startTime = requestStats[i].StartTime
		}
		if requestStats[i].EndTime.After(summary.endTime) {
			summary.endTime = requestStats[i].EndTime
		}

		if requestStats[i].DataTransferred > summary.maxDataTransferred {
			summary.maxDataTransferred = requestStats[i].DataTransferred
		}
		if requestStats[i].DataTransferred < summary.minDataTransferred {
			summary.minDataTransferred = requestStats[i].DataTransferred
		}

		totalDurations += requestStats[i].Duration
		summary.statusCodes[requestStats[i].StatusCode]++
		summary.totalDataTransferred += requestStats[i].DataTransferred
		if requestStats[i].Error == nil {
			nonErrCount++
		}
	}
	//kinda ugly to calculate average, then convert into nanoseconds
	if nonErrCount == 0 {
		summary.avgDuration = 0
	} else {
		avgNs := totalDurations.Nanoseconds() / int64(nonErrCount)
		newAvg, _ := time.ParseDuration(fmt.Sprintf("%d", avgNs) + "ns")
		summary.avgDuration = newAvg
	}

	if nonErrCount == 0 {
		summary.avgDataTransferred = 0
	} else {
		summary.avgDataTransferred = summary.totalDataTransferred / nonErrCount
	}

	summary.avgRPS = float64(len(requestStats)) / float64(summary.endTime.Sub(summary.startTime))
	return summary
}
