API
===

## List wallets
Lists all wallet accounts in the system.

**Method**: `GET`

**URL**: `/wallets[?currency=USD]`

**Query String Params**:
Optional
- currency: string

### Success response
**Status Code**: `200`
```json
[
    {
        "id": "alice123",
        "balance": 800.0,
        "currency": "USD",
        "created_at": "2021-01-02T08:30:00Z",
        "updated_at": "2021-01-02T08:30:00Z"
    },
    {
        "id": "alice123",
        "balance": 800.0,
        "currency": "USD",
        "created_at": "2021-01-02T08:30:00Z",
        "updated_at": "2021-01-02T08:30:00Z"
    }
]
```

### Error response
**Status Code**: `400` | `500`
```json
{
    "error": "some short description"
}
```

## Get wallet
Get wallet account with ID.

**Method**: `GET`

**URL**: `/wallets/{id}`

**URL Params**:
Required
- id: string

### Success response
**Status Code**: `200`
```json
{
  "id": "sato-101",
  "balance": 5000,
  "currency": "JPY",
  "created_at": "2021-10-19T23:20:59.929457+08:00",
  "updated_at": "2021-10-19T23:20:59.929457+08:00"
}
```

### Error response
**Status Code**: `400` | `404` | `500`
```json
{
  "error": "sql: no rows in result set"
}
```

## Create wallet
Create a wallet account.

**Method**: `POST`

**URL**: `/wallets`

**Data Params**:
Required
- id: string
- init_amt: float
- currency: string

### Success response
**Status Code**: `200`
```json
{
  "id": "alice-123",
  "balance": 800.0,
  "currency": "USD",
  "created_at": "2021-10-19T23:20:59.929457+08:00",
  "updated_at": "2021-10-19T23:20:59.929457+08:00"
}
```

### Error response
**Status Code**: `400` | `500`
```json
{
  "error": "invalid currency code (ISO 4217)"
}
```

## Get wallet payments
List incoming and outgoing transfers to wallet account.

**Method**: `GET`

**URL**: `/wallets/{id}/payments`

**URL Params**:
Required
- id: string

### Success response
**Status Code**: `200`
```json
{
  "id": "alice-123",
  "balance": 800.0,
  "currency": "USD",
  "created_at": "2021-10-19T23:20:59.929457+08:00",
  "updated_at": "2021-10-19T23:20:59.929457+08:00"
}
```

### Error response
**Status Code**: `400` | `500`
```json
{
  "error": "invalid currency code (ISO 4217)"
}
```

## Create payment
Transfer from a wallet account to another of the same currency.
Payment should fail if balance has less than requested amount.

**Method**: `POST`

**URL**: `/wallets/{id}/payments`

**URL Params**:
Required
- id: string

**Data Params**:
Required
- to_account: string
- amount: float64

### Success response
**Status Code**: `200`
```json
{
  "account": "bob-456",
  "to_account": "alice-123",
  "amount": 50,
  "direction": 2,
  "created_at": "0001-01-01T00:00:00Z"
}
```

### Error response
**Status Code**: `400` | `500`
```json
{
  "error": "existing balance less than requested transfer amount"
}
```

## List all transfers
List all transfers. 
The endpoint can be passed a `from` and/or a `to` (wallet account IDs) which work together as a `where... OR` query.
It is up to the API client (e.g. web, mobile) to filter on their end for cases where the user explicitly wants a `where... AND`.

**Method**: `GET`

**URL**: `/transfers[?currency=JPY][&from=alice-123][&to=bob-456]`

**Query String Params**:
Optional
- currency: string
- from: string
- to: string

### Success response
**Status Code**: `200`
```json
[
  {
    "id": 1,
    "from": "alice-123",
    "to": "nil-000",
    "currency": "USD",
    "amount": 50,
    "created_at": "2021-10-20T07:28:19.576098+08:00"
  },
  {
    "id": 2,
    "from": "alice-123",
    "to": "bob-456",
    "currency": "USD",
    "amount": 100,
    "created_at": "2021-10-20T07:30:35.882997+08:00"
  },
  {
    "id": 3,
    "from": "bob-456",
    "to": "alice-123",
    "currency": "USD",
    "amount": 50,
    "created_at": "2021-10-20T07:31:10.542693+08:00"
  }
]
```

### Error response
**Status Code**: `400` | `500`
```json
{
  "error": "some malformed field in the universe"
}
```