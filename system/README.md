# system

This package contains our `services` or as they are called on code, `interactors`. It is the **business logic**, the use
cases of our domain logic in action; they are the services we provide. Known as `interfaces` in the `clean architecture`
design model.

Note here the difference between the **system** and **domain** layer. A `Skiller` is a domain object, is a reference to
a particular type of person captured in relation to our product. A `User`, on the other hand, is a system-related entity
used for system-related operations like authentication and authorization. A `Skiller` does not care about permissions,
just about the fields related to its domain. A `User` needs to know about its permissions, its role, and the domain.

I call this package `system` and not `services`, `application`, or `tech` because I feel it brings the focus to a more
correct conception, that being of managing system-level objects and defining system-level behavior. It is the technical
aspect or layer atop the domain layer, and it does not care about who or what outside itself will implement it. It just
defines how the system will be used and how it will interact with the domain objects.

### Defining Repository Interfaces

#### Single Data Source

Why do the System-Interactors not define their own Repository-Interfaces? We let API-Interactors define their own
System-Interfaces. If we wanted to build something else that needs Repository methods, it shouldn't need to go into the
System Layer to get it (but should it just define its own?). By having it at the domain layer we are saying that every
Repository-Implementation needs to define all the Domain-level Repository methods for it to work with our
System-Interactor.

##### tldr: the domain-bottom layer defines the interface. the system-top layer fulfills it.

```
// domain
type SkillRepo interface {
    FindAll() []Skill
    FindById(id string) Skill
}

// system
type SkillInteractor struct {
    SkillRepo domain.SkillRepo
}
```

This is backwards from the API-System interfaces where we allow the top-layer to define the interface it'll use and
allow the bottom-layer to fulfill it.

##### tldr: the api-top layer defines its interface. the system-bottom layer fulfills it.

```
// api
type SkillController struct {
    SkillInteractor SkillInteractorInterface
}

type SkillInteractorInterface interface {
    FindAll() []domain.Skill
}

// system
type SkillInteractor struct {
    SkillRepo domain.SkillRepo
}

func (s *SkillInteractor) FindAll() []Skill {
    return s.SkillRepo.FindAll()
}
```

What else would use our Repository interface? The way it is now, the System-Interactor uses the Repository-Interface and
allows us to pass in any Repository that fulfills that contract. So we can swap out Data Sources by creating new
API-level repository interactors that fit the domain repository interface.

##### tldr: with the repository interface in the domain-bottom layer, we can swap out individual data sources

```
// domain
type SkillRepo interface {
    FindAll() []Skill
}

// api
type AnyRepository struct {
    connection DB
}

func (a *AnyRepository) FindAll() []Skill {
    return a.connection.doSomething()
}

// system
type SkillInteractor struct {
   SkillRepo domain.SkillRepo
}

&SkillInteractor{
    SkillRepo AnyRepository
}
```

#### Multiple Data Sources with Multiple Interactors

But what if we wanted to add multiple data sources? Should a new Interactor use the existing repository interface? Say I
add `ElasticSearch` to our infrastructure. In our APIs I build out the Interactor, or define the service and services it
will provide, i.e. Search for all Skills. The Interactor will interface with ElasticSearch, get back data, convert that
data into domain objects and return. Say I have another SkillInteractor to handle this, an `ElasticSkillInteractor`.
What Repository interface is it implementing? If it implements `domain.SkillRepository`, then
the `ElasticSearchRepoInteractor` will need to implement every single method. If the `ElasticSkillInteractor`
defines the Repository-Interface it will be using instead of using the one defined in the domain layer, then we can keep
the `ElasticSearchRepoInteractor` to a minimum.

##### tldr: defining the repository interface (api) at the system level allows us to define only what we use

```
// api
type ElasticSearchRepoInteractor struct {
    ES ElasticSearch
}

func (e *ElasticSearchInteractor) FindAll() []Skill {
    return e.ES.findAllSkills()
}

// system
type ElasticSkillInteractor struct {
    // using this would mean our ElasticSearch API has to implement all of the methods defined
    DomainSkillRepository domain.SkillRepo 

    // defining a custom repo means we can only define what our ElasticSearch API is going to use
    ESSkillRepository CustomSkillRepo 
}

type CustomSkillRepo interface {
    FindAll() []Skill
}
```

#### Multiple Data Sources in One Interactor

What if I just want to swap out one of the System-Interactor's methods with a call to a different API? Continuing with
the above example, I want my `FindAll()` call to hit ElasticSearch, but keep all my other calls hitting Postgres. What
do I do? Well my System-Interactor can:

1. Hard code an ElasticSearch-Interactor into its struct
2. Define a new SkillSearch-Interface into its struct
3. Split out its Repository interface into smaller interfaces

##### tldr: (2) define new smaller interfaces to handle api feature sets without touching other interfaces

Let's look at each of them:

1. Hard coded struct. Here we break principle. We are having system-layer code depend on api-layer layer code. No go.

```
// system
type SkillInteractor struct {
    PGSkillRepo domain.SkillRepo
    ESSkillRepo api.ESInteractor
}
```

2. Define a new, smaller Repository-Interface that just defines what we're going to be using from our outside API. We
   can abstract even further on `SkillSearch` and create an interface for handling all of our search capabilities.

```
// system
type SkillInteractor struct {
    PGSkillRepo domain.SkillRepo
    ESSkillRepo SkillSearch
}

type SkillSearch interface {
    FindAll() []Skill
}
```

3. Breaking out the original interface. This doesn't accomplish much and is possibly dangerous if you have other
   implementations that are relying on the SkillRepo's `FindAll()` method that you just took out. Our original domain
   repository can keep all of its functions and our `SkillInteractor` will just call a different receiver value.

```
// domain
type SkillRepo interface {
    // This no longer requires the `FindAll()`
    FindById(id string) Skill
}

// system
type SkillInteractor struct {
    PGSkillRepo domain.SkillRepo
    ESSkillRepo ElasticSearchFindAllSkillsRepository
}

type ElasticSearchFindAllSkillsRepository interface {
    FindAll() []Skill
}
```

