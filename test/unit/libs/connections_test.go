package libs

import (
	"testing"

	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/anand-kay/linkedin-clone/models"
)

func TestExtractSentAndReceivedReqs(t *testing.T) {
	connections := []models.Connection{
		{User1: "4", User2: "9"},
		{User1: "6", User2: "4"},
		{User1: "1", User2: "4"},
	}

	sentReqs, receivedReqs := libs.ExtractSentAndReceivedReqs(connections, "4")

	if len(sentReqs) != 1 {
		t.Error("ExtractSentAndReceivedReqs - Failed")
	}
	if sentReqs[0] != "9" {
		t.Error("ExtractSentAndReceivedReqs - Failed")
	}

	if len(receivedReqs) != 2 {
		t.Error("ExtractSentAndReceivedReqs - Failed")
	}
	if receivedReqs[0] != "6" && receivedReqs[0] != "1" {
		t.Error("ExtractSentAndReceivedReqs - Failed")
	}
	if receivedReqs[1] != "6" && receivedReqs[1] != "1" {
		t.Error("ExtractSentAndReceivedReqs - Failed")
	}
}
