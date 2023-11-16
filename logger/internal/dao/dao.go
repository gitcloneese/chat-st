package dao

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewDao, NewDB, NewKafkaReceiver)

// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
}

// dao dao.
type dao struct {
	db *gorm.DB
}

// New new a dao and return.
func NewDao(db *gorm.DB) (d Dao, cf func(), err error) {
	return newDao(db)
}

//nolint:all
func newDao(db *gorm.DB) (d *dao, cf func(), err error) {
	d = &dao{
		db: db,
	}
	cf = d.Close
	return d, cf, nil
}

// Close close the resource.
func (d *dao) Close() {

}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}
