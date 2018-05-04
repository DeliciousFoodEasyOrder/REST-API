package models

// Seat Model
type Seat struct {
	SeatID     int    `xorm:"AUTOINCR" json:"seat_id"`
	Number     string `json:"number"`
	QRCodeURL  string `json:"qr_code_url"`
	MerchantID int    `json:"merchant_id"`
}

type SeatDataAccessObject struct{}

var SeatDAO *SeatDataAccessObject

// InsertOne inserts a seat to database
func (*SeatDataAccessObject) InsertOne(s *Seat) {
	_, err := orm.InsertOne(s)
	if err != nil {
		panic(err)
	}
}

// FindByNumber finds a seat by id
func (*SeatDataAccessObject) FindByID(seatID int) *Seat {
	var seat Seat
	has, err := orm.Table(seat).ID(seatID).Get(&seat)
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &seat
}

func (*SeatDataAccessObject) DeleteBySeatID(seatID int) {
	var seat Seat
	_, err := orm.Table(seat).ID(seatID).Delete(&seat)
	if err != nil {
		panic(err)
	}
}

func (*SeatDataAccessObject) UpdateOne(seat *Seat) {
	_, err := orm.Table(seat).ID(seat.SeatID).Update(seat)
	if err != nil {
		panic(err)
	}
}
