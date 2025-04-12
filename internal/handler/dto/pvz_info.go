package dto

type PvzInfoRequest struct {
	StartDate string `schema:"startDate" validate:"omitempty"`
	EndDate   string `schema:"endDate"   validate:"omitempty"`
	Page      int    `schema:"page"      validate:"omitempty"`
	Limit     int    `schema:"limit"     validate:"omitempty"`
}

type PvzInfoResponse struct {
	Data []SinglePvzInfoResponse
}

type SinglePvzInfoResponse struct {
	PvzData    PvzResponse     `json:"pvz"`
	Receptions []ReceptionInfo `json:"receptions"`
}
type ReceptionInfo struct {
	ReceptionData ReceptionResponse `json:"reception"`
	Products      []ProductResponse `json:"products"`
}
