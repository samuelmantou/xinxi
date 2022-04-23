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

func (p *PinTuan) insertPool(po *model.Pool) {
	if err := p.db.Create(po).Error; err != nil {
		log.Println(err)
	}
}

func (p *PinTuan) insertNew(DestProductId int) {
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
		po := &model.Pool{
			Status: model.PoolNormal,
			DestProductId: n.DestProductId,
			OrderId: n.OrderId,
		}
		p.insertPool(po)
		if err := p.db.Model(&model.New{}).Where("id = ?", n.Id).Update("status", model.NewStatusFinish).Error; err != nil {
			log.Println(err)
		}
	}
}

func (p *PinTuan) insertLost(DestProductId int) {
	var w model.Win
	err := p.db.Where("status = ? AND position = ? AND dest_product_id = ?", model.WinStatusLost, p.lastPosition, DestProductId).Find(&w).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
	po := &model.Pool{
		Status: model.PoolNormal,
		DestProductId: w.DestProductId,
		OrderId: w.OrderId,
	}
	p.insertPool(po)
	if err := p.db.Model(&model.Win{}).Where("id = ?", w.Id).Update("status", model.WinStatusLostFinish).Error; err != nil {
		log.Println(err)
	}
}

//func (p *PinTuan) zj(DestProductId int) {
//	round := 1
//	var lastZj model.Zj
//	err := p.db.Where("dest_product_id = ?", DestProductId).Find(&lastZj).Error
//	if err == nil {
//		round = lastZj.Round
//	}
//	if err != nil && err != gorm.ErrRecordNotFound {
//		log.Println(err)
//	}
//	round = round + 1
//
//	var pdArr []model.Wait
//	err = p.db.Where(
//		"dest_product_id = ? AND round = ?",
//		DestProductId,
//		round,
//	).Order("id asc").Find(&pdArr).Error
//	if err == gorm.ErrRecordNotFound || len(pdArr) == 0 {
//		return
//	}
//
//	j := 0
//	for i := 0; i < len(pdArr) / 4; i++ {
//		for j = i * 4; j < i * 4 + 4; j++ {
//			pd := pdArr[j]
//			z := &model.Zj{
//				Round: pd.Round,
//				Group: pd.Group,
//				Position: pd.Position,
//				OrderId: pd.OrderId,
//				DestProductId: DestProductId,
//				Index: pd.Index,
//				Status: model.ZjStatusWait,
//			}
//
//			p.db.Create(&z)
//		}
//	}
//	for ; j < len(pdArr); j++ {
//		pd := pdArr[j]
//		err := p.db.Model(&model.Wait{}).Where("id = ?", pd.Id).Update("status", model.WaitStatusMiss).Error
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}
//
//
//func (p *PinTuan) kj(DestProductId, win int) {
//	round := 1
//	var lastZj model.Zj
//	err := p.db.Where("dest_product_id = ?", DestProductId).Find(&lastZj).Error
//	if err == nil {
//		round = lastZj.Round
//	}
//	if err != nil && err != gorm.ErrRecordNotFound {
//		log.Println(err)
//	}
//	round = round + 1
//
//	var pdArr []model.Wait
//	err = p.db.Where(
//		"dest_product_id = ? AND round = ?",
//		DestProductId,
//		round,
//	).Order("id asc").Find(&pdArr).Error
//	if err == gorm.ErrRecordNotFound || len(pdArr) == 0 {
//		return
//	}
//
//	j := 0
//	for i := 0; i < len(pdArr) / 4; i++ {
//		for j = i * 4; j < i * 4 + 4; j++ {
//			pd := pdArr[j]
//			z := &model.Zj{
//				Round: pd.Round,
//				Group: pd.Group,
//				Position: pd.Position,
//				OrderId: pd.OrderId,
//				DestProductId: DestProductId,
//				Index: pd.Index,
//			}
//			if j == win {
//				z.Status = model.ZjStatusWin
//			}else{
//				z.Status = model.ZjStatusLost
//			}
//			p.db.Create(&z)
//		}
//	}
//	for ; j < len(pdArr); j++ {
//		pd := pdArr[j]
//		err := p.db.Model(&model.Wait{}).Where("id = ?", pd.Id).Update("status", model.WaitStatusMiss).Error
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}