# Database Schema Documentation

This directory contains the technical details of the database schema.

## Entity Relationship Diagram

The following diagram illustrates the relationships between the core entities: **Cinemas**, **Screens**, **Movies**, **Shows**, **Seats**, **Bookings**, and **Payments**.

```mermaid
erDiagram
    cinemas ||--o{ screens : "contains"
    screens ||--o{ shows : "hosts"
    screens ||--o{ seats : "has"
    movies ||--o{ shows : "features"
    shows ||--o{ show_seat : "has seats"
    seats ||--o{ show_seat : "instances"
    users ||--o{ bookings : "makes"
    shows ||--o{ bookings : "booked for"
    bookings ||--o{ show_seat : "reserves"
    bookings ||--|| payments : "paid by"

    cinemas {
        integer cinema_id PK
        varchar name
        varchar city
    }

    screens {
        integer screen_id PK
        integer cinema_id FK
        varchar name
    }

    shows {
        integer show_id PK
        integer movie_id FK
        integer screen_id FK
        timestamp start_time
    }

    users {
        integer id PK
        varchar username
        timestamp created_at
    }

    movies {
        integer movie_id PK
        varchar movie_name
        varchar movie_description
        integer movie_duration
    }

    seats {
        integer seat_id PK
        varchar row
        integer column
        integer screen_id FK
        tier_variant tier_variant
    }

    show_seat {
        integer show_seat_id PK
        integer seat_id FK
        integer show_id FK
        Status status
        timestamp locked_at
        timestamp lock_expires_at
        decimal price
        integer booking_id FK
    }

    bookings {
        integer booking_id PK
        integer user_id FK
        integer show_id FK
        BookingStatus status
        timestamp created_at
    }

    payments {
        integer payment_id PK
        integer booking_id FK
        decimal amount
        payment_status status
        varchar provider_reference
        timestamp created_at
    }
```

## Table Descriptions

### Key Entities

- **cinemas**: Represents physical theater locations.
- **screens**: Individual halls within a cinema.
- **movies**: Film metadata.
- **shows**: Specific instances of a movie playing on a screen at a time.
- **seats**: Physical seats in a screen (static configuration).

### Booking Flow Entities

- **users**: Registered customers.
- **show_seat**: The dynamic inventory table. Links a `seat` to a `show` and tracks its availability status (`Unlocked`, `Locked`, `Booked`).
- **bookings**: Represents a user's intent to reserve seats.
- **payments**: Records financial transactions for bookings.

## Enums

- **Status**: `Unlocked`, `Locked`, `Booked` (Used in `show_seat`).
- **BookingStatus**: `Created`, `Confirmed`, `Expired`, `Cancelled`.
- **payment_status**: `Pending`, `Success`, `Failed`.
- **tier_variant**: `VIP`, `Premium`, `Regular`.

## DBML Definition

The raw DBML definition is available in [schema.dbml](./schema.dbml). You can copy the content of that file and paste it into [dbdiagram.io](https://dbdiagram.io) for an interactive view.
