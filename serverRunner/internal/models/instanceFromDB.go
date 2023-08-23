package models

import (
	"database/sql"

	"github.com/3AM-Developer/server-runner/internal/database"
	"github.com/3AM-Developer/server-runner/internal/instance"
)

// Still figuring this out
// Also, update for Gorm

var InstanceDB = DB{database.Db}

var (
	InvalidInstanceError error
)

type DB struct {
	Db *sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

func (d *DB) GetInstanceById(id int) (*instance.Instance, error) {
	query := "SELECT name, dir FROM instances WHERE id = ?"
	row := d.Db.QueryRow(query, id)

	qRes := &instance.Instance{
		Id: id,
	}
	err := row.Scan(&qRes.Name, &qRes.Dir)
	if err != nil {
		return nil, err
	}

	return qRes, err

}

func (d *DB) GetInstanceByName(name string) (*instance.Instance, error) {
	query := "SELECT id, dir FROM instances WHERE name = ?"
	row := d.Db.QueryRow(query, name)
	qRes := &instance.Instance{
		Name: name,
	}

	err := row.Scan(&qRes.Id, &qRes.Dir)
	if err != nil {
		return nil, err
	}

	return qRes, err

}

func (d *DB) NewInstance(i *instance.Instance) error {
	if !i.VerifyInstance() {
		return InvalidInstanceError
	}

	// Start a new transaction
	tx, err := d.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw the panic
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	insert := "INSERT INTO instances (name, dir) VALUES (?, ?)"
	_, err = tx.Exec(insert, i.Name, i.Dir)
	if err != nil {
		return err
	}

	err = i.Write()
	if err != nil {
		return err
	}

	return nil
}
