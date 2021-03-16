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

	// crea attività
	SpendereSoldi := Attivita{
		Id:          uuid.NewString(),
		UO:          "CTIO.5GDT.PDT",
		Num:         1,
		Ruolo:       R,
		Titolo:      "Impiego Budget",
		Descrizione: "Spendere e spandere a vanvera",
	}

	SpendereAltriSoldi := Attivita{
		Id:          uuid.NewString(),
		UO:          "CTIO.5GDT.PDO",
		Num:         2,
		Ruolo:       R,
		Titolo:      "Impiego Budget in operazioni nere",
		Descrizione: "Spendere e spandere sempre più a vanvera",
	}

	// aggiunge attività a processo
	CicloPassivo.Attivitas = append(CicloPassivo.Attivitas, &SpendereSoldi, &SpendereAltriSoldi)

	fmt.Println(CicloPassivo.UOCoinvolte())

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
