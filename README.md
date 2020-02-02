Go Feature Flags
================

GoDoc: (https://godoc.org/github.com/jimdoescode/feature)

This is a simple feature flagging library. It allows you to enable a feature for a certain percentage of a group or a population.

Create a new feature flag by using the `feature.NewFlag` method. Provide the feature name and the threshold for how often the feature should be enabled. 
```go
flag := feature.NewFlag("my new great feature", 0.25)
```

How's it work?
--------------

Feature flagging works by taking the group's identifier combining it with the name of the flag and reducing that hash to a value between 0 and 1. If that value is less than the threshold set for the flag then the feature is enabled. Otherwise it's disabled. This allows you to execute experiments on certain groups or cautiously roll out a new feature.

Groups
------

The groups interface is what will let you consistently bucket a feature.

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
user := &User{123, false, ...} // This would be fetched from a db or something
if flag.EnabledFor(user) {
	// Do feature
} else {
	// Don't do feature
}
```

In our example the feature is enabled for 25% of the users. Those users within the threshold will **always** have the feature enabled unless the sample size is lowered by changing the flag's threshold percentage. Increasing the threshold percentage will enable the feature flag for new users that fall into the larger sample size.

Certain groups might need to always have a feature flag enabled. This can be done by returning true for the `AlwaysEnabled` method of the `feature.Group` interface. In our example if a user has the isAdmin flag set to true then all feature flags will be enabled for that user.
```go
admin := &User{123, true, ...} // This would be fetched from a db or something
if flag.EnabledFor(admin) {
	// Do feature
} else {
	// Don't do feature
}
```

Random flagging
---------------

If you don't need consistent bucketing then you can use the the `Enabled` function. This function will randomly return that a feature is enabled but the number of enabled vs disabled results will converge on the flag's threshold. In our case that's 25%. Which means that if we called `Enabled` 100 times then ~25 of those calls would return true.
```go
if flag.Enabled() {
	// Do feature
} else {
	// Don't do feature
}
```

Credit where it's due
---------------------

The technique for reducing an identifier to a boolean comes from [this Etsy feature library](https://github.com/etsy/feature) which was written in PHP.

All the code in this repo is licensed under the [MIT license](https://opensource.org/licenses/MIT)
