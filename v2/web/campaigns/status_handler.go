package campaigns

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/go-sql-driver/mysql"
	"github.com/ryanmoran/stack"
)

type campaignStatusGetter interface {
	Get(connection collections.ConnectionInterface, campaignID string) (collections.CampaignStatus, error)
}

type StatusHandler struct {
	campaignStatuses campaignStatusGetter
}

func NewStatusHandler(campaignStatuses campaignStatusGetter) StatusHandler {
	return StatusHandler{
		campaignStatuses: campaignStatuses,
	}
}

func (h StatusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	splitURL := strings.Split(req.URL.Path, "/")
	campaignID := splitURL[len(splitURL)-2]
	conn := context.Get("database").(collections.DatabaseInterface).Connection()

	status, err := h.campaignStatuses.Get(conn, campaignID)
	if err != nil {
		panic(err)
	}

	output := map[string]interface{}{
		"id":              status.CampaignID,
		"status":          status.Status,
		"total_messages":  status.TotalMessages,
		"sent_messages":   status.SentMessages,
		"retry_messages":  status.RetryMessages,
		"failed_messages": status.FailedMessages,
		"queued_messages": status.QueuedMessages,
		"start_time":      status.StartTime,
		"completed_time":  nil,
	}

	if (status.CompletedTime != mysql.NullTime{}) {
		completedTimeValue, err := status.CompletedTime.Time.MarshalText()
		if err != nil {
			panic(err)
		}

		output["completed_time"] = string(completedTimeValue)
	}

	response, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}