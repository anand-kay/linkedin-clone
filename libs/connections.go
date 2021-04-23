package libs

import (
	"github.com/anand-kay/linkedin-clone/models"
)

// ExtractSentAndReceivedReqs - Extracts sent and received requests from pending requests
func ExtractSentAndReceivedReqs(connections []models.Connection, userID string) ([]string, []string) {
	var sentReqs, receivedReqs []string

	for _, connection := range connections {
		if connection.User1 == userID {
			sentReqs = append(sentReqs, connection.User2)
		} else if connection.User2 == userID {
			receivedReqs = append(receivedReqs, connection.User1)
		}
	}

	return sentReqs, receivedReqs
}
