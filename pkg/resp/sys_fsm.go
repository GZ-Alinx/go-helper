package resp

import "github.com/golang-module/carbon"

type FsmApprovalLog struct {
	End             bool `json:"end"`             // is ended?
	WaitingConfirm  bool `json:"waitingConfirm"`  // is waiting submitter confirm?
	WaitingResubmit bool `json:"waitingResubmit"` // is waiting submitter resubmit?
	Cancel          bool `json:"cancel"`          // is submitter canceled?
}

type FsmLogTrack struct {
	CreatedAt carbon.ToDateTimeString `json:"createdAt"`
	UpdatedAt carbon.ToDateTimeString `json:"updatedAt"`
	Name      string                  `json:"name"`
	Opinion   string                  `json:"opinion"`
	Status    uint                    `json:"status"`
	End       bool                    `json:"end"`
	Cancel    bool                    `json:"cancel"`
}

type FsmMachine struct {
	Base
	Category                   uint   `json:"category"`
	Name                       string `json:"name"`
	SubmitterName              string `json:"submitterName"`
	SubmitterEditFields        string `json:"submitterEditFields"`
	SubmitterConfirm           uint   `json:"submitterConfirm"`
	SubmitterConfirmEditFields string `json:"submitterConfirmEditFields"`
	EventsJson                 string `json:"eventsJson"`
}
