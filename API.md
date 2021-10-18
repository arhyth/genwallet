API
===

## List wallets
Lists all wallet accounts in the system

**Method**: `GET`

**URL**: `/wallets`

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
**Status Code**: `400 | 500`
```json
{
    "error": "some short description"
}
```