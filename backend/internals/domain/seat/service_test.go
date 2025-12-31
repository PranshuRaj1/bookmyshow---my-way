package seat

import (
	"testing"
	"time"
)

func TestLockSeats_AllOrNothing(t *testing.T) {
	now := time.Now()

	seats := []ShowSeat{
		{ShowSeatID: 1, ShowID: 1, Status: StatusUnlocked},
		{ShowSeatID: 2, ShowID: 1, Status: StatusUnlocked},
		{ShowSeatID: 3, ShowID: 1, Status: StatusUnlocked},
	}

	repo := newFakeRepository(seats)
	service := NewService(repo)

	err := service.LockSeatsForBooking(
		1,
		[]int64{1, 2, 3},
		101,
		10*time.Minute,
		now,
	)

	if err != nil {
		t.Fatalf("expected lock to succeed, got error: %v", err)
	}
}

func TestLockSeats_FailsIfOneSeatLocked(t *testing.T) {
	now := time.Now()
	lockUntil := now.Add(10 * time.Minute)

	seats := []ShowSeat{
		{ShowSeatID: 1, ShowID: 1, Status: StatusUnlocked},
		{
			ShowSeatID:    2,
			ShowID:        1,
			Status:        StatusLocked,
			LockExpiresAt: &lockUntil,
		},
	}

	repo := newFakeRepository(seats)
	service := NewService(repo)

	err := service.LockSeatsForBooking(
		1,
		[]int64{1, 2},
		101,
		10*time.Minute,
		now,
	)

	if err == nil {
		t.Fatalf("expected lock to fail due to locked seat")
	}
}

func TestConfirmBooking_Idempotent(t *testing.T) {
	now := time.Now()
	lockUntil := now.Add(10 * time.Minute)
	bookingID := int64(101)

	seats := []ShowSeat{
		{
			ShowSeatID:    1,
			ShowID:        1,
			Status:        StatusLocked,
			BookingID:     &bookingID,
			LockExpiresAt: &lockUntil,
		},
	}

	repo := newFakeRepository(seats)
	service := NewService(repo)

	err1 := service.ConfirmBookingSeats(1, bookingID, now)
	err2 := service.ConfirmBookingSeats(1, bookingID, now)

	if err1 != nil || err2 != nil {
		t.Fatalf("confirmation should be idempotent")
	}
}
