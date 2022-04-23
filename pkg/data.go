package pkg

import (
	"gorm.io/gorm"
	"log"
	"xinxi/pkg/model"
)

func (p *PinTuan) getDistPidArr() []model.Product {
	var pArr []model.Product
	err := p.db.Debug().Where("type = 1").Find(&pArr).Error
	if err != nil {
		log.Println(err)
	}
	return pArr
}

func (p *PinTuan) insertPool(po *model.Pool) {
	if err := p.db.Debug().Create(po).Error; err != nil {
		log.Println(err)
	}
}

func (p *PinTuan) insertNew(DestProductId int) {
	var nArr []model.New
	err := p.db.Debug().Where("status = ? AND dest_product_id = ?", model.NewStatusNormal, DestProductId).Find(&nArr).Error
	if err == gorm.ErrRecordNotFound || len(nArr) == 0 {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	for _, n := range nArr {
		po := &model.Pool{
			Status: model.PoolNormal,
			DestProductId: n.DestProductId,
			OrderId: n.OrderId,
		}
		p.insertPool(po)
		if err := p.db.Debug().Model(&model.New{}).Where("id = ?", n.Id).Update("status", model.NewStatusFinish).Error; err != nil {
			log.Println(err)
		}
	}
}

func (p *PinTuan) insertLost(DestProductId int) {
	var w model.Win
	err := p.db.Debug().Where("status = ? AND position = ? AND dest_product_id = ?", model.WinStatusLost, p.lastPosition, DestProductId).
		Order("id asc").
		Find(&w).Error
	if err == gorm.ErrRecordNotFound || w.Id == 0{
		return
	}
	po := &model.Pool{
		Status: model.PoolNormal,
		DestProductId: w.DestProductId,
		OrderId: w.OrderId,
	}
	p.insertPool(po)
	if err := p.db.Debug().Model(&model.Win{}).Where("id = ?", w.Id).Update("status", model.WinStatusLostFinish).Error; err != nil {
		log.Println(err)
	}
}

func (p *PinTuan) open(DestProductId, win int) {
	round := 0
	var lastWin model.Win
	err := p.db.Debug().Where("dest_product_id = ?", DestProductId).Find(&lastWin).Error
	if lastWin.Id > 0 {
		round = lastWin.Round
	}
	round++

	var poolArr []model.Pool
	err = p.db.Debug().Where(
		"dest_product_id = ? AND status = ?",
		DestProductId,
		model.PoolNormal,
	).Order("id asc").Find(&poolArr).Error
	if err == gorm.ErrRecordNotFound || len(poolArr) == 0 {
		return
	}

	j := 0
	idx := 1
	for i := 0; i < len(poolArr) / 4; i++ {
		k := 1
		for j = i * 4; j < i * 4 + 4; j++ {
			po := poolArr[j]
			z := &model.Win{
				Round: round,
				Group: i + 1,
				Position: k,
				OrderId: po.OrderId,
				DestProductId: DestProductId,
				Index: idx,
			}

			if k == win {
				z.Status = model.WinStatusWin
			}else{
				z.Status = model.WinStatusLost
			}

			p.db.Debug().Create(&z)
			idx++
			k++
		}
	}
	
	for ; j < len(poolArr); j++ {
		po := poolArr[j]
		err := p.db.Debug().Model(&model.Win{}).Where("id = ?", po.Id).Update("status", model.WinStatusMiss).Error
		if err != nil {
			log.Println(err)
		}
	}
}