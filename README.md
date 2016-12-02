Go Feature Flags
================

This is a simple feature flagging library. It allows you to enable a feature for a certain percentage of a group or a population.

Create a new feature flag by using the `feature.NewFlag` method. Provide the feature name and the threshold for how often the feature should be *enabled*. 
```go
flag := feature.NewFlag("my new great feature", 0.25)
```

To randomly enable a feature flag use the `Enabled` method. This flag will be enabled or not each time it's called. Over a number of calls the number of enabled vs disabled results will converge on the flag threshold. In our case that's 25%.
```go
if flag.Enabled() {
	// Do feature
} else {
	// Don't do feature
}
```

Groups
------

If you'd like to consistently enable a feature flag, meaning it will always be enabled for certain members that are within the flag threshold, you must satisfy the `feature.Group` interface.
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
	return u.isAdmin //Admins should have all feature flags enabled.
}
```

It's important that the `GetGroupIdentifier` method should return consistent byte slices that are unique to each member. In the example the group identifier is the id field given to each User. This works great because it won't change for a user but is also unique to that user.

Once you've satified the interface you can use the `EnabledFor` method.
```go
user := &User{123, false, ...}
if flag.EnabledFor(user) {
	// Do feature
} else {
	// Don't do feature
}
```

In our example the feature is enabled for 25% of the users. Those users within the threshold will always have the feature enabled unless the sample size is lowered by changing the flag's threshold percentage. Increasing the threshold percentage will enable the feature flag for new users that fall into the larger sample size.

Certain groups might need to always have a feature flag enabled. This can be done by returning true for the `AlwaysEnabled` method of the `feature.Group` interface. In our example if a user has the isAdmin flag set to true then all feature flags will be enabled for that user.
```go
admin := &User{123, true, ...}
if flag.EnabledFor(admin) {
	// Do feature
} else {
	// Don't do feature
}
```

How's it work?
--------------

Feature flagging works by taking the group identifier and reducing it to a value between 0 and 1. If that value is less than the threshold set for a flag then the feature is enabled. Otherwise it's disabled. This allows you to execute experiments on certain groups or cautiously roll out a new feature.
