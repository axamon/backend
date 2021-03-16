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
var database, collection = "dbProcessi", "processi"
var mongoURL = "mongodb://localhost"

func main() {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Print("session", err)
	}

	c = session.DB(database).C(collection)

	var p = NewProcesso()

	p.Titolo = "Cloud"

	p.Save()

	fmt.Println(p)

	var id = p.Id

	t, _ := GetProcesso(id)

	t.Autori = []string{"Alberto Bregliano"}

	UpdateProcesso(id, t)

	p, _ = GetProcesso(id)

	fmt.Printf("%v\n", p)

	// pp, err := GetAllProcessi()
	// for i, p := range pp {
	// 	fmt.Println(i, p)
	// }

	//	p.Delete()

}

// NewProcesso crea un nuovo processo.
func NewProcesso() (p Processo) {
	p.Id = uuid.NewString()
	p.Versione = 1
	p.Status = Nuovo
	p.Created_at = time.Now()
	return p
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

func (p Processo) Approva() {
	p.Status = Approvato
}

func (p Processo) Ver() uint {
	return p.Versione
}

func (p Processo) Delete() {
	err := c.Remove(bson.M{"id": p.Id})
	if err != nil {
		log.Print(err)
	}
}

func (p Processo) Save() {
	err := c.Insert(&p)
	if err != nil {
		log.Print(err)
	}
}
