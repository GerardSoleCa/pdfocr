package pdf

import "testing"

func TestProcessor(t *testing.T) {
	processor := Processor{
		InName:  "in.pdf",
		OutName: "out.pdf",
	}

	if processor.InName != "in.pdf" {
		t.Error("Expected in.pdf, got", processor.InName)
	}
	if processor.OutName != "out.pdf" {
		t.Error("Expected in.pdf, got", processor.OutName)
	}
}

func TestProcessor_Process(t *testing.T) {
	processor := Processor{
		InName:  "in.pdf",
		OutName: "out.pdf",
	}

	processor.Process()
}
