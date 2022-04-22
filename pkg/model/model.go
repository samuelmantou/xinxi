package model

import "time"

type Product struct {
	Id int `gorm:"id"`
	Name string `gorm:"name"`
	Type int `gorm:"type"`
}

func (p *Product) TableName() string {
	return "yd_product"
}

type WaitStatus int

const (
	WaitStatusNormal WaitStatus = iota
	WaitStatusFinish
)

type Wait struct {
	Id int `gorm:"id"`
	Status int `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (w *Wait) TableName() string {
	return "xinxi2_pin_tuan"
}

type PdStatus int

const (
	PdStatusNew PdStatus = iota + 1
	PdStatusLost
	PdStatusMiss
	PdStatusMissIn
	PdStatusMissFinish
)


type Pd struct {
	Id int `gorm:"id"`
	Index int `gorm:"index"`
	Round int `gorm:"round"`
	Group int `gorm:"group"`
	Position int `gorm:"position"`
	Status PdStatus `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (p *Pd) TableName() string {
	return "xinxi2_pin_tuan_history"
}

type ZjStatus int

const (
	ZjStatusWait ZjStatus = iota
	ZjStatusWin
	ZjStatusLost
	ZjStatusLostFinish
)

type Zj struct {
	Id int `gorm:"id"`
	Index int `gorm:"index"`
	Round int `gorm:"round"`
	Group int `gorm:"group"`
	Position int `gorm:"position"`
	Status ZjStatus `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (z *Zj) TableName() string {
	return "xinxi2_pin_tuan_history2"
}
