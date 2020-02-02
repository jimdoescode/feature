Go Feature Flags
================

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

It's important that the `GetGroupIdentifier` method should return byte slices that are unique to each member. In the example, the group identifier is the id field given to each User in the database. This works great because it shouldn't change and is unique to each user.

Once you've satified the interface you can use the `EnabledFor` method.
```go
flag := feature.NewFlag("my new great feature", 0.25)
user := &User{123, false, ...} // This would be fetched from a db or something
if flag.EnabledFor(user) {
	// Do feature
} else {
	// Don't do feature
}
```

In our example the feature is enabled for 25% of the users. Those users within the 25% threshold will **always** have the feature enabled. Increasing the threshold percentage will increase the number of users who can see the feature, decreasing will reduce the number of users who can see the feature. Decreasing the threshold to 10% will mean that the 15% of users who saw the feature when the threshold was at 25% will now *not* see the feature.

Certain groups might need to always have a feature flag enabled. This can be done by returning true for the `AlwaysEnabled` method of the `feature.Group` interface. In our example if a user has the isAdmin flag set to true then that user will have access to all features.

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

Additional Documentation
------------------------

GoDoc: https://godoc.org/github.com/jimdoescode/feature

Credit where it's due
---------------------

The technique for reducing an identifier to a boolean comes from [this Etsy feature library](https://github.com/etsy/feature) which was written in PHP.

All the code in this repo is licensed under the [MIT license](https://opensource.org/licenses/MIT)
