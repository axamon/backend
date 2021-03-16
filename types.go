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

type Processo struct {
	Id           string
	Titolo       string
	Descrizione  string
	Testo        string
	Autori       []string
	Verificatori []string
	Approvatori  []string
	Versione     uint
	Input        *Processo
	Output       *Processo
	Raci         []Attivita
	Status       StatusType
	Kpis         []Kpi
	Created_at   time.Time
	Updated_at   time.Time
}

type Attivita struct {
	Id          string
	Num         int
	UO          string
	Titolo      string
	Descrizione string
	Ruolo       string
	Input       *Attivita
	Output      *Attivita
}

type Kpi struct {
}
