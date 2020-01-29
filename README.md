# gonforce
Gonforce ( Read as Go Enforce ) enforces the given package follows the guidelines defined in "gonforce.yaml" file.

Sample file may look something like the following

```
package: github.com/hasitpbhatt/gonforce

default:
  allow:
    - gopkg.in/yaml.v2
  disallow:
    - gopkg.in/yaml.v1
```

Here we are ensuring the package can't use "gopkg.in/yaml.v1" but is allowed to use "gopkg.in/yaml.v2".
