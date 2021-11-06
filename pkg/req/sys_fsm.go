package req

import "github.com/piupuer/go-helper/pkg/resp"

type FsmCreateMachine struct {
	Category                   NullUint         `json:"category"`
	Name                       string           `json:"name"`
	SubmitterName              string           `json:"submitterName"`
	SubmitterEditFields        string           `json:"submitterEditFields"`
	SubmitterConfirm           NullUint         `json:"submitterConfirm"`
	SubmitterConfirmEditFields string           `json:"submitterConfirmEditFields"`
	Levels                     []FsmCreateEvent `json:"levels"`
}

type FsmCreateEvent struct {
	Name       string   `json:"name" form:"name"`
	Edit       NullUint `json:"edit" form:"edit"`
	Refuse     NullUint `json:"refuse" form:"refuse"`
	EditFields string   `json:"editFields" form:"editFields"`
	Roles      IdsStr   `json:"roles" form:"roles"`
	Users      IdsStr   `json:"users" form:"users"`
}

type FsmUpdateMachine struct {
	Name                       *string          `json:"name"`
	SubmitterName              *string          `json:"submitterName"`
	SubmitterEditFields        *string          `json:"submitterEditFields"`
	SubmitterConfirm           *NullUint        `json:"submitterConfirm"`
	SubmitterConfirmEditFields *string          `json:"submitterConfirmEditFields"`
	Levels                     []FsmCreateEvent `json:"levels"`
}

type FsmCreateLog struct {
	Category        NullUint `json:"category" form:"category"`
	Uuid            string   `json:"uuid" form:"uuid"`
	SubmitterRoleId uint     `json:"submitterRoleId" form:"submitterRoleId"`
	SubmitterUserId uint     `json:"submitterUserId" form:"submitterUserId"`
}

type FsmApproveLog struct {
	Category        NullUint `json:"category" form:"category"`
	Uuid            string   `json:"uuid" form:"uuid"`
	ApprovalRoleId  uint     `json:"approvalRoleId" form:"approvalRoleId"`
	ApprovalUserId  uint     `json:"approvalUserId" form:"approvalUserId"`
	ApprovalOpinion string   `json:"approvalOpinion" form:"approvalOpinion"`
	Approved        NullUint `json:"approved" form:"approved"`
}

type FsmSubmitterDetail struct {
	Category NullUint `json:"category" form:"category"`
	Uuid     string   `json:"uuid" form:"uuid"`
}

type FsmSubmitterDetailField struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type UpdateFsmSubmitterDetail struct {
	FsmSubmitterDetail
	Fields []FsmSubmitterDetailField `json:"fields"`
	Keys   []string                  `json:"-"`
	Vals   []string                  `json:"-"`
}

func (d *UpdateFsmSubmitterDetail) Parse() {
	k := make([]string, len(d.Fields))
	v := make([]string, len(d.Fields))
	for i, field := range d.Fields {
		k[i] = field.Key
		v[i] = field.Val
	}
	d.Keys = k
	d.Vals = v
}

type FsmPermissionLog struct {
	Category       NullUint `json:"category" form:"category"`
	Uuid           string   `json:"uuid" form:"uuid"`
	ApprovalRoleId uint     `json:"approvalRoleId" form:"approvalRoleId"`
	ApprovalUserId uint     `json:"approvalUserId" form:"approvalUserId"`
	Approved       uint     `json:"approved" form:"approved"`
}

type FsmPendingLog struct {
	ApprovalRoleId uint     `json:"approvalRoleId" form:"approvalRoleId"`
	ApprovalUserId uint     `json:"approvalUserId" form:"approvalUserId"`
	Category       NullUint `json:"category" form:"category"`
	resp.Page
}

type FsmLog struct {
	Category NullUint `json:"category" form:"category"`
	Uuid     string   `json:"uuid" form:"uuid"`
}

type FsmMachine struct {
	Category         *NullUint `json:"category" form:"category"`
	Name             string    `json:"name" form:"name"`
	SubmitterName    string    `json:"submitterName" form:"submitterName"`
	SubmitterConfirm *NullUint `json:"submitterConfirm" form:"submitterConfirm"`
	resp.Page
}
