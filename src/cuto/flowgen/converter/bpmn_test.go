package converter

import "testing"

func TestNewDefinitions(t *testing.T) {
	d := NewDefinitions()
	if d.Process == nil {
		t.Error("Created Definitions has nil Process.")
	}
}

func TestAppendServiceTask(t *testing.T) {
	d := NewDefinitions()
	d.AppendServiceTask(NewServiceTask(NewJob("test1")))
	d.AppendServiceTask(NewServiceTask(NewJob("test2")))
	if len(d.Process.Tasks) != 2 {
		t.Fatalf("Number of tasks[%d] is not equal to expected[%d].",
			len(d.Process.Tasks), 2)
	}
}

func TestAppendParallelGateway(t *testing.T) {
	d := NewDefinitions()
	ogw, cgw := NewParallelGatewayPair(NewGateway())
	d.AppendParallelGateway(ogw)
	d.AppendParallelGateway(cgw)
	if len(d.Process.Gateways) != 2 {
		t.Fatalf("Number of gateway[%d] is not equal to expected[%d].",
			len(d.Process.Gateways), 2)
	}
}

func TestAppendSequenceFlow(t *testing.T) {
	d := NewDefinitions()
	d.AppendSequenceFlow(NewSequenceFlow("from1", "to1"))
	d.AppendSequenceFlow(NewSequenceFlow("from2", "to2"))
	if len(d.Process.Flows) != 2 {
		t.Fatalf("Number of flow[%d] is not equal to expected[%d].",
			len(d.Process.Flows), 2)
	}
}

func TestAppendJob(t *testing.T) {
	d := NewDefinitions()
	j := NewJob("test")
	d.AppendJob(j, "pre")
	if len(d.Process.Tasks) != 1 {
		t.Fatalf("ServiceTask was not appended.")
	}
	if len(d.Process.Flows) != 1 {
		t.Fatalf("SequenceFlow was not appended.")
	}
}

func TestAppendGateway(t *testing.T) {
	d := NewDefinitions()
	g := NewGateway()
	g.AddPathHead(NewJob("test1"))
	g.AddPathHead(NewJob("test2"))
	g.AddPathHead(NewJob("test3"))
	g.PathHeads[2].SetNext(NewJob("test4"))
	d.AppendGateway(g, "pre")
	if len(d.Process.Tasks) != 4 {
		t.Logf("Number of ServiceTask is %d", len(d.Process.Tasks))
		t.Fatalf("ServiceTask was not appended enough.")
	}
	if len(d.Process.Gateways) != 2 {
		t.Logf("Number of ParallelGateways is %d", len(d.Process.Gateways))
		t.Fatalf("ParallelGateways was not appended enough.")
	}
	if len(d.Process.Flows) != 8 {
		t.Logf("Number of SequenceFlow is %d", len(d.Process.Flows))
		t.Fatalf("SequenceFlow was not appended enough.")
	}
}

func TestNewProcess(t *testing.T) {
	p := NewProcess()
	if p.Start == nil {
		t.Error("Created Process has nil StartEvent.")
	}
	if p.End == nil {
		t.Error("Created Process has nil EndEvent.")
	}
	if len(p.Tasks) != 0 {
		t.Error("Created Process has unexpected ServiceTasks.")
	}
	if len(p.Gateways) != 0 {
		t.Error("Created Process has unexpected ParallelGateways.")
	}
	if len(p.Flows) != 0 {
		t.Error("Created Process has unexpected SequenceFlows.")
	}
}

func TestNewStartEvent(t *testing.T) {
	s := NewStartEvent()
	if !isIrregalJobName(s.ID) {
		t.Error("StartEvent ID[%s] has possibility of using for job ID or Name.", s.ID)
	}
}

func TestNewEndEvent(t *testing.T) {
	e := NewEndEvent()
	if !isIrregalJobName(e.ID) {
		t.Error("EndEvent ID[%s] has possibility of using for job ID or Name.", e.ID)
	}
}

func TestNewServiceTask(t *testing.T) {
	initIndexes()
	j := NewJob("test")
	s := NewServiceTask(j)
	if s.ID != "job1" {
		t.Errorf("Created ServiceTasks has unexpected ID[%s].", s.ID)
	}
	if s.Name != "test" {
		t.Errorf("Created ServiceTasks has unexpected ID[%s].", s.Name)
	}
}

func TestNewParallelGatewayPair(t *testing.T) {
	initIndexes()
	g := NewGateway()
	openGW, closeGW := NewParallelGatewayPair(g)
	if openGW.ID != "gw1_open" {
		t.Errorf("Created ParallelGateway for open has unexpected ID[%s].", openGW.ID)
	}
	if closeGW.ID != "gw1_close" {
		t.Errorf("Created ParallelGateway for close has unexpected ID[%s].", closeGW.ID)
	}
}

func TestNewSequenceFlow(t *testing.T) {
	s := NewSequenceFlow("fromid", "toid")
	if s.From != "fromid" {
		t.Errorf("Created SequenceFlow has unexpected From[%s].", s.From)
	}
	if s.To != "toid" {
		t.Errorf("Created SequenceFlow has unexpected To[%s].", s.To)
	}
}
