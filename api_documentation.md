# Nanny for API Documentation
About
This is our pet sitting service API that we built for the final project. It connects pet owners with sitters in Almaty.

## Base URL

```
http://localhost:8080
```

## Authentication
For protected endpoints, you need to include JWT token in header:

```
Authorization: Bearer <your_token>
```

## Authentication Endpoints

### Register as Owner
`POST /api/auth/register/owner`

No auth needed

**Request:**
```json
{
  "full_name": "Anara Armankyzy",
  "email": "anara.arman@kbtu.kz",
  "phone": "+77770296982",
  "password": "SecurePass123"
}
```

**Response (200):**
```json
{
  "message": "Owner registered successfully"
}
```
### Register as Sitter
`POST /api/auth/register/sitter`

No auth needed

**Request:**
```json
{
  "full_name": "Anara Armanik",
  "email": "anara.armanik@kbtu.kz",
  "phone": "+77770296982",
  "password": "SecurePass123",
  "experience_years": 5,
  "certificates": "Pet Care Certificate 2022",
  "preferences": "Prefer dogs and cats",
  "location": "Almaty, Kazakhstan"
}
```

**Response (200):**
```json{
  "message": "Sitter registered successfully, awaiting admin approval"
}
```

*Note: New sitters need admin approval before they can create services*

### Login
`POST /api/auth/login`

**Request:**
```json
{
  "email": "anara.arman@kbtu.kz",
  "password": "SecurePass123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": 1,
    "full_name": "Anara Armankyzy",
    "email": "anara.arman@kbtu.kz",
    "phone": "+77770296982",
    "role": "owner"
  }
}
```

Token is valid for 72 hours

## Pets
### Create Pet
`POST /api/pets`

Needs auth (Owner only)

**Request:**
```json
{
  "name": "Buddy",
  "type": "dog",
  "age": 3,
  "notes": "Friendly golden retriever"
}
```

Valid types: cat, dog, rodent

**Response (201):**
```json
{
  "pet_id": 1,
  "owner_id": 1,
  "name": "Buddy",
  "type": "dog",
  "age": 3,
  "notes": "Friendly golden retriever"
}
```

### Get Pet
`GET /api/pets/{id}`
Public endpoint

**Response (200):**
```json
{
  "pet_id": 1,
  "owner_id": 1,
  "name": "Buddy",
  "type": "dog",
  "age": 3,
  "notes": "Friendly golden retriever"
}
```

### Update Pet
`PUT /api/pets/{id}`

Needs auth (Owner only, must be your pet)

**Request:**
```json
{
  "name": "Buddy Jr.",
  "type": "dog",
  "age": 4,
  "notes": "Now 4 years old"
}
```

### Delete Pet
**DELETE** `/api/pets/{id}`

Needs auth (Owner only)

### Get Owner's Pets
**GET** `/api/owners/{owner_id}/pets`

Public endpoint

Returns array of all pets for that owner

---

## Services

### Search Services
**GET** `/api/services/search`

Public endpoint

Query params (all optional):
- `type` - walking, boarding, or home-care
- `location` - filter by location
- `min_price` - minimum price
- `max_price` - maximum price

Example:
`/api/services/search?type=walking&location=Almaty&max_price=3000`

Response (200):
```json[
  {
    "service_id": 1,
    "sitter_id": 2,
    "sitter_name": "Jane Smith",
    "location": "Almaty",
    "type": "walking",
    "price_per_hour": 2500.00,
    "description": "1-hour walk with your dog",
    "rating": 4.8,
    "total_reviews": 25
  }
]
```

## Get Service Details
`GET /api/services/{id}`

Public endpoint

## Get Sitter's Services
GET `/api/sitters/{sitter_id}/services`

Public endpoint

## Create Service
POST `/api/services`

Needs auth (Sitter only, must be approved)

**Request:**
```json
{
  "type": "boarding",
  "price_per_hour": 7000.00,
  "description": "Pet stays at my place overnight"
}
```

## Update Service
PUT `/api/services/{id}`
Needs auth (Sitter only, must be your service)

## Delete Service
DELETE `/api/services/{id}`
Needs auth (Sitter only)

## Bookings
Create Booking
POST `/api/bookings`
Needs auth (Owner only)

