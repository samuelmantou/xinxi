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

type NewStatus int

const (
	NewStatusNormal NewStatus = iota
	NewStatusFinish
)

type New struct {
	Id int `gorm:"id"`
	Status int `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (w *New) TableName() string {
	return "xinxi2_pin_tuan"
}

type WaitStatus int

const (
	WaitStatusNew WaitStatus = iota + 1
	WaitStatusLost
	WaitStatusMiss
	WaitStatusMissIn
	WaitStatusMissFinish
)

type Wait struct {
	Id int `gorm:"id"`
	Index int `gorm:"index"`
	Round int `gorm:"round"`
	Group int `gorm:"group"`
	Position int `gorm:"position"`
	Status WaitStatus `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (p *Wait) TableName() string {
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
