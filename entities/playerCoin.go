package entities

import (
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
	"time"
)

type PlayerCoin struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;"`
	PlayerID  string    `gorm:"type:varchar(64);not null;"`
	Amount    int64     `gorm:"not null;"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime;"`
}

func (p *PlayerCoin) ToPlayerCoinModel() *playerCoinModel.PlayerCoin {
	return &playerCoinModel.PlayerCoin{
		ID:        p.ID,
		PlayerID:  p.PlayerID,
		Amount:    p.Amount,
		CreatedAt: p.CreatedAt,
	}
}
