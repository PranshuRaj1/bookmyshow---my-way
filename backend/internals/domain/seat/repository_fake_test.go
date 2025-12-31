package seat

import (
	"errors"
	"sync"
	"time"
)

type fakeRepository struct {
	mu    sync.Mutex
	seats map[int64]ShowSeat // key = showSeatID
}

func newFakeRepository(seats []ShowSeat) *fakeRepository {
	m := make(map[int64]ShowSeat)
	for _, s := range seats {
		m[s.ShowSeatID] = s
	}
	return &fakeRepository{seats: m}
}

func (f *fakeRepository) ListSeatsForShow(showID int64, now time.Time) ([]ShowSeat, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var result []ShowSeat
	for _, s := range f.seats {
		if s.ShowID != showID {
			continue
		}

		// Treat expired locks as unlocked
		if s.Status == StatusLocked && s.LockExpiresAt != nil && now.After(*s.LockExpiresAt) {
			s.Status = StatusUnlocked
			s.LockedAt = nil
			s.LockExpiresAt = nil
			s.BookingID = nil
		}

		result = append(result, s)
	}
	return result, nil
}

func (f *fakeRepository) LockSeats(
	showID int64,
	seatIDs []int64,
	bookingID int64,
	lockUntil time.Time,
) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()

	// First pass: validate
	for _, id := range seatIDs {
		seat, ok := f.seats[id]
		if !ok || seat.ShowID != showID {
			return errors.New("seat not found")
		}

		if seat.Status == StatusBooked {
			return errors.New("seat already booked")
		}

		if seat.Status == StatusLocked && seat.LockExpiresAt != nil && now.Before(*seat.LockExpiresAt) {
			return errors.New("seat locked by another booking")
		}
	}

	// Second pass: apply
	// := time.Now()
	for _, id := range seatIDs {
		seat := f.seats[id]
		seat.Status = StatusLocked
		seat.LockedAt = &now
		seat.LockExpiresAt = &lockUntil
		seat.BookingID = &bookingID
		f.seats[id] = seat
	}

	return nil
}

func (f *fakeRepository) BookLockedSeats(showID int64, bookingID int64, now time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, seat := range f.seats {
		if seat.ShowID != showID {
			continue
		}

		if seat.BookingID == nil || *seat.BookingID != bookingID {
			continue
		}

		if seat.Status == StatusBooked {
			continue
		}

		// Must be locked and not expired
		if seat.Status != StatusLocked {
			return errors.New("seat not locked")
		}

		if seat.LockExpiresAt != nil && now.After(*seat.LockExpiresAt) {
			return errors.New("lock expired")
		}
	}

	for id, seat := range f.seats {
		if seat.BookingID != nil && *seat.BookingID == bookingID {
			seat.Status = StatusBooked
			f.seats[id] = seat
		}
	}

	return nil
}

func (f *fakeRepository) ReleaseExpiredLocks(now time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for id, seat := range f.seats {
		if seat.Status == StatusLocked && seat.LockExpiresAt != nil && now.After(*seat.LockExpiresAt) {
			seat.Status = StatusUnlocked
			seat.LockedAt = nil
			seat.LockExpiresAt = nil
			seat.BookingID = nil
			f.seats[id] = seat
		}
	}
	return nil
}
