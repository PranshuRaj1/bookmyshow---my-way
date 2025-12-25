# Database Design Showcase & Project Roadmap

> I learned about database design from Caleb's playlist and took this as a project to showcase my skills. I will be creating a BookMyShow-type application using **Next.js** and **Golang** (which I will learn and implement).

Caleb's playlist: https://www.youtube.com/watch?v=h0j0QN2b57M&list=PL_c9BZzLwBRK0Pc28IdvPQizD2mJlgoID

This project demonstrates a robust database design for a cinema booking system. Rather than just a table overview, this documentation explains the **feature map**—what real-world capabilities this design enables and why they matter.

---

## Design Philosophy

This design implements a **correct, concurrent, time-aware seat booking system** that can survive real-world failures: retries, crashes, races, and payment inconsistencies. It prioritizes **inventory truth**, enforces **atomic multi-seat booking**, separates **intent from money**, and prevents deadlocks by construction. What you’ve built is not just a schema — it’s a **set of enforceable invariants** that make the system trustworthy under stress.

That’s what real users care about, even if they never see the tables.

---

## 1. Real-time seat availability (correct under load)

### Feature

Users can see which seats are:

* available
* temporarily locked
* fully booked

### Why your design supports this

* Physical seats (`Seats`) never change
* Availability lives in `show_seat`
* `status + lock_expires_at` allows the backend to treat **expired locks as free**

This means:

* no cron jobs guessing state
* no stale “ghost locks”
* availability is always computed from truth

### Real-world effect

On a Friday night spike, users don’t see seats that are already half-booked or stuck forever.

---

## 2. Time-bound seat locking (anti-double-booking)

### Feature

When a user selects seats, those seats are **reserved for a limited time**.

### How it’s implemented

* `Locked` status
* `locked_at`
* `lock_expires_at`

Locks are:

* enforced per seat
* enforced per show
* enforced by time, not by hope

### Real-world effect

* Two users can’t race to grab the same seat
* Abandoned carts don’t poison inventory
* The system self-heals after crashes

---

## 3. Multi-seat atomic booking (all or nothing)

### Feature

A booking with multiple seats:

* either succeeds completely
* or fails completely

No partial bookings.

### How your design enables this

* Seats are updated **inside one transaction**
* Each seat transitions conditionally: `Locked → Booked`
* If *any* seat fails, the whole transaction rolls back

### Real-world effect

Users never end up with:

* “you booked 2 seats out of 3”
* inconsistent pricing
* confusing receipts

This is critical for trust.

---

## 4. Concurrency safety under retries and races

### Feature

The system survives:

* duplicate API calls
* payment retries
* webhook duplication
* slow gateways

### How

* Conditional state transitions decide winners
* “0 rows affected” = safe failure
* Booking confirmation is **idempotent**

You don’t rely on:

* frontend timing
* redirects
* isolation levels alone

### Real-world effect

Even if the payment gateway misbehaves, your inventory does not.

---

## 5. Grace window + retry support (good UX without lying)

### Feature

Users get:

* a grace window if payment is initiated
* retry attempts without losing seats immediately

### Why it works

* Lock expiration is explicit and adjustable
* Eligibility is time-based
* Final authority is still the seat state

You improve UX **without weakening correctness**.

### Real-world effect

Fewer refunds, fewer angry users, same safety guarantees.

---

## 6. Clean separation of concerns (maintainability)

### Feature

Each concept has one responsibility:

* **Seats** → physical layout
* **show_seat** → inventory + truth
* **Booking** → user intent + lifecycle
* **Payment** → money + gateway state

### Why this matters

* Bugs are localized
* Recovery logic is clear
* Future changes don’t ripple unpredictably

### Real-world effect

When something breaks at 9 PM on release night, you know *where* to look.

---

## 7. Payment reconciliation and recovery

### Feature

The system can handle:

* payment success but booking failure
* delayed callbacks
* refunds

### How

* Payment is event-driven
* Booking has its own lifecycle (`Created`, `Confirmed`, `Expired`)
* Seat ownership is checked *after* payment

### Real-world effect

Money and inventory are reconciled safely instead of pretending they move together.

---

## 8. Deadlock prevention under high contention

### Feature

Multiple users booking overlapping seats don’t freeze the system.

### How

* Seats are processed in a canonical order
* All transactions acquire locks consistently
* Circular waits are mathematically impossible

### Real-world effect

Traffic spikes slow things down — they don’t crash the platform.

---

## 9. Deterministic, debuggable system behavior

### Feature

Given the same inputs and timing, the system behaves predictably.

### Why

* State transitions are explicit
* Time is modeled, not inferred
* Ownership is single-sourced

### Real-world effect

You can:

* replay incidents
* explain failures to users
* answer support tickets confidently

---

## 10. Scales logically before it scales technically

### Feature

This design scales **conceptually** before you add:

* Redis
* queues
* caches
* microservices

### Why this matters

Most systems fail from *incorrect logic*, not slow databases.

Your design ensures:

* correctness first
* optimization later



NOTE: This is my deduction as of now, and I will be updating this as I learn more.


