package main

import "time"

type StatusType int

const (
	Nuovo StatusType = iota
	Verificato
	Approvato
	InVigore
	Superato
)

type RaciResp string

const (
	R RaciResp = "R"
	A RaciResp = "A"
	C RaciResp = "C"
	I RaciResp = "I"
)

type Processo struct {
	Id           string
	Titolo       string
	Descrizione  string
	Testo        string
	Autori       []string
	Verificatori []string
	Approvatori  []string
	Versione     uint
	Input        []string
	Output       []string
	Attivitas    []*Attivita
	Status       StatusType
	Kpis         []Kpi
	Created_at   time.Time
	Updated_at   time.Time
	costo        float64
	tmedio       float64
	devstd       float64
}

type Attivita struct {
	Id          string
	Num         int
	UO          string
	Titolo      string
	Descrizione string
	Ruolo       RaciResp
	Input       []string
	Output      []string
	tmedio      float64
	devstd      float64
	costo       float64
}

type Kpi struct {
}
