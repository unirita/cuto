package converter

import (
	"bytes"
	"testing"
)

func getOpenGatewayId(g *Gateway) string {
	return g.ID() + "_open"
}

func getCloseGatewayId(g *Gateway) string {
	return g.ID() + "_close"
}

func isExistSequenceFlow(d *Definitions, from string, to string) bool {
	for _, flow := range d.Process.Flows {
		if flow.From == from && flow.To == to {
			return true
		}
	}

	return false
}

func assertSequenceFlow(t *testing.T, d *Definitions, from string, to string) {
	if !isExistSequenceFlow(d, from, to) {
		t.Errorf("Expected SequenceFlow[FROM:%s,TO:%s] not exists.", from, to)
	}
}

func TestGenerateDefinitions(t *testing.T) {
	initIndexes()

	// job1->job2->[job3,job4->job5]->[job6,job7]->job8
	job1 := NewJob("test1")
	job2 := NewJob("test2")
	job3 := NewJob("test3")
	job4 := NewJob("test4")
	job5 := NewJob("test5")
	job6 := NewJob("test6")
	job7 := NewJob("test7")
	job8 := NewJob("test8")
	gw1 := NewGateway()
	gw2 := NewGateway()
	job1.SetNext(job2)
	job2.SetNext(gw1)
	gw1.AddPathHead(job3)
	gw1.AddPathHead(job4)
	job4.SetNext(job5)
	gw1.SetNext(gw2)
	gw2.AddPathHead(job6)
	gw2.AddPathHead(job7)
	gw2.SetNext(job8)

	d := GenerateDefinitions(job1)
	assertSequenceFlow(t, d, job1.ID(), job2.ID())
	assertSequenceFlow(t, d, job2.ID(), getOpenGatewayId(gw1))
	assertSequenceFlow(t, d, getOpenGatewayId(gw1), job3.ID())
	assertSequenceFlow(t, d, getOpenGatewayId(gw1), job4.ID())
	assertSequenceFlow(t, d, job4.ID(), job5.ID())
	assertSequenceFlow(t, d, job3.ID(), getCloseGatewayId(gw1))
	assertSequenceFlow(t, d, job5.ID(), getCloseGatewayId(gw1))
	assertSequenceFlow(t, d, getCloseGatewayId(gw1), getOpenGatewayId(gw2))
	assertSequenceFlow(t, d, getOpenGatewayId(gw2), job6.ID())
	assertSequenceFlow(t, d, getOpenGatewayId(gw2), job7.ID())
	assertSequenceFlow(t, d, job6.ID(), getCloseGatewayId(gw2))
	assertSequenceFlow(t, d, job7.ID(), getCloseGatewayId(gw2))
	assertSequenceFlow(t, d, getCloseGatewayId(gw2), job8.ID())
}

func TestExport(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<definitions>
    <process>
        <startEvent id=":start"></startEvent>
        <endEvent id=":end"></endEvent>
        <serviceTask id="job1" name="test1"></serviceTask>
        <serviceTask id="job2" name="test2"></serviceTask>
        <serviceTask id="job3" name="test3"></serviceTask>
        <serviceTask id="job4" name="test4"></serviceTask>
        <parallelGateway id="gw1_open"></parallelGateway>
        <parallelGateway id="gw1_close"></parallelGateway>
        <sequenceFlow sourceRef=":start" targetRef="job1"></sequenceFlow>
        <sequenceFlow sourceRef="job1" targetRef="gw1_open"></sequenceFlow>
        <sequenceFlow sourceRef="gw1_open" targetRef="job2"></sequenceFlow>
        <sequenceFlow sourceRef="job2" targetRef="gw1_close"></sequenceFlow>
        <sequenceFlow sourceRef="gw1_open" targetRef="job3"></sequenceFlow>
        <sequenceFlow sourceRef="job3" targetRef="gw1_close"></sequenceFlow>
        <sequenceFlow sourceRef="gw1_close" targetRef="job4"></sequenceFlow>
        <sequenceFlow sourceRef="job4" targetRef=":end"></sequenceFlow>
    </process>
</definitions>`

	initIndexes()
	job1 := NewJob("test1")
	job2 := NewJob("test2")
	job3 := NewJob("test3")
	job4 := NewJob("test4")
	gw1 := NewGateway()
	job1.SetNext(gw1)
	gw1.AddPathHead(job2)
	gw1.AddPathHead(job3)
	gw1.SetNext(job4)
	d := GenerateDefinitions(job1)

	buf := new(bytes.Buffer)
	err := Export(buf, d)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}

	if buf.String() != expected {
		t.Errorf("Exported BPMN is not expected.")
		t.Logf("Output:\n%s", expected)
	}
}
