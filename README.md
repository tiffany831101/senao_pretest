# Installation
`docker-compose up -d --build`

# API Documentation

## Create User Account

### Request

- **URL:** `POST /account`
- **Description:** Create a new user account with the provided username and password.

#### Request Body

```json
{
    "username": "string",
    "password": "string"
}
```

- `username` (string, required, min=3, max=32): The username for the new user account.
- `password` (string, required, min=8, max=32): The password for the new user account.

### Response

#### Successful Response

```json
{
    "success": true
}
```

- `success` (boolean): Indicates whether the account creation was successful.

#### Error Responses

- **400 Bad Request:**
  - Invalid JSON payload:
    ```json
    {
        "success": false,
        "reason": "Invalid JSON payload"
    }
    ```
  - Password requirements not met:
    ```json
    {
        "success": false,
        "reason": "Password must be 8-32 characters long and include at least 1 uppercase letter, 1 lowercase letter, and 1 number."
    }
    ```
  - Username already exists:
    ```json
    {
        "success": false,
        "reason": "Username already exists"
    }
    ```

## Verify Account and Password

### Request

- **URL:** `POST /account/:username/validate`
- **Description:** Verify the provided username and password.

#### Request Body

```json
{
    "username":"string",
    "password": "string"
}
```
- `username` (string, required) 
- `password` (string, required, min=8, max=32): The password to be verified.

### Response

#### Successful Response

```json
{
    "success": true
}
```

- `success` (boolean): Indicates whether the username and password verification was successful.

#### Error Responses

- **400 Bad Request:**
  - Invalid JSON payload:
    ```json
    {
        "success": false,
        "reason": "Invalid JSON payload"
    }
    ```
- **401 Unauthorized:**
  - Username not found:
    ```json
    {
        "success": false,
        "reason": "Username not found"
    }
    ```
  - Incorrect password:
    ```json
    {
        "success": false,
        "reason": "Incorrect password"
    }
    ```
- **429 Too Many Requests:**
  - Rate limit exceeded:
    ```json
    {
        "success": false,
        "reason": "Rate limit exceeded. Try again later."
    }
    ```
- **500 Failed to reset attempts:**
  - Failed to reset the rate limiter in Redis
  ```json
    {
        "success": false,
        "reason": "Failed to reset attempts."
    }
  ```
