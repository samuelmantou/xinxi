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
	IsRefund bool `gorm:"is_refund"`
	Uid int `json:"uid"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (w *New) TableName() string {
	return "xinxi2_pin_tuan"
}

type PoolStatus int

const (
	PoolNormal PoolStatus = iota + 1
	PoolFinish
)

type Pool struct {
	Id int `gorm:"id"`
	Round int `gorm:"round"`
	Position int `gorm:"position"`
	Group int `gorm:"group"`
	Status PoolStatus `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	IsRefund bool `gorm:"is_refund"`
	Uid int `json:"uid"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (p *Pool) TableName() string {
	return "xinxi2_pin_tuan_pool"
}

type WinStatus int

const (
	WinStatusWin WinStatus = iota + 1
	WinStatusLost
	WinStatusLostFinish
	WinStatusMiss
)

type Win struct {
	Id int `gorm:"id"`
	Index int `gorm:"index"`
	Round int `gorm:"round"`
	Group int `gorm:"group"`
	Position int `gorm:"position"`
	Status WinStatus `gorm:"status"`
	OrderId int `gorm:"order_id"`
	DestProductId int `gorm:"dest_product_id"`
	IsRefund bool `gorm:"is_refund"`
	Uid int `json:"uid"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (z *Win) TableName() string {
	return "xinxi2_pin_tuan_win"
}

type WinLog struct {
	Id int `gorm:"id"`
	Round int `gorm:"round"`
	DestProductId int `gorm:"dest_product_id"`
	IsSuccess bool `gorm:"is_success"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (z *WinLog) TableName() string {
	return "xinxi2_pin_tuan_zj"
}

type Running struct {
	Id int `gorm:"id"`
	DestProductId int `gorm:"dest_product_id"`
	IsOpen bool `gorm:"is_open"`
}

func (r *Running) TableName() string {
	return "xinxi2_pin_tuan_running"
}

type ConfigData struct {
	Id string `gorm:"id"`
	RootKey string `gorm:"root_key"`
	Key string `gorm:"key"`
	Value string `gorm:"value"`
}

func (d *ConfigData) TableName() string {
	return "config_data"
}
