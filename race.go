package main

type Race struct {
	Status    string `json:"status"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
