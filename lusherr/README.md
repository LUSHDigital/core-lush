# LUSH Core Errors
This package is used to streamline dealing with errors and error messages within the LUSH infrastructure. Using the errors provided by this package will collect a lot of useful debug information about where exactly the error occurred.

## Error Types
These are the standard error types that can be used within a project's domain logic to aid with debugging and error reporting to any of its API consumers.

### Internal Error
`InternalError` can be used to wrap any error e.g. Trying to generate a random UUID, but the generation failed.

```go
id, err := uuid.NewV4()
if err != nil {
    return NewInternalError(err)
}
```

### Unauthorized Error
`UnauthorizedError` should be used when an action is performed by a user that they don't have permission to do e.g. Someone tried to access something they were not allowed to according to a permission policy.

```go
if err := policy.Permit(consumer); err != nil {
    return NewUnauthorizedError(err)
}
```

### Validation Error
`ValidationError` should be used to detail what user generated information is incorrect and why e.g. Someone set the name field for a user to be empty, but the validation requires it to be present.

```go
type ProductRevision struct {
    plu string
}

func (r ProductRevision) validate() error {
    if plu == "" {
        return NewValidationError("product revision", "plu", fmt.Errorf("must be present"))
    }
}
```

### Database Query Error
`DatabaseQueryError` should be used to provide detail about a failed database query e.g. Trying to query the database, but the database rejects the query.

```go
const stmt = `SELECT * FROM user`
rows, err := qu.Query(stmt)
if err != nil {
    return nil, NewDatabaseQueryError(stmt, err)
}
```

### Not Found Error
`NotFoundError` should be used when an entity cannot be found e.g. Someone tries to retrieve a user, but the user for the given ID does not exist in the database.

```go
const stmt = `SELECT * FROM user WHERE id = $1`
rows, err := qu.Query(stmt, id)
if err != nil {
    switch err {
    case sql.ErrNoRows:
        return nil, NewNotFoundError("user", id, err)
    default:
        return nil, NewDatabaseQueryError(stmt, err)
    }
}
```

### Not Allowed Error
`NotAllowedError` should be used when an certain action is not allowed e.g. Someone tries to delete something, but the record has been marked as permanent.

```go
if product.permanent {
    return NewNotAllowedError(fmt.Errorf("not allowed to remove permanent products"))
}
```

## Locate
Errors produced with the `lusherr` package can be located to return the `runtime.Frame` of where it occurred.

```go
frame, found := lusherr.Locate(err)
if found {
    log.Println(err, frame)
} else {
    log.Println(err, "frame could not be found")
}
```

### Locator Interface
Implement the [`Locator`](#locator-interface) interface on an error type to return its caller frame to be used in conjunction with the `lusherr` package and associated tooling.

```go
type Locator interface {
    Locate() runtime.Frame
}
```

## Pin
Call `Pin` to wrap an error with information about the caller frame of where `Pin` was invoked for an error that does not already carry a caller frame. The resulting error can be cast to [`Locate`](#locator-interface).

```go
func UploadToBucket(w io.Writer) error {
    if err := upload(w); err != nil {
        return lusherr.Pin(err)
    }
}
```

## Pinner Interface
Implement the [`Pinner`](#pinner-interface) interface on an error type to prevent errors that already implements [`Locator`](#locator-interface) to be wrapped multiple times.

```go
type Pinner interface {
    Pin(runtime.Frame) error
}
```