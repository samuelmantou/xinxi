package pkg

import (
	"gorm.io/gorm"
	"log"
	"xinxi/pkg/model"
)

func (p *PinTuan) getDistPidArr() []model.Product {
	var pArr []model.Product
	err := p.db.Where("type = 1").Find(&pArr).Error
	if err != nil {
		log.Println(err)
	}
	return pArr
}

func (p *PinTuan) addWait(w *model.Wait) {
	round := 1
	var lastZj model.Zj
	err := p.db.Where("dest_product_id = ?", w.DestProductId).Find(&lastZj).Error
	if err == nil {
		round = lastZj.Round
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
	}
	w.Round = round + 1
	var lastPd model.Wait
	err = p.db.Where("dest_product_id = ? AND round = ?", w.DestProductId, w.Round).Find(&lastPd).Error
	if err == gorm.ErrRecordNotFound || lastPd.Id == 0 {
		w.Round = round
		w.Index = 1
	}else{
		w.Index = lastPd.Index + 1
	}

	if err := p.db.Create(w).Error; err != nil {
		log.Println(err)
	}
}

func (p *PinTuan) startNew(DestProductId int) {
	var nArr []model.New
	err := p.db.Where("status = ? AND dest_product_id = ?", model.NewStatusNormal, DestProductId).Find(&nArr).Error
	if err == gorm.ErrRecordNotFound || len(nArr) == 0 {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	for _, n := range nArr {
		w := &model.Wait{
			Status: model.WaitStatusNew,
			DestProductId: n.DestProductId,
			OrderId: n.OrderId,
		}
		p.addWait(w)
		if err := p.db.Model(&model.Wait{}).Where("id = ?", w.Id).Update("status", model.NewStatusFinish).Error; err != nil {
			log.Println(err)
		}
	}
}

func (p *PinTuan) startMiss(DestProductId int) {
	var msArr []model.Wait
	err := p.db.Where("status = ? AND dest_product_id = ?", model.WaitStatusMiss, DestProductId).Order("id asc").Find(&msArr).Error
	if err == gorm.ErrRecordNotFound || len(msArr) == 0 {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	for _, w := range msArr {
		pd := &model.Wait{
			Status: model.WaitStatusMissIn,
			DestProductId: w.DestProductId,
			OrderId: w.OrderId,
		}
		p.addWait(pd)
		if err := p.db.Model(&model.Wait{}).Where("id = ?", w.Id).Update("status", model.WaitStatusMissFinish).Error; err != nil {
			log.Println(err)
		}
	}
}

func (p *PinTuan) startLost(DestProductId int) {
	var l model.Zj
	err := p.db.Where("status = ? AND position = ? AND dest_product_id = ?", model.ZjStatusLost, p.lastPosition, DestProductId).Find(&l).Error
	if err == gorm.ErrRecordNotFound || l.Id == 0 {
		p.changeNextC<- struct{}{}
		return
	}
	w := &model.Wait{
		Status: model.WaitStatusMissIn,
		DestProductId: l.DestProductId,
		OrderId: l.OrderId,
	}
	p.addWait(w)
	if err := p.db.Model(&model.Zj{}).Where("id = ?", l.Id).Update("status", model.ZjStatusLostFinish).Error; err != nil {
		log.Println(err)
	}
}

func (p *PinTuan) zj(DestProductId int) {
	round := 1
	var lastZj model.Zj
	err := p.db.Where("dest_product_id = ?", DestProductId).Find(&lastZj).Error
	if err == nil {
		round = lastZj.Round
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
	}
	round = round + 1

	var pdArr []model.Wait
	err = p.db.Where(
		"dest_product_id = ? AND round = ?",
		DestProductId,
		round,
	).Order("id asc").Find(&pdArr).Error
	if err == gorm.ErrRecordNotFound || len(pdArr) == 0 {
		return
	}

	j := 0
	for i := 0; i < len(pdArr) / 4; i++ {
		for j = i * 4; j < i * 4 + 4; j++ {
			pd := pdArr[j]
			z := &model.Zj{
				Round: pd.Round,
				Group: pd.Group,
				Position: pd.Position,
				OrderId: pd.OrderId,
				DestProductId: DestProductId,
				Index: pd.Index,
				Status: model.ZjStatusWait,
			}

			p.db.Create(&z)
		}
	}
	for ; j < len(pdArr); j++ {
		pd := pdArr[j]
		err := p.db.Model(&model.Wait{}).Where("id = ?", pd.Id).Update("status", model.WaitStatusMiss).Error
		if err != nil {
			log.Println(err)
		}
	}
}


func (p *PinTuan) kj(DestProductId, win int) {
	round := 1
	var lastZj model.Zj
	err := p.db.Where("dest_product_id = ?", DestProductId).Find(&lastZj).Error
	if err == nil {
		round = lastZj.Round
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
	}
	round = round + 1

	var pdArr []model.Wait
	err = p.db.Where(
		"dest_product_id = ? AND round = ?",
		DestProductId,
		round,
	).Order("id asc").Find(&pdArr).Error
	if err == gorm.ErrRecordNotFound || len(pdArr) == 0 {
		return
	}

	j := 0
	for i := 0; i < len(pdArr) / 4; i++ {
		for j = i * 4; j < i * 4 + 4; j++ {
			pd := pdArr[j]
			z := &model.Zj{
				Round: pd.Round,
				Group: pd.Group,
				Position: pd.Position,
				OrderId: pd.OrderId,
				DestProductId: DestProductId,
				Index: pd.Index,
			}
			if j == win {
				z.Status = model.ZjStatusWin
			}else{
				z.Status = model.ZjStatusLost
			}
			p.db.Create(&z)
		}
	}
	for ; j < len(pdArr); j++ {
		pd := pdArr[j]
		err := p.db.Model(&model.Wait{}).Where("id = ?", pd.Id).Update("status", model.WaitStatusMiss).Error
		if err != nil {
			log.Println(err)
		}
	}
}