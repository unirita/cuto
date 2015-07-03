package converter

import "encoding/xml"

// Definitions element in BPMN.
type Definitions struct {
	XMLName xml.Name `xml:"definitions"`
	Process *Process `xml:"process"`
}

// NewDefinitions create empty Definitions object.
func NewDefinitions() *Definitions {
	d := new(Definitions)
	d.Process = NewProcess()
	return d
}

// AppendServiceTask appends ServiceTask object to Definitions#Process.
func (d *Definitions) AppendServiceTask(task *ServiceTask) {
	d.Process.Tasks = append(d.Process.Tasks, task)
}

// AppendParallelGateway appends ParallelGateway object to Definitions#Process.
func (d *Definitions) AppendParallelGateway(gateway *ParallelGateway) {
	d.Process.Gateways = append(d.Process.Gateways, gateway)
}

// AppendSequenceFlow appends SequenceFlow object to Definitions#Process.
func (d *Definitions) AppendSequenceFlow(flow *SequenceFlow) {
	d.Process.Flows = append(d.Process.Flows, flow)
}

// AppendJob appends job element to Definitions object.
// This function returns job ID to create connection with next element.
func (d *Definitions) AppendJob(job *Job, pre string) string {
	s := NewServiceTask(job)
	d.AppendServiceTask(s)
	d.AppendSequenceFlow(NewSequenceFlow(pre, s.ID))
	return s.ID
}

// AppendGateway appends gateway element to Definitions object.
// This function returns close gateway ID to create connection with next element.
func (d *Definitions) AppendGateway(gw *Gateway, pre string) string {
	openGW, closeGW := NewParallelGatewayPair(gw)
	d.AppendParallelGateway(openGW)
	d.AppendParallelGateway(closeGW)
	d.AppendSequenceFlow(NewSequenceFlow(pre, openGW.ID))

	// Append inner paths.
	for _, pathHead := range gw.PathHeads {
		preInPath := openGW.ID
		current := pathHead
		for {
			d.AppendJob(current, preInPath)
			preInPath = current.ID()
			var ok bool

			next := current.Next()
			if next == nil {
				break
			}

			current, ok = next.(*Job)
			if !ok {
				panic("Nest of gateway detected.")
			}
		}
		d.AppendSequenceFlow(NewSequenceFlow(preInPath, closeGW.ID))
	}
	return closeGW.ID
}

// Process element in BPMN.
type Process struct {
	Start    *StartEvent        `xml:"startEvent"`
	End      *EndEvent          `xml:"endEvent"`
	Tasks    []*ServiceTask     `xml:"serviceTask"`
	Gateways []*ParallelGateway `xml:"parallelGateway"`
	Flows    []*SequenceFlow    `xml:"sequenceFlow"`
}

// Create empty Process object.
func NewProcess() *Process {
	p := new(Process)
	p.Start = NewStartEvent()
	p.End = NewEndEvent()
	p.Tasks = make([]*ServiceTask, 0)
	p.Gateways = make([]*ParallelGateway, 0)
	p.Flows = make([]*SequenceFlow, 0)
	return p
}

// StartEvent element in BPMN.
type StartEvent struct {
	ID string `xml:"id,attr"`
}

// Create StartEvent object with unique ID.
func NewStartEvent() *StartEvent {
	s := new(StartEvent)
	s.ID = ":start"
	return s
}

// EndEvent element in BPMN.
type EndEvent struct {
	ID string `xml:"id,attr"`
}

// Create EndEvent object with unique ID.
func NewEndEvent() *EndEvent {
	e := new(EndEvent)
	e.ID = ":end"
	return e
}

// ServiceTask element in BPMN.
type ServiceTask struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// Create ServiceTask object from job element.
func NewServiceTask(job *Job) *ServiceTask {
	s := new(ServiceTask)
	s.ID = job.ID()
	s.Name = job.Name()
	return s
}

// ParallelGateway element in BPMN.
type ParallelGateway struct {
	ID string `xml:"id,attr"`
}

// Create pair of ParallelGateway to open and close parallel paths.
// This function has two return values.
// First one is gateway for open, second one is gateway for close.
func NewParallelGatewayPair(gw *Gateway) (*ParallelGateway, *ParallelGateway) {
	const (
		openSuffix  = "_open"
		closeSuffix = "_close"
	)

	openGW := new(ParallelGateway)
	openGW.ID = gw.ID() + openSuffix
	closeGW := new(ParallelGateway)
	closeGW.ID = gw.ID() + closeSuffix
	return openGW, closeGW
}

// SequenceFlow element in BPMN.
type SequenceFlow struct {
	From string `xml:"sourceRef,attr"`
	To   string `xml:"targetRef,attr"`
}

// Create SequenceFlow
func NewSequenceFlow(from, to string) *SequenceFlow {
	s := new(SequenceFlow)
	s.From = from
	s.To = to
	return s
}
