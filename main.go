package main

import (
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

var c *mgo.Collection
var database, collection = "db", "processi"
var mongoURL = "mongodb://localhost"

func main() {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Print("session", err)
	}

	c = session.DB(database).C(collection)

	deleteAllProcessi()

	// Crea processi
	CicloPassivo, err := NewProcesso("Ciclo passivo")
	if err != nil {
		log.Println(err)
	}
	VerificaBudget, err := NewProcesso("Verifica budget")
	if err != nil {
		log.Println(err)
	}

	// Collega processi tra loro
	CicloPassivo.HaAValle(VerificaBudget)
	VerificaBudget.HaAMonte(CicloPassivo)

	p, _ := GetProcesso(CicloPassivo.Id)
	p2, _ := GetProcesso(VerificaBudget.Id)
	fmt.Printf("%+v\n", p)
	fmt.Printf("%+v\n", p2)

	CicloPassivo.Delete()

	// pp, err := GetAllProcessi()
	// for i, p := range pp {
	// 	fmt.Println(i, p)
	// }

	//	p.Delete()

}

// NewProcesso crea un nuovo processo.
func NewProcesso(titolo string) (p Processo, err error) {
	// Verifica che non esistano processi con lo stesso nome.
	processi, err := GetAllProcessi()
	if err != nil {
		return Processo{}, err
	}
	for _, processo := range processi {
		if processo.Titolo == titolo {
			return Processo{}, fmt.Errorf("titolo \"%s\" già esistente con id %s", titolo, processo.Id)
		}
	}
	p.Id = uuid.NewString()
	p.Titolo = titolo
	p.Versione = 1
	p.Status = Nuovo
	p.Created_at = time.Now()
	err = p.Save()
	return p, err
}

// GetProcesso recupera il processo con id.
func GetProcesso(id string) (Processo, error) {
	var p Processo
	err := c.Find(bson.M{"id": id}).One(&p)
	if err != nil {
		log.Printf("GetProcesso per id: %s in errore: %v \n", id, err)
	}
	return p, err
}

// GetAllProcessi recupera tutti i processi.
func GetAllProcessi() ([]Processo, error) {
	var processi []Processo
	err := c.Find(nil).All(&processi)
	if err != nil {
		log.Printf("GetAllProcessi in errore: %v \n", err)
	}
	return processi, err
}

// UpdateProcesso modifica un processo.
func UpdateProcesso(id string, p Processo) {
	p.Id = id
	p.Versione++
	p.Updated_at = time.Now()

	c.Update(bson.M{"id": id}, &p)
}

// DeleteProcesso cancella un processo dal db.
func DeleteProcesso(id string) error {
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		log.Print(err)
	}
	return err
}

// deleteAllProcessi recupera tutti i processi.
func deleteAllProcessi() error {
	var processi []Processo
	err := c.Find(nil).All(&processi)
	if err != nil {
		log.Printf("GetAllProcessi in errore: %v \n", err)
	}
	for _, p := range processi {
		c.Remove(bson.M{"id": p.Id})
	}
	return err
}

func find(s []string, val string) bool {
	for _, item := range s {
		if item == val {
			return true
		}
	}
	return false
}

// Metodi

// UOCoinvolte restituisce la lista delle Unità organizzative
// coinvolte in un processo.
func (p Processo) UOCoinvolte() []string {
	m := make(map[string]struct{})
	for _, a := range p.Raci {
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
