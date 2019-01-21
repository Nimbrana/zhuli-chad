// Package signaling establish the connection between clients of WebRTC protocol
// This project is based on signaling-go project from sfRTC-framework of Karina Romero (@KarinaRomero)
// The original github repository is https://github.com/KarinaRomero/signaling-go
package signaling

type loginMessage struct {
	Type   string `json:"type"`
	Succes bool   `json:"success"`
}

type offer struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

type offerMessage struct {
	Type  string `json:"type"`
	Offer offer  `json:"offer"`
	Name  string `json:"name"`
}

type answer struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

type answerMessage struct {
	Type   string `json:"type"`
	Answer answer `json:"answer"`
}

type candidate struct {
	Candidate     string  `json:"candidate"`
	SdpMid        string  `json:"sdpMid"`
	SdpMLineIndex float64 `json:"sdpMLineIndex"`
}

type candidateMessage struct {
	Type      string    `json:"type"`
	Candidate candidate `json:"candidate"`
}

type leaveMessage struct {
	Type string `json:"type"`
}

type unrecognizedMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
