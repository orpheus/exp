# system

This package contains our `services` or as they are called on code, `interactors`. It is the **business logic**, the use
cases of our domain logic in action; they are the services we provide. Known as `interfaces` in the `clean architecture`
design model.

Note here the difference between the **system** and **domain** layer. A `Skiller` is a domain object, is a reference to
a particular type of person captured in relation to our product. A `User`, on the other hand, is a system-related entity
used for system-related operations like authentication and authorization. A `Skiller` does not care about permissions,
just about the fields related to its domain. A `User` needs to know about its permissions, its role, and the relevant
domain fields.

I call this package `system` and not `services`, `application`, or `tech` because I feel it brings the focus to a more
correct concept of its value, that being of managing system-level objects and defining system-level behavior. It is the
technical aspect or layer atop the domain layer, and it does not care about who or what outside itself will implement
it. It just defines how the system will be used and how it will interact with the domain objects. 