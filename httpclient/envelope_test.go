package httpclient

import "testing"

func TestUpstreamErrorUsesFirstMessage(t *testing.T) {
	err := UpstreamError("restaurant", "", "not found")
	if err == nil || err.Error() != "restaurant upstream error: not found" {
		t.Fatalf("UpstreamError() = %v", err)
	}
}

func TestUpstreamErrorFallsBackToServiceName(t *testing.T) {
	err := UpstreamError("pay", "", "")
	if err == nil || err.Error() != "pay upstream error" {
		t.Fatalf("UpstreamError() = %v", err)
	}
}

func TestEnvelopeShape(t *testing.T) {
	envelope := Envelope[string]{
		Code:      0,
		Message:   "success",
		Data:      "ok",
		RequestID: "rid-1",
	}

	if envelope.Code != 0 || envelope.Message != "success" || envelope.Data != "ok" || envelope.RequestID != "rid-1" {
		t.Fatalf("envelope = %+v", envelope)
	}
}
