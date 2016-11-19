Go Feature Flags
================

This is a simple feature flagging library. It allows you to enable a feature for a certain percentage of a group or a population.

Create a new feature flag by using the `feature.NewFlag` method. Provide the feature name and the percentage of times the feature will be *enabled*. 
```go
flag := feature.NewFlag("feature", 0.75)
```

To randomly enable a feature use the `Enabled` method. This feature will be enabled or not each time it's called. Over a number of calls the number of enabled vs disabled features will converge on the flag percentage.
```go
if flag.Enabled() {
	// Do feature
} else {
	// Don't do feature
}
```

Groups
------

If you'd like to consistently apply a feature, meaning it will always be enabled for group members that are within the flag threshold, you must satisfy the `feature.Group` interface.
```go
type User struct {
	id uint
	isAdmin bool
	...
}

func (u *User) GetGroupIdentifier() []byte {
	buf := make([]byte, 8)
	binary.PutUvarint(buf, u.id)
	return buf
}

func (u *User) AlwaysEnabled() bool {
	return u.isAdmin
}
```

It's important that the `GetGroupIdentifier` method should return consistent byte slices for specific members. In our example the group identifier is the id field given to each User entry.

Once you've satified the interface you can use the `EnabledFor` method.
```go
user := &User{123, false, ...}
if flag.EnabledFor(user) {
	// Do feature
} else {
	// Don't do feature
}
```

In our example the feature is enabled for half of the users. Those users will always have the feature enabled unless the flag percentage is lowered. 

Certain groups might need to always have a feature enabled. This can be done by returning true for the `AlwaysEnabled` method of the `feature.Group` interface.
```go
user := &User{123, true, ...}
if flag.EnabledFor(user) {
	// Do feature
} else {
	// Don't do feature
}
```

How's it work?
--------------

Feature flagging works by taking the group identifier and reducing it to a value between 0 and 1. If that value is less than the percentage set for a flag then the feature is enabled. Otherwise it's disabled. This allows you to execute experiments on certain groups or slowly ramp up a feature.
