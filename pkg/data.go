package pkg

import (
	"gorm.io/gorm"
	"log"
	"strconv"
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
			Uid: n.Uid,
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
		Uid: w.Uid,
	}
	p.insertPool(po)
	if err := p.db.Debug().Model(&model.Win{}).Where("id = ?", w.Id).Update("status", model.WinStatusLostFinish).Error; err != nil {
		log.Println(err)
	}
}

func (p *PinTuan) open(destProductId, win int, reward Reward) {
	round := 0
	var lastWin model.Win
	err := p.db.Debug().Where("dest_product_id = ?", destProductId).Order("round desc").Find(&lastWin).Error
	if lastWin.Id > 0 {
		round = lastWin.Round
	}
	round++

	var poolArr []model.Pool
	err = p.db.Debug().Where(
		"dest_product_id = ? AND status = ?",
		destProductId,
		model.PoolNormal,
	).Order("id asc").Find(&poolArr).Error
	if err == gorm.ErrRecordNotFound || len(poolArr) == 0 {
		return
	}
	if len(poolArr) < 4 {
		return
	}
	j := 0
	idx := 1
	var winIds []string
	for i := 0; i < len(poolArr) / 4; i++ {
		k := 1
		for j = i * 4; j < i * 4 + 4; j++ {
			po := poolArr[j]
			w := &model.Win{
				Round: round,
				Group: i + 1,
				Position: k,
				OrderId: po.OrderId,
				DestProductId: destProductId,
				Index: idx,
				Uid: po.Uid,
			}

			if k == win {
				w.Status = model.WinStatusWin
			}else{
				w.Status = model.WinStatusLost
			}

			p.db.Debug().Create(&w)
			if w.Id > 0 {
				winIds = append(winIds, strconv.Itoa(w.Id))
			}
			p.db.Debug().Model(&model.Pool{}).
				Where("id = ?", po.Id).
				Updates(map[string]interface{}{
					"status": model.PoolFinish,
					"round": w.Round,
					"group": w.Group,
					"position": w.Position,
				})
			idx++
			k++
		}
	}
	if len(winIds) > 0 {
		reward(winIds)
	}
	
	for ; j < len(poolArr); j++ {
		po := poolArr[j]
		w := model.Win{
			Index: idx,
			Round: round,
			OrderId: po.OrderId,
			DestProductId: destProductId,
			Status: model.WinStatusMiss,
			Uid: po.Uid,
		}
		idx++
		err := p.db.Debug().Create(&w)
		if err != nil {
			log.Println(err)
		}
	}

	l := model.WinLog{
		Round: round,
		DestProductId: destProductId,
	}
	if err := p.db.Debug().Create(&l).Error; err != nil {
		log.Println(err)
	}
}