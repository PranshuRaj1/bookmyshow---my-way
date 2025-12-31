package seat

import (
	"errors"
	"slices"
	"time"
)

var (
	ErrNoSeatsProvided      = errors.New("no seats provided")
	ErrDuplicateSeatIDs     = errors.New("duplicate seat ids provided")
	ErrSeatLockFailed       = errors.New("one or more seats could not be locked")
	ErrSeatConfirmationFail = errors.New("one or more seats could not be booked")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetSeatsForShow returns seat inventory for a show.
// Expired locks must be treated as unlocked by repository.
func (s *Service) GetSeatsForShow(showID int64, now time.Time) ([]ShowSeat, error) {
	return s.repo.ListSeatsForShow(showID, now)
}

// LockSeatsForBooking attempts to lock all given seats for a booking.
// Either all seats are locked or none are.
func (s *Service) LockSeatsForBooking(
	showID int64,
	seatIDs []int64,
	bookingID int64,
	lockDuration time.Duration,
	now time.Time,
) error {

	if len(seatIDs) == 0 {
		return ErrNoSeatsProvided
	}

	if hasDuplicates(seatIDs) {
		return ErrDuplicateSeatIDs
	}

	// Canonical ordering (deadlock prevention)
	slices.Sort(seatIDs)

	lockUntil := now.Add(lockDuration)

	err := s.repo.LockSeats(
		showID,
		seatIDs,
		bookingID,
		lockUntil,
	)

	if err != nil {
		return ErrSeatLockFailed
	}

	return nil
}

// ConfirmBookingSeats confirms all seats locked for a booking.
// Operation is idempotent.
func (s *Service) ConfirmBookingSeats(
	showID int64,
	bookingID int64,
	now time.Time,
) error {

	err := s.repo.BookLockedSeats(showID, bookingID, now)
	if err != nil {
		return ErrSeatConfirmationFail
	}

	return nil
}

// ReleaseExpiredLocks is a cleanup operation.
// Can be run periodically.
func (s *Service) ReleaseExpiredLocks(now time.Time) error {
	return s.repo.ReleaseExpiredLocks(now)
}

// ---- helpers ----

func hasDuplicates(ids []int64) bool {
	seen := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return true
		}
		seen[id] = struct{}{}
	}
	return false
}
