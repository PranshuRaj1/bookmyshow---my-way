package seat

import "time"


type SeatStatus string

const (
	StatusUnlocked SeatStatus = "Unlocked"
	StatusLocked SeatStatus = "Locked"
	StatusBooked SeatStatus  = "Booked"
)

type ShowSeat struct {
	ShowSeatID int64
	ShowID	   int64
	SeatID     int64
	Status     SeatStatus
	LockedAt   *time.Time
	LockExpiresAt *time.Time
	BookingID *int64
	Price int64
}


// Repository defines what the seat domain needs from persistence
type Repository interface {
ListSeatsForShow(showID int64, now time.Time) ([]ShowSeat, error) 

// LockSeats attempts to lock the given seats for a booking.
	// It must:
	// - respect canonical ordering
	// - fail if any seat cannot be locked
	// - be atomic from the caller's perspective
	LockSeats(
		showID int64,
		seatIDs []int64,
		bookingID int64,
		lockUntil time.Time,
	)error

	// BookLockedSeats attempts to transition seats from Locked → Booked.
	// It must:
	// - only succeed if seats are still locked for this booking
	// - be safe under retries
	BookLockedSeats(
		showID int64,
		bookingID int64,
		now time.Time,
	) error

	// ReleaseExpiredLocks releases locks that have passed expiry.
	
	ReleaseExpiredLocks(now time.Time) error

}