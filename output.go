package main

import (
	"fmt"
	"strings"
	"time"
)

// PingResultToJson convert PingSatistic to  RingResult Json format for API
func PingResultToJson(stat *PingSatistic) DataPoint1MinResultJSON {
	_, mean, _, _, _ := stat.TimeMeasure.Statistic()
	var errorShow string
	loss := fmt.Sprintf("%3.1f%s", stat.Loss*100, "%")
	if stat.Count == 0 {
		errorShow = "No Data"
		loss = fmt.Sprintf("%3.1f%s", float64(0), "%")
	}
	jsonResult := DataPoint1MinResultJSON{Submitted: int(stat.Submitted), Confirmed: int(stat.Confirmed), Mean: int(mean), Error: errorShow}

	jsonResult.Loss = loss
	ts := time.Unix(stat.TimeStamp, 0).UTC()
	jsonResult.TimeStamp = ts.Format(time.RFC3339)
	return jsonResult
}

// ReportPayload get the report within specified minutes
func (s *SlackPayload) ReportPayload(c Cluster, data *GroupsAllStatistic) {
	// Header Block
	headerText := fmt.Sprintf("total-submitted: %3.0f, total-confirmed:%3.0f, average-loss:%3.1f%s",
		data.GetGroupsAllStatistic(false).Submitted,
		data.GetGroupsAllStatistic(false).Confirmed,
		data.GetGroupsAllStatistic(false).Loss*100, "%")
	header := Block{
		BlockType: "section",
		BlockText: SlackText{
			SType: "mrkdwn",
			SText: headerText,
		},
	}
	s.Blocks = append(s.Blocks, header)
	// BodyBlock
	body := Block{}
	records := reportRecordBlock(data)
	description := "( Submitted, Confirmed, Loss, min/mean/max/stddev ms )"
	memo := "* rpc error : context deadline exceeded does not count as trascation\n** BlockhashNotFound does not show in Error List"
	errorRecords := reportErrorBlock(data)
	body = Block{
		BlockType: "section",
		BlockText: SlackText{SType: "mrkdwn", SText: fmt.Sprintf("```%s\n%s\n%s\n%s```", description, records, memo, errorRecords)},
	}
	s.Blocks = append(s.Blocks, body)
}

func (s *SlackPayload) AlertPayload(c Cluster, gStat *GlobalStatistic) {
	var text, timeStatis string
	if gStat.TimeStatistic.Stddev <= 0 {
		timeStatis = fmt.Sprintf(" %d/%3.0f/%d/%s ", gStat.TimeStatistic.Min, gStat.TimeStatistic.Mean, gStat.TimeStatistic.Max, "NaN")
	} else {
		timeStatis = fmt.Sprintf(" %d/%3.0f/%d/%3.0f ", gStat.TimeStatistic.Min, gStat.TimeStatistic.Mean, gStat.TimeStatistic.Max, gStat.TimeStatistic.Stddev)
	}
	text = fmt.Sprintf("{ hostname: %s, submitted: %3.0f,confirmed:%3.0f, loss: %3.1f%s, confirmation:%s}",
		config.HostName, gStat.Submitted, gStat.Confirmed, gStat.Loss*100, "%", timeStatis)

	header := Block{
		BlockType: "section",
		BlockText: SlackText{
			SType: "mrkdwn",
			SText: text,
		},
	}
	s.Blocks = append(s.Blocks, header)
}

func reportRecordBlock(data *GroupsAllStatistic) string {
	text := ""
	timeStatis := ""
	for _, ps := range data.PingStatisticList {
		if ps.TimeStatistic.Stddev <= 0 {
			timeStatis = fmt.Sprintf(" %d/%3.0f/%d/%s ", ps.TimeStatistic.Min, ps.TimeStatistic.Mean, ps.TimeStatistic.Max, "NaN")
		} else {
			timeStatis = fmt.Sprintf(" %d/%3.0f/%d/%3.0f ", ps.TimeStatistic.Min, ps.TimeStatistic.Mean, ps.TimeStatistic.Max, ps.TimeStatistic.Stddev)
		}
		lossPercentage := ps.Loss * 100
		if ps.Count > 0 {
			text = fmt.Sprintf("%s( %3.0f, %3.0f, %3.1f%s, %s )\n", text, ps.Submitted, ps.Confirmed, lossPercentage, "%", timeStatis)
		}
	}
	return text
}

func reportErrorBlock(data *GroupsAllStatistic) string {
	var exceededText, errorText string
	if len(data.GlobalErrorStatistic) == 0 {
		return ""
	}
	for k, v := range data.GlobalErrorStatistic {
		if strings.Contains(k, string(RPCServerDeadlineExceededKey)) {
			exceededText = fmt.Sprintf("*(count:%d) RPC Server context deadline exceed\n", v)
		} else if !strings.Contains(k, string(BlockhashNotFoundKey)) {
			errorText = fmt.Sprintf("%s\n(count: %d) %s", errorText, v, k)
		}
	}
	if len(exceededText) > 0 || len(errorText) > 0 {
		return fmt.Sprintf("Error List:%s\n%s", exceededText, errorText)
	}
	return ""
}

func reportRawErrorBlock(data *GroupsAllStatistic) string {
	var errorText string
	if len(data.GlobalErrorStatistic) == 0 {
		return ""
	}
	for k, v := range data.GlobalErrorStatistic {
		errorText = fmt.Sprintf("%s\n(count: %d) %s", errorText, v, k)
	}
	return fmt.Sprintf("Error List:%s", errorText)
}
