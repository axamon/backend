package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Metodi

// UOCoinvolte restituisce la lista delle Unità organizzative
// coinvolte in un processo.
func (p Processo) UOCoinvolte() []string {
	m := make(map[string]struct{})
	for _, a := range p.Attivitas {
		m[a.UO] = struct{}{}
	}
	var uos []string
	for uo := range m {
		uos = append(uos, uo)
	}
	return uos
}

func (p Processo) HaAMonte(p2 Processo) {
	// aggiorna processo a monte
	if !find(p2.Output, p.Id) { // se non è già presente
		p2.Output = append(p2.Output, p.Id)
		p2.Update()
	}
	// aggiorna processo a valle
	if !find(p.Input, p2.Id) {
		p.Input = append(p.Input, p2.Id)
		p.Update()
	}
}

func (p Processo) HaAValle(p2 Processo) {
	if !find(p.Output, p2.Id) {
		p.Output = append(p.Output, p2.Id)
		p.Update()
	}
	if !find(p2.Input, p.Id) {
		p2.Input = append(p2.Input, p.Id)
		p2.Update()
	}
}

func (p Processo) Approva() {
	p.Status = Approvato
}

func (p Processo) Ver() uint {
	return p.Versione
}

func (p Processo) Delete() {
	err := c.Remove(bson.M{"id": p.Id})
	if err != nil {
		log.Printf("Delete di %s in errore: %v\n", p.Titolo, err)
	}
}

func (p Processo) Update() error {
	p.Versione++
	p.Updated_at = time.Now()

	err := c.Update(bson.M{"id": p.Id}, &p)
	if err != nil {
		log.Printf("Update di %s in errore: %v\n", p.Id, err)
	}
	return err
}

func (p Processo) Save() error {
	err := c.Insert(&p)
	if err != nil {
		log.Printf("Save di %s in errore: %v\n", p.Id, err)

	}
	return err
}
