# LUSH Core Authentication
This package is used to deal with authenticated requests and responses within the LUSH infrastructure.

## Policies
Policies should be used to structure access control inside a project's domain logic.

### Roles Policy
Given that you have an action that can only be performed by administrators, you can set up a `RolePolicy` with the `"admin"` role. This will permit the consumer access only if they possess one or more of the roles specified in the policy.

```go
policy := lushauth.RolePolicy{"admin", "staff"}
policy.Permit(consumer)
```

### Grants Policy
For more granular access control you can require grants to be set in the consumer. Given that you've got an action to delete a user, you can set up a `GrantPolicy` with the `"users.delete"` grant. This will permit the consumer access only if they possess one or more grants specified in the policy.

```go
policy := lushauth.GrantPolicy{"users.delete"}
policy.Permit(consumer)
```

### User Policy
Sometimes actions are bound to only be performed by a very specific user. Given that you have an action for a user to update their own profile, you can set up a `UserPolicy` with the UUID of the user. This will permit the consumer access only if their UUID match any of the UUIDs specified in the policy.

```go
policy := lushauth.UserPolicy{
    UserID, // UserID: "5d4b32f9-5954-41c3-a470-7d76317635a7"
}
policy.Permit(consumer)
```

### Market Policy
For certain things you need to allow access per market and roles in those markets. Given that you've got an action to set a price for a product in the British market, you can set up a `MarketPolicy` with the `"gb"` market id and the `"digital_manager"` market role. This will permit the consumer access only if they belong to the given market and that they possess one or more of the roles for the given market.

```go
policy := lushauth.MarketPolicy{
    ID: "gb",
    Roles: []string{
        "digital_manager",
    },
}
policy.Permit(consumer)
```

### Any Policy
Sometimes you might have different access criteria for a given action. Given that you have an action to update a page for the Swedish market which can only be done by a digital manager within that market, _OR_ by a global administrator, you can set up multiple policies. This will permit the consumer access only if they're permitted access within any of the policies.

```go
policy := lushauth.AnyPolicy{
    lushauth.MarketPolicy{
        ID: "gb",
        Roles: []string{
            "digital_manager",
        },
    },
    lushauth.RolePolicy{"admin"},
}
policy.Permit(consumer)
```

### All Policy
Sometimes you need more than one criteria for access criteria for a given action. Given that you have an action for going into maintenance mode in the Netherlands where you want to ensure only a specific intersection of people with the global `"admin"` role and the `"digital_manager"` role within the market. This will permit the consumer access only if they're permitted by all policies.

```go
policy := lushauth.AllPolicy{
    lushauth.MarketPolicy{
        ID: "nl",
        Roles: []string{
            "digital_manager",
        },
    },
    lushauth.RolePolicy{"admin"},
}
policy.Permit(consumer)
```