**Request:**
```json
{
  "sitter_id": 2,
  "pet_id": 1,
  "service_id": 1,
  "start_time": "2025-12-20T10:00:00Z",
  "end_time": "2025-12-20T11:00:00Z"
}
```

**Response (201):**
```json{
  "booking_id": 1,
  "status": "pending",
  "message": "Booking created successfully"
}
```

Booking flow: pending -> confirmed -> completed (or cancelled anytime)

### Get Booking

**GET**  `/api/bookings/{id}`
Public endpoint

**Confirm Booking**

**POST** `/api/bookings/{id}/confirm`
Needs auth (Sitter only)
Only sitter assigned to booking can confirm

**Cancel Booking**

**POST** `/api/bookings/{id}/cancel`
Needs auth (Owner or Sitter)

**Complete Booking**

**POST** `/api/bookings/{id}/complete`
Needs auth (Sitter only)
Can only complete after end_time has passed

**Get Owner's Bookings**

GET `/api/owners/{owner_id}/bookings`
Optional query param: `status` (pending/confirmed/cancelled/completed)

**Get Sitter's Bookings**

GET `/api/sitters/{sitter_id}/bookings`
Same as owner's bookings

## Reviews
### Create Review

POST `/api/reviews`

Needs auth (Owner only)

Request:
```json
{
  "booking_id": 1,
  "rating": 5,
  "comment": "Excellent service!"
}
```
**Requirements:**

- Rating must be 1-5
- Booking must be completed
- Can only review your own bookings
- One review per booking

## Get Review
GET `/api/reviews/{id}`

Public endpoint
### Update Review

PUT `/api/reviews/{id}`

Needs auth (Owner only, must be your review)

### Delete Review
DELETE /api/reviews/{id}

Needs auth (Owner only)

### Get Sitter's Reviews
GET `/api/sitters/{sitter_id}/reviews`
Public endpoint returns all reviews for that sitter

### Get Sitter Rating
GET `/api/sitters/{sitter_id}/rating`

Response:
```json
{
  "sitter_id": 2,
  "average_rating": 4.75,
  "total_reviews": 24
}
```

### Get Booking Review
GET `/api/bookings/{booking_id}/review`
Returns the review for specific booking (if exists)

## Admin Endpoints
All admin endpoints need admin role
### Get Pending Sitters
GET `/api/admin/sitters/pending`
Returns list of sitters waiting for approval
### Approve Sitter
POST `/api/admin/sitters/{sitter_id}/approve`
Changes sitter status to "approved"
### Reject Sitter
POST `/api/admin/sitters/{sitter_id}/reject`
Changes sitter status to "rejected"
### Get Sitter Details
GET `/api/admin/sitters/{sitter_id}`
Returns detailed info about sitter including stats
### Get All Users
GET `/api/admin/users`

## Query params:

- role - filter by owner/sitter/admin
- limit - how many results (default 50)
- offset - for pagination

**Get User Details**

GET `/api/admin/users/{user_id}`

**Delete User**
DELETE `/api/admin/users/{user_id}`


Warning: This deletes everything related to user (pets, bookings, etc)

## **Error Responses**

All errors return JSON:
```json{
  "error": "Error message here"
}
```

### Common status codes:

200 - OK
201 - Created
400 - Bad request (invalid data)
401 - Unauthorized (need to login)
403 - Forbidden (don't have permission)
404 - Not found
500 - Server error


## Testing

Quick test with curl:

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register/owner \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Test","email":"test@test.com","phone":"+77001234567","password":"test123"}'
```

# Login and save token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}' | jq -r '.token')


# Create pet
```
curl -X POST http://localhost:8080/api/pets \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Max","type":"dog","age":2,"notes":"Good boy"}'
```

## Using Postman:

1. Create collection
2. Add environment variables: base_url and token
3. Use {{base_url}} and {{token}} in requests


## Notes

1. Tokens expire after 3 days
2. Passwords are hashed with bcrypt
3. Pet types are limited to cat, dog, rodent
4. Service types are walking, boarding, home-care
5. Start time must be in future when creating booking
6. Can't review until booking is completed


### If you have any questions please meet the team: Anara, Nuray, Arukhan, Sabukhi
