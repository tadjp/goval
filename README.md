## Summary

## Guide

### Example

#### GetAll

```go
type Member struct {
    Name string
}
type Team struct {
    Name    string
    Members []*Member
}
team := Team{
    Name: "TEAM-A",
    Members: []*Member{
        {
            Name: "Alice",
        },
        {
            Name: "Bob",
        },
    },
}
path, _ := goval.Parse("Name")
fmt.Println(goval.GetAll[string](&team, path)) // [TEAM-A]
path, _ := goval.Parse("Members[0].Name")
fmt.Println(goval.GetAll[string](&team, path)) // [Alice]
path, _ := goval.Parse("Members[*].Name")
fmt.Println(goval.GetAll[string](&team, path)) // [Alice Bob]
```

### Set/SetFunc

```go
path, _ = goval.Parse("Name")
goval.Set[string](&team, path, "TEAM-B")
fmt.Println(team.Name) // TEAM-B

path, _ = goval.Parse("Members[*].Name")
goval.SetFunc[string](&team, path, func (v string, pathInfo goval.PathInfo) string{
    return strings.ToLower(v)
})
fmt.Println(team.Members[0].Name) // alice
fmt.Println(team.Members[1].Name) // bob
```

## Feature

- [ ] Map field support
- [ ] Reduce function

## Licence

This software is released under the MIT License, see LICENSE.